package currency

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/isaac"
	"github.com/spikeekips/mitum/util"
)

var createAccountsItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateAccountsItemProcessor)
	},
}

var createAccountsProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateAccountsProcessor)
	},
}

type CreateAccountsItemProcessor struct {
	h    util.Hash
	item currency.CreateAccountsItem
	ns   base.StateMergeValue
	nb   map[currency.CurrencyID]base.StateMergeValue
}

func (opp *CreateAccountsItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]

		policy, err := existsCurrencyPolicy(am.Currency(), getStateFunc)
		if err != nil {
			return err
		}

		if am.Big().Compare(policy.NewAccountMinBalance()) < 0 {
			return errors.Errorf("amount should be over minimum balance, %v < %v", am.Big(), policy.NewAccountMinBalance())
		}
	}

	target, err := opp.item.Address()
	if err != nil {
		return err
	}

	st, err := notExistsState(currency.StateKeyAccount(target), "key of target account", getStateFunc)
	if err != nil {
		return err
	}
	opp.ns = currency.NewAccountStateMergeValue(st.Key(), st.Value())

	nb := map[currency.CurrencyID]base.StateMergeValue{}
	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]
		switch _, found, err := getStateFunc(currency.StateKeyBalance(target, am.Currency())); {
		case err != nil:
			return err
		case found:
			return isaac.ErrStopProcessingRetry.Errorf("target balance already exists, %q", target)
		default:
			nb[am.Currency()] = currency.NewBalanceStateMergeValue(currency.StateKeyBalance(target, am.Currency()), currency.NewBalanceStateValue(currency.NewZeroAmount(am.Currency())))
		}
	}
	opp.nb = nb

	return nil
}

func (opp *CreateAccountsItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	nac, err := currency.NewAccountFromKeys(opp.item.Keys())
	if err != nil {
		return nil, err
	}

	sts := make([]base.StateMergeValue, len(opp.item.Amounts())+1)
	sts[0] = currency.NewAccountStateMergeValue(opp.ns.Key(), currency.NewAccountStateValue(nac))

	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]
		v, ok := opp.nb[am.Currency()].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, errors.Errorf("expected BalanceStateValue, not %T", opp.nb[am.Currency()].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Add(am.Big())))
		sts[i+1] = currency.NewBalanceStateMergeValue(opp.nb[am.Currency()].Key(), stv)
	}

	return sts, nil
}

func (opp *CreateAccountsItemProcessor) Close() error {
	opp.h = nil
	opp.item = nil
	opp.ns = nil
	opp.nb = nil

	createAccountsItemProcessorPool.Put(opp)

	return nil
}

type CreateAccountsProcessor struct {
	*base.BaseOperationProcessor
}

func NewCreateAccountsProcessor() GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new CreateAccountsProcessor")

		nopp := createAccountsProcessorPool.Get()
		opp, ok := nopp.(*CreateAccountsProcessor)
		if !ok {
			return nil, e(nil, "expected CreateAccountsProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *CreateAccountsProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(currency.CreateAccountsFact)
	if !ok {
		return ctx, nil, errors.Errorf("expected CreateAccountsFact, not %T", op.Fact())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot be create-account sender, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *CreateAccountsProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, ok := op.Fact().(currency.CreateAccountsFact)
	if !ok {
		return nil, nil, errors.Errorf("expected CreateAccountsFact, not %T", op.Fact())
	}

	required, err := opp.calculateItemsFee(op, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee: %w", err), nil
	}

	sb, err := CheckEnoughBalance(fact.Sender(), required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	ns := make([]*CreateAccountsItemProcessor, len(fact.Items()))
	for i := range fact.Items() {
		cip := createAccountsItemProcessorPool.Get()
		c, ok := cip.(*CreateAccountsItemProcessor)
		if !ok {
			return nil, nil, errors.Errorf("expected CreateAccountsItemProcessor, not %T", cip)
		}

		c.h = op.Hash()
		c.item = fact.Items()[i]

		if err := c.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess CreateAccountsItem: %w", err), nil
		}

		ns[i] = c
	}

	var sts []base.StateMergeValue // nolint:prealloc
	for i := range ns {
		s, err := ns[i].Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process CreateAccountsItem: %w", err), nil
		}
		sts = append(sts, s...)

		ns[i].Close()
	}

	for i := range sb {
		v, ok := sb[i].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, errors.Errorf("expected BalanceStateValue, not %T", sb[i].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currency.NewBalanceStateMergeValue(sb[i].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *CreateAccountsProcessor) Close() error {
	createAccountsProcessorPool.Put(opp)

	return nil
}

func (opp *CreateAccountsProcessor) calculateItemsFee(op base.Operation, getStateFunc base.GetStateFunc) (map[currency.CurrencyID][2]currency.Big, error) {
	fact, ok := op.Fact().(currency.CreateAccountsFact)
	if !ok {
		return nil, errors.Errorf("expected CreateAccountsFact, not %T", op.Fact())
	}

	items := make([]currency.AmountsItem, len(fact.Items()))
	for i := range fact.Items() {
		items[i] = fact.Items()[i]
	}

	return CalculateItemsFee(getStateFunc, items)
}

func CalculateItemsFee(getStateFunc base.GetStateFunc, items []currency.AmountsItem) (map[currency.CurrencyID][2]currency.Big, error) {
	required := map[currency.CurrencyID][2]currency.Big{}

	for i := range items {
		it := items[i]

		for j := range it.Amounts() {
			am := it.Amounts()[j]

			rq := [2]currency.Big{currency.ZeroBig, currency.ZeroBig}
			if k, found := required[am.Currency()]; found {
				rq = k
			}

			policy, err := existsCurrencyPolicy(am.Currency(), getStateFunc)
			if err != nil {
				return nil, err
			}

			switch k, err := policy.Feeer().Fee(am.Big()); {
			case err != nil:
				return nil, err
			case !k.OverZero():
				required[am.Currency()] = [2]currency.Big{rq[0].Add(am.Big()), rq[1]}
			default:
				required[am.Currency()] = [2]currency.Big{rq[0].Add(am.Big()).Add(k), rq[1].Add(k)}
			}
		}
	}

	return required, nil
}

func CheckEnoughBalance(
	holder base.Address,
	required map[currency.CurrencyID][2]currency.Big,
	getStateFunc base.GetStateFunc,
) (map[currency.CurrencyID]base.StateMergeValue, error) {
	sb := map[currency.CurrencyID]base.StateMergeValue{}

	for cid := range required {
		rq := required[cid]

		st, err := existsState(currency.StateKeyBalance(holder, cid), "key of holder balance", getStateFunc)
		if err != nil {
			return nil, err
		}

		am, err := currency.StateBalanceValue(st)
		if err != nil {
			return nil, errors.Errorf("not enough balance of sender, %q: %w", holder, err)
		}

		if am.Big().Compare(rq[0]) < 0 {
			return nil, errors.Errorf("not enough balance of sender, %q; %v !> %v", holder, am.Big(), rq[0])
		}
		sb[cid] = currency.NewBalanceStateMergeValue(st.Key(), currency.NewBalanceStateValue(am))
	}

	return sb, nil
}
