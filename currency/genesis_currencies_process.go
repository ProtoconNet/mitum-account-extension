package currency

import (
	"context"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
)

func (GenesisCurrencies) PreProcess(
	ctx context.Context, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	return ctx, nil, nil
}

func (op GenesisCurrencies) Process(
	ctx context.Context, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process GenesisCurrencies")

	fact, ok := op.Fact().(GenesisCurrenciesFact)
	if !ok {
		return nil, nil, e(nil, "expected GenesisCurrenciesFact, not %T", op.Fact())
	}

	newAddress, err := fact.Address()
	if err != nil {
		return nil, nil, e(err, "")
	}

	ns, err := notExistsState(currency.StateKeyAccount(newAddress), "key of genesis account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("genesis account already exists, %q: %w", newAddress, err), nil
	}

	cs := make([]CurrencyDesign, len(fact.cs))
	gas := map[currency.CurrencyID]base.StateMergeValue{}
	sts := map[currency.CurrencyID]base.StateMergeValue{}
	for i := range fact.cs {
		c := fact.cs[i]
		c.genesisAccount = newAddress
		cs[i] = c

		st, err := notExistsState(StateKeyCurrencyDesign(c.Currency()), "currency", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("currency design already exists, %q: %w", c.Currency(), err), nil
		}

		sts[c.Currency()] = NewCurrencyDesignStateMergeValue(st.Key(), NewCurrencyDesignStateValue(c))

		st, err = notExistsState(currency.StateKeyBalance(newAddress, c.Currency()), "key of genesis balance", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("account balance already exists, %q: %w", newAddress, err), nil
		}
		gas[c.Currency()] = currency.NewBalanceStateMergeValue(st.Key(), currency.NewBalanceStateValue(currency.NewZeroAmount(c.Currency())))
	}

	var smvs []base.StateMergeValue
	if ac, err := currency.NewAccount(newAddress, fact.keys); err != nil {
		return nil, nil, e(err, "")
	} else {
		smvs = append(smvs, currency.NewAccountStateMergeValue(ns.Key(), currency.NewAccountStateValue(ac)))
	}

	for i := range cs {
		c := cs[i]
		v, ok := gas[c.Currency()].Value().(currency.BalanceStateValue)
		if !ok {
			return nil, nil, e(nil, "expected BalanceStateValue, not %T", gas[c.Currency()].Value())
		}

		gst := currency.NewBalanceStateMergeValue(gas[c.Currency()].Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Add(c.amount.Big()))))
		dst := NewCurrencyDesignStateMergeValue(sts[c.Currency()].Key(), NewCurrencyDesignStateValue(c))
		smvs = append(smvs, gst, dst)

		sts, err := createZeroAccount(c.Currency(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to create zero account, %q: %w", c.Currency(), err), nil
		}

		smvs = append(smvs, sts...)
	}

	return smvs, nil, nil
}

func createZeroAccount(
	cid currency.CurrencyID,
	getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	sts := make([]base.StateMergeValue, 2)

	ac, err := currency.ZeroAccount(cid)
	if err != nil {
		return nil, err
	}
	ast, err := notExistsState(currency.StateKeyAccount(ac.Address()), "key of zero account", getStateFunc)
	if err != nil {
		return nil, err
	}

	sts[0] = currency.NewAccountStateMergeValue(ast.Key(), currency.NewAccountStateValue(ac))

	bst, err := notExistsState(currency.StateKeyBalance(ac.Address(), cid), "key of zero account balance", getStateFunc)
	if err != nil {
		return nil, err
	}

	sts[1] = currency.NewBalanceStateMergeValue(bst.Key(), currency.NewBalanceStateValue(currency.NewZeroAmount(cid)))

	return sts, nil
}
