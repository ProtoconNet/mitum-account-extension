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

var createContractAccountsItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateContractAccountsItemProcessor)
	},
}

var createContractAccountsProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateContractAccountsProcessor)
	},
}

func (CreateContractAccounts) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	// NOTE Process is nil func
	return nil, nil, nil
}

type CreateContractAccountsItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   CreateContractAccountsItem
	ns     base.StateMergeValue
	oas    base.StateMergeValue
	oac    currency.Account
	nb     map[currency.CurrencyID]base.StateMergeValue
}

func (opp *CreateContractAccountsItemProcessor) PreProcess(
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

	st, err = notExistsState(StateKeyContractAccount(target), "key of target contract account", getStateFunc)
	if err != nil {
		return err
	}
	opp.oas = NewContractAccountStateMergeValue(st.Key(), st.Value())

	st, err = existsState(currency.StateKeyAccount(opp.sender), "key of sender account", getStateFunc)
	if err != nil {
		return err
	}
	oac, err := currency.LoadStateAccountValue(st)
	if err != nil {
		return err
	}
	opp.oac = oac

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

func (opp *CreateContractAccountsItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, len(opp.item.Amounts())+2)

	nac, err := currency.NewAccountFromKeys(opp.item.Keys())
	if err != nil {
		return nil, err
	}
	ks := NewContractAccountKeys()
	ncac, err := nac.SetKeys(ks)
	if err != nil {
		return nil, err
	}
	sts[0] = currency.NewAccountStateMergeValue(opp.ns.Key(), currency.NewAccountStateValue(ncac))

	cas := NewContractAccount(opp.oac.Address(), true)
	sts[1] = NewContractAccountStateMergeValue(opp.oas.Key(), NewContractAccountStateValue(cas))

	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]
		v, ok := opp.nb[am.Currency()].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, errors.Errorf("expected BalanceStateValue, not %T", opp.nb[am.Currency()].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Add(am.Big())))
		sts[i+2] = currency.NewBalanceStateMergeValue(opp.nb[am.Currency()].Key(), stv)
	}

	return sts, nil
}

func (opp *CreateContractAccountsItemProcessor) Close() error {
	opp.h = nil
	opp.item = nil
	opp.ns = nil
	opp.nb = nil
	opp.sender = nil
	opp.oas = nil
	opp.oac = currency.Account{}

	createContractAccountsItemProcessorPool.Put(opp)

	return nil
}

type CreateContractAccountsProcessor struct {
	*base.BaseOperationProcessor
}

func NewCreateContractAccountsProcessor() GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new CreateContractAccountsProcessor")

		nopp := createContractAccountsProcessorPool.Get()
		opp, ok := nopp.(*CreateContractAccountsProcessor)
		if !ok {
			return nil, e(nil, "expected CreateContractAccountsProcessor, not %T", nopp)
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

func (opp *CreateContractAccountsProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(CreateContractAccountsFact)
	if !ok {
		return ctx, nil, errors.Errorf("expected CreateContractAccountsFact, not %T", op.Fact())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.sender), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.sender, err), nil
	}

	if err := checkNotExistsState(StateKeyContractAccount(fact.sender), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot be create-contract-account sender, %q: %w", fact.sender, err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *CreateContractAccountsProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, ok := op.Fact().(CreateContractAccountsFact)
	if !ok {
		return nil, nil, errors.Errorf("expected CreateContractAccountsFact, not %T", op.Fact())
	}

	required, err := opp.calculateItemsFee(op, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee: %w", err), nil
	}

	sb, err := CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	ns := make([]*CreateContractAccountsItemProcessor, len(fact.items))
	for i := range fact.items {
		cip := createContractAccountsItemProcessorPool.Get()
		c, ok := cip.(*CreateContractAccountsItemProcessor)
		if !ok {
			return nil, nil, errors.Errorf("expected CreateContractAccountsItemProcessor, not %T", cip)
		}

		c.h = op.Hash()
		c.item = fact.items[i]
		c.sender = fact.sender

		if err := c.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess CreateContractAccountsItem: %w", err), nil
		}

		ns[i] = c
	}

	var sts []base.StateMergeValue // nolint:prealloc
	for i := range ns {
		s, err := ns[i].Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process CreateContractAccountsItem: %w", err), nil
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

func (opp *CreateContractAccountsProcessor) Close() error {
	createContractAccountsProcessorPool.Put(opp)

	return nil
}

func (opp *CreateContractAccountsProcessor) calculateItemsFee(op base.Operation, getStateFunc base.GetStateFunc) (map[currency.CurrencyID][2]currency.Big, error) {
	fact, ok := op.Fact().(CreateContractAccountsFact)
	if !ok {
		return nil, errors.Errorf("expected CreateContractAccountsFact, not %T", op.Fact())
	}

	items := make([]currency.AmountsItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateItemsFee(getStateFunc, items)
}
