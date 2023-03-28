package currency

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var keyUpdaterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(KeyUpdaterProcessor)
	},
}

type KeyUpdaterProcessor struct {
	*base.BaseOperationProcessor
}

func NewKeyUpdaterProcessor() GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new KeyUpdaterProcessor")

		nopp := keyUpdaterProcessorPool.Get()
		opp, ok := nopp.(*KeyUpdaterProcessor)
		if !ok {
			return nil, errors.Errorf("expected KeyUpdaterProcessor, not %T", nopp)
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

func (opp *KeyUpdaterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess KeyUpdater")

	fact, ok := op.Fact().(currency.KeyUpdaterFact)
	if !ok {
		return ctx, nil, e(nil, "expected KeyUpdaterFact, not %T", op.Fact())
	}

	st, err := existsState(currency.StateKeyAccount(fact.Target()), "key of target account", getStateFunc)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("target not found, %q: %w", fact.Target(), err), nil
	}

	if err := checkNotExistsState(StateKeyContractAccount(fact.Target()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("contract account already exists, %q: %w", fact.Target(), err), nil
	}

	ks, err := currency.StateKeysValue(st)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("failed to get keys value, %q: %w", fact.Keys().Hash(), err), nil
	}
	if ks.Equal(fact.Keys()) {
		return ctx, base.NewBaseOperationProcessReasonError("same Keys as existing, %q: %w", fact.Keys().Hash(), err), nil
	}

	if err := checkFactSignsByState(fact.Target(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *KeyUpdaterProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process KeyUpdater")

	fact, ok := op.Fact().(currency.KeyUpdaterFact)
	if !ok {
		return nil, nil, e(nil, "expected KeyUpdaterFact, not %T", op.Fact())
	}

	st, err := existsState(currency.StateKeyAccount(fact.Target()), "key of target account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("target not found, %q: %w", fact.Target(), err), nil
	}
	sa := currency.NewAccountStateMergeValue(st.Key(), st.Value())

	policy, err := existsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}
	fee, err := policy.Feeer().Fee(currency.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err = existsState(currency.StateKeyBalance(fact.Target(), fact.Currency()), "key of target balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("target balance not found, %q: %w", fact.Target(), err), nil
	}
	sb := currency.NewBalanceStateMergeValue(st.Key(), st.Value())

	switch b, err := currency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %q: %w", currency.StateKeyBalance(fact.Target(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("not enough balance of target, %q", fact.Target()), nil
	}

	var sts []base.StateMergeValue // nolint:prealloc

	v, ok := sb.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts = append(sts, currency.NewBalanceStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee)))))

	a, err := currency.NewAccountFromKeys(fact.Keys())
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to create new account from keys"), nil
	}
	sts = append(sts, currency.NewAccountStateMergeValue(sa.Key(), currency.NewAccountStateValue(a)))

	return sts, nil, nil
}

func (opp *KeyUpdaterProcessor) Close() error {
	keyUpdaterProcessorPool.Put(opp)

	return nil
}
