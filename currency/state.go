package currency

import (
	"fmt"
	"strings"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var CurrencyDesignStateValueHint = hint.MustNewHint("currency-design-state-value-v0.0.1")

var StateKeyCurrencyDesignPrefix = "currencydesign:"

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

func StateKeyCurrencyDesign(cid mitumcurrency.CurrencyID) string {
	return fmt.Sprintf("%s%s", StateKeyCurrencyDesignPrefix, cid)
}

var ContractAccountStateValueHint = hint.MustNewHint("contract-account-state-value-v0.0.1")

var StateKeyContractAccountSuffix = ":contractaccount"

type ContractAccountStateValue struct {
	hint.BaseHinter
	account ContractAccount
}

func NewContractAccountStateValue(account ContractAccount) ContractAccountStateValue {
	return ContractAccountStateValue{
		BaseHinter: hint.NewBaseHinter(ContractAccountStateValueHint),
		account:    account,
	}
}

func (c ContractAccountStateValue) Hint() hint.Hint {
	return c.BaseHinter.Hint()
}

func (c ContractAccountStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid ContractAccountStateValue")

	if err := c.BaseHinter.IsValid(ContractAccountStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false, c.account); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (c ContractAccountStateValue) HashBytes() []byte {
	return c.account.Bytes()
}

func StateKeyContractAccount(a base.Address) string {
	return fmt.Sprintf("%s%s", a.String(), StateKeyContractAccountSuffix)
}

func IsStateContractAccountKey(key string) bool {
	return strings.HasSuffix(key, StateKeyContractAccountSuffix)
}

func StateContractAccountValue(st base.State) (ContractAccount, error) {
	v := st.Value()
	if v == nil {
		return ContractAccount{}, util.ErrNotFound.Errorf("contract account status not found in State")
	}

	cs, ok := v.(ContractAccountStateValue)
	if !ok {
		return ContractAccount{}, errors.Errorf("invalid contract account status value found, %T", v)
	}
	return cs.account, nil
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

type ContractAccountStateValueMerger struct {
	*base.BaseStateValueMerger
}

func NewContractAccountStateValueMerger(height base.Height, key string, st base.State) *ContractAccountStateValueMerger {
	s := &ContractAccountStateValueMerger{
		BaseStateValueMerger: base.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewContractAccountStateMergeValue(key string, stv base.StateValue) base.StateMergeValue {
	return base.NewBaseStateMergeValue(
		key,
		stv,
		func(height base.Height, st base.State) base.StateValueMerger {
			return NewContractAccountStateValueMerger(height, key, st)
		},
	)
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

func checkNotExistsState(
	key string,
	getState base.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case found:
		return base.NewBaseOperationProcessReasonError("state, %q already exists", key)
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

func existsCurrencyPolicy(cid mitumcurrency.CurrencyID, getStateFunc base.GetStateFunc) (CurrencyPolicy, error) {
	var policy CurrencyPolicy
	switch i, found, err := getStateFunc(StateKeyCurrencyDesign(cid)); {
	case err != nil:
		return CurrencyPolicy{}, err
	case !found:
		return CurrencyPolicy{}, errors.Errorf("currency not found, %v", cid)
	default:
		currencydesign, ok := i.Value().(CurrencyDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return CurrencyPolicy{}, errors.Errorf("expected CurrencyDesignStateValue, not %T", i.Value())
		}
		policy = currencydesign.CurrencyDesign.policy
	}
	return policy, nil
}
