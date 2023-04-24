package currency

import (
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var (
	CurrencyPolicyHint = hint.MustNewHint("mitum-currency-currency-policy-v0.0.1")
)

type CurrencyPolicy struct {
	hint.BaseHinter
	newAccountMinBalance mitumcurrency.Big
	feeer                Feeer
}

func NewCurrencyPolicy(newAccountMinBalance mitumcurrency.Big, feeer Feeer) CurrencyPolicy {
	return CurrencyPolicy{
		BaseHinter:           hint.NewBaseHinter(CurrencyPolicyHint),
		newAccountMinBalance: newAccountMinBalance, feeer: feeer,
	}
}

func (po CurrencyPolicy) Bytes() []byte {
	return util.ConcatBytesSlice(po.newAccountMinBalance.Bytes(), po.feeer.Bytes())
}

func (po CurrencyPolicy) IsValid([]byte) error {
	if !po.newAccountMinBalance.OverNil() {
		return util.ErrInvalid.Errorf("NewAccountMinBalance under zero")
	}

	if err := util.CheckIsValiders(nil, false, po.BaseHinter, po.feeer); err != nil {
		return util.ErrInvalid.Errorf("invalid currency policy: %w", err)
	}

	return nil
}

func (po CurrencyPolicy) NewAccountMinBalance() mitumcurrency.Big {
	return po.newAccountMinBalance
}

func (po CurrencyPolicy) Feeer() Feeer {
	return po.feeer
}
