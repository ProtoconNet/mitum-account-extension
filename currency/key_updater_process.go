package currency

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
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
	fact, ok := op.Fact().(currency.KeyUpdaterFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError("expected KeyUpdaterFact, not %T", op.Fact()), nil
	}

	if err := checkFactSignsByState(fact.Target(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	st, err := existsState(currency.StateKeyAccount(fact.Target()), "target keys", getStateFunc)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("failed to check existence of target %v: %w", fact.Target(), err), nil
	}

	if err := checkNotExistsState(StateKeyContractAccount(fact.Target()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("no keys in contract account %v: %w", fact.Target(), err), nil
	}

	ks, err := currency.StateKeysValue(st)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("failed to get state value of keys %v: %w", fact.Keys().Hash(), err), nil
	}
	if ks.Equal(fact.Keys()) {
		return ctx, base.NewBaseOperationProcessReasonError("same Keys as existing %v: %w", fact.Keys().Hash(), err), nil
	}

	return ctx, nil, nil
}

func (opp *KeyUpdaterProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, ok := op.Fact().(currency.KeyUpdaterFact)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected KeyUpdaterFact, not %T", op.Fact()), nil
	}

	st, err := existsState(currency.StateKeyAccount(fact.Target()), "target keys", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check existence of target %v: %w", fact.Target(), err), nil
	}
	sa := currency.NewAccountStateMergeValue(st.Key(), st.Value())

	policy, err := existsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check existence of currency %v: %w", fact.Currency(), err), nil
	}
	fee, err := policy.Feeer().Fee(currency.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency %v: %w", fact.Currency(), err), nil
	}

	st, err = existsState(currency.StateKeyBalance(fact.Target(), fact.Currency()), "balance of target", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check existence of targe balance %v: %w", fact.Target(), err), nil
	}
	sb := currency.NewBalanceStateMergeValue(st.Key(), st.Value())

	switch b, err := currency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to check existence of target balance %v,%v: %w", fact.Currency(), fact.Target(), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("insufficient balance with fee %v,%v", fact.Currency(), fact.Target()), nil
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
