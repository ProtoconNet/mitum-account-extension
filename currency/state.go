package currency

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

var CurrencyDesignStateValueHint = hint.MustNewHint("currency-design-state-value-v0.0.1")

var StateKeyCurrencyDesignPrefix = "extensioncurrencydesign:"

type CurrencyDesignStateValue struct {
	hint.BaseHinter
	CurrencyDesign CurrencyDesign
}

func NewCurrencyDesignStateValue(currencyDesign CurrencyDesign) CurrencyDesignStateValue {
	return CurrencyDesignStateValue{
		BaseHinter:     hint.NewBaseHinter(CurrencyDesignStateValueHint),
		CurrencyDesign: currencyDesign,
	}
}

func (c CurrencyDesignStateValue) Hint() hint.Hint {
	return c.BaseHinter.Hint()
}

func (c CurrencyDesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid CurrencyDesignStateValue")

	if err := c.BaseHinter.IsValid(CurrencyDesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false, c.CurrencyDesign); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (c CurrencyDesignStateValue) HashBytes() []byte {
	return c.CurrencyDesign.Bytes()
}

func StateCurrencyDesignValue(st base.State) (CurrencyDesign, error) {
	v := st.Value()
	if v == nil {
		return CurrencyDesign{}, util.ErrNotFound.Errorf("currency design not found in State")
	}

	de, ok := v.(CurrencyDesignStateValue)
	if !ok {
		return CurrencyDesign{}, errors.Errorf("invalid currency design value found, %T", v)
	}

	return de.CurrencyDesign, nil
}

func IsStateCurrencyDesignKey(key string) bool {
	return strings.HasPrefix(key, StateKeyCurrencyDesignPrefix)
}

func StateKeyCurrencyDesign(cid currency.CurrencyID) string {
	return fmt.Sprintf("%s%s", StateKeyCurrencyDesignPrefix, cid)
}

func checkExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return base.NewBaseOperationProcessReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func notExistsState(
	k,
	name string,
	getState base.GetStateFunc,
) (base.State, error) {
	var st base.State
	switch _, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, base.NewBaseOperationProcessReasonError("%s already exists", name)
	case !found:
		st = base.NewBaseState(base.NilHeight, k, nil, nil, nil)
	}
	return st, nil
}

func existsCurrencyPolicy(cid currency.CurrencyID, getStateFunc base.GetStateFunc) (CurrencyPolicy, error) {
	var policy CurrencyPolicy
	switch i, found, err := getStateFunc(StateKeyCurrencyDesign(cid)); {
	case err != nil:
		return CurrencyPolicy{}, err
	case !found:
		return CurrencyPolicy{}, base.NewBaseOperationProcessReasonError("currency not found, %v", cid)
	default:
		currencydesign, ok := i.Value().(CurrencyDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return CurrencyPolicy{}, errors.Errorf("expected CurrencyDesignStateValue, not %T", i.Value())
		}
		policy = currencydesign.CurrencyDesign.policy
	}
	return policy, nil
}

type CurrencyDesignStateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewCurrencyDesignStateValueMerger(height base.Height, key string, st base.State) *CurrencyDesignStateValueMerger {
	s := &CurrencyDesignStateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewCurrencyDesignStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	return base.NewBaseStateMergeValue(
		key,
		stv,
		func(height base.Height, st base.State) base.StateValueMerger {
			return NewCurrencyDesignStateValueMerger(height, key, st)
		},
	)
}
