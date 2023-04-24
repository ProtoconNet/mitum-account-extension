package currency

import (
	"context"
	"sync"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var transfersItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransfersItemProcessor)
	},
}

var transfersProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransfersProcessor)
	},
}

type TransfersItemProcessor struct {
	h    util.Hash
	item mitumcurrency.TransfersItem
	rb   map[mitumcurrency.CurrencyID]base.StateMergeValue
}

func (opp *TransfersItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	if _, err := existsState(mitumcurrency.StateKeyAccount(opp.item.Receiver()), "key of receiver account", getStateFunc); err != nil {
		return err
	}

	rb := map[mitumcurrency.CurrencyID]base.StateMergeValue{}
	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]

		_, err := existsCurrencyPolicy(am.Currency(), getStateFunc)
		if err != nil {
			return err
		}

		st, err := existsState(mitumcurrency.StateKeyBalance(opp.item.Receiver(), am.Currency()), "key of receiver balance", getStateFunc)
		if err != nil {
			return nil
		}

		balance, err := mitumcurrency.StateBalanceValue(st)
		if err != nil {
			return err
		}

		rb[am.Currency()] = mitumcurrency.NewBalanceStateMergeValue(st.Key(), mitumcurrency.NewBalanceStateValue(balance))
	}

	opp.rb = rb

	return nil
}

func (opp *TransfersItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, len(opp.item.Amounts()))
	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]
		v, ok := opp.rb[am.Currency()].Value().(mitumcurrency.BalanceStateValue)
		if !ok {
			return nil, errors.Errorf("expect BalanceStateValue, not %T", opp.rb[am.Currency()].Value())
		}
		stv := mitumcurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Add(am.Big())))
		sts[i] = mitumcurrency.NewBalanceStateMergeValue(opp.rb[am.Currency()].Key(), stv)
	}

	return sts, nil
}

func (opp *TransfersItemProcessor) Close() error {
	opp.h = nil
	opp.item = nil
	opp.rb = nil

	transfersItemProcessorPool.Put(opp)

	return nil
}

type TransfersProcessor struct {
	*base.BaseOperationProcessor
}

func NewTransfersProcessor() GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new TransfersProcessor")

		nopp := transfersProcessorPool.Get()
		opp, ok := nopp.(*TransfersProcessor)
		if !ok {
			return nil, e(nil, "expected TransfersProcessor, not %T", nopp)
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

func (opp *TransfersProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess Transfers")

	fact, ok := op.Fact().(mitumcurrency.TransfersFact)
	if !ok {
		return ctx, nil, e(nil, "expected TransfersFact, not %T", op.Fact())
	}

	if err := checkExistsState(mitumcurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := checkNotExistsState(StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account cannot transfer amounts, %q: %w", fact.Sender(), err), nil
	}

	if err := checkFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *TransfersProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process Transfers")

	fact, ok := op.Fact().(mitumcurrency.TransfersFact)
	if !ok {
		return nil, nil, e(nil, "expected TransfersFact, not %T", op.Fact())
	}

	required, err := opp.calculateItemsFee(op, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to calculate fee: %w", err), nil
	}

	sb, err := CheckEnoughBalance(fact.Sender(), required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance: %w", err), nil
	}

	ns := make([]*TransfersItemProcessor, len(fact.Items()))
	for i := range fact.Items() {
		cip := transfersItemProcessorPool.Get()
		c, ok := cip.(*TransfersItemProcessor)
		if !ok {
			return nil, nil, e(nil, "expected TransfersItemProcessor, not %T", cip)
		}

		c.h = op.Hash()
		c.item = fact.Items()[i]

		if err := c.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("fail to preprocess transfer item: %w", err), nil
		}

		ns[i] = c
	}

	var sts []base.StateMergeValue // nolint:prealloc
	for i := range ns {
		s, err := ns[i].Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to process transfer item: %w", err), nil
		}
		sts = append(sts, s...)
	}

	for k := range required {
		rq := required[k]
		v, ok := sb[k].Value().(mitumcurrency.BalanceStateValue)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError("failed to process transfer"), nil
		}
		stv := mitumcurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(rq[0])))
		sts = append(sts, mitumcurrency.NewBalanceStateMergeValue(sb[k].Key(), stv))
	}

	return sts, nil, nil
}

func (opp *TransfersProcessor) Close() error {
	transfersProcessorPool.Put(opp)

	return nil
}

func (opp *TransfersProcessor) calculateItemsFee(op base.Operation, getStateFunc base.GetStateFunc) (map[mitumcurrency.CurrencyID][2]mitumcurrency.Big, error) {
	fact, ok := op.Fact().(mitumcurrency.TransfersFact)
	if !ok {
		return nil, errors.Errorf("expected TransfersFact, not %T", op.Fact())
	}
	items := make([]mitumcurrency.AmountsItem, len(fact.Items()))
	for i := range fact.Items() {
		items[i] = fact.Items()[i]
	}

	return CalculateItemsFee(getStateFunc, items)
}
