package currency

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
)

var withdrawsItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(WithdrawsItemProcessor)
	},
}

var withdrawsProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(WithdrawsProcessor)
	},
}

func (Withdraws) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	// NOTE Process is nil func
	return nil, nil, nil
}

type WithdrawsItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   WithdrawsItem
	tb     map[currency.CurrencyID]base.StateMergeValue
}

func (opp *WithdrawsItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	if err := checkExistsState(currency.StateKeyAccount(opp.item.Target()), getStateFunc); err != nil {
		return err
	}

	st, err := existsState(StateKeyContractAccount(opp.item.Target()), "key of target contract account", getStateFunc)
	if err != nil {
		return err
	}
	v, err := StateContractAccountValue(st)
	if err != nil {
		return err
	}
	if !v.owner.Equal(opp.sender) {
		return errors.Errorf("contract account owner is not matched with %q", opp.sender)
	}

	tb := map[currency.CurrencyID]base.StateMergeValue{}
	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]

		_, err := existsCurrencyPolicy(am.Currency(), getStateFunc)
		if err != nil {
			return err
		}

		st, _, err := getStateFunc(currency.StateKeyBalance(opp.item.Target(), am.Currency()))
		if err != nil {
			return err
		}

		balance, err := currency.StateBalanceValue(st)
		if err != nil {
			return err
		}

		tb[am.Currency()] = currency.NewBalanceStateMergeValue(st.Key(), currency.NewBalanceStateValue(balance))
	}

	opp.tb = tb

	return nil
}

func (opp *WithdrawsItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, len(opp.item.Amounts()))
	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]
		v, ok := opp.tb[am.Currency()].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, errors.Errorf("expect BalanceStateValue, not %T", opp.tb[am.Currency()].Value())
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(am.Big())))
		sts[i] = currency.NewBalanceStateMergeValue(opp.tb[am.Currency()].Key(), stv)
	}

	return sts, nil
}

func (opp *WithdrawsItemProcessor) Close() error {
	opp.h = nil
	opp.sender = nil
	opp.item = nil
	opp.tb = nil

	withdrawsItemProcessorPool.Put(opp)

	return nil
}

type WithdrawsProcessor struct {
	*base.BaseOperationProcessor
}

func NewWithdrawsProcessor() GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new WithdrawsProcessor")

		nopp := withdrawsProcessorPool.Get()
		opp, ok := nopp.(*WithdrawsProcessor)
		if !ok {
			return nil, e(nil, "expected WithdrawsProcessor, not %T", nopp)
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

func (opp *WithdrawsProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess Withdraws")

	fact, ok := op.Fact().(WithdrawsFact)
	if !ok {
		return ctx, nil, e(nil, "expected WithdrawsFact, not %T", op.Fact())
	}

	if err := checkExistsState(currency.StateKeyAccount(fact.sender), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.sender, err), nil
	}

	if err := checkNotExistsState(StateKeyContractAccount(fact.sender), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot be ca withdraw sender, %q: %w", fact.sender, err), nil
	}

	if err := checkFactSignsByState(fact.sender, op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *WithdrawsProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process Withdraws")

	fact, ok := op.Fact().(WithdrawsFact)
	if !ok {
		return nil, nil, e(nil, "expected WithdrawsFact, not %T", op.Fact())
	}

	required, err := opp.calculateItemsFee(op, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee: %w", err), nil
	}
	sb, err := CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	ns := make([]*WithdrawsItemProcessor, len(fact.items))
	for i := range fact.items {
		cip := withdrawsItemProcessorPool.Get()
		c, ok := cip.(*WithdrawsItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected WithdrawsItemProcessor, not %T", cip)
		}

		c.h = op.Hash()
		c.sender = fact.sender
		c.item = fact.items[i]

		if err := c.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess WithdrawsItem: %w", err), nil
		}

		ns[i] = c
	}

	var sts []base.StateMergeValue // nolint:prealloc
	for i := range ns {
		s, err := ns[i].Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process WithdrawsItem: %w", err), nil
		}
		sts = append(sts, s...)

		ns[i].Close()
	}

	for k := range required {
		rq := required[k]
		v, ok := sb[k].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError("failed to process Withdraws: expected BalanceStateValue, not %T", sb[k].Value()), nil
		}
		stv := currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Add(rq[0]).Sub(rq[1].MulInt64(2))))
		sts = append(sts, currency.NewBalanceStateMergeValue(sb[k].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *WithdrawsProcessor) Close() error {
	withdrawsProcessorPool.Put(opp)

	return nil
}

func (opp *WithdrawsProcessor) calculateItemsFee(op base.Operation, getStateFunc base.GetStateFunc) (map[currency.CurrencyID][2]currency.Big, error) {
	fact, ok := op.Fact().(WithdrawsFact)
	if !ok {
		return nil, errors.Errorf("expected WithdrawsFact, not %T", op.Fact())
	}
	items := make([]currency.AmountsItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateItemsFee(getStateFunc, items)
}
