package currency

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	AmountValueType   = hint.Type("mitum-currency-amount-value")
	AmountValueHint   = hint.NewHint(AmountValueType, "v0.0.1")
	AmountValueHinter = AmountValue{BaseHinter: hint.NewBaseHinter(AmountValueHint)}
)

type AmountValue struct {
	hint.BaseHinter
	amount     currency.Amount
	contractid ContractID
}

func NewAmountValue(big currency.Big, cid currency.CurrencyID, id ContractID) AmountValue {
	am := currency.NewAmount(big, cid)
	amv := AmountValue{BaseHinter: hint.NewBaseHinter(AmountValueHint), amount: am, contractid: id}

	return amv
}

func NewAmountValuefromAmount(am currency.Amount, id ContractID) AmountValue {
	amv := AmountValue{BaseHinter: hint.NewBaseHinter(AmountValueHint), amount: am, contractid: id}

	return amv
}

func MustNewAmountValue(big currency.Big, cid currency.CurrencyID, id ContractID) AmountValue {
	amv := NewAmountValue(big, cid, id)
	if err := amv.IsValid(nil); err != nil {
		panic(err)
	}

	return amv
}

func (am AmountValue) Bytes() []byte {
	return util.ConcatBytesSlice(
		am.amount.Bytes(),
		am.contractid.Bytes(),
	)
}

func (am AmountValue) Hash() valuehash.Hash {
	return am.GenerateHash()
}

func (am AmountValue) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(am.Bytes())
}

func (am AmountValue) IsEmpty() bool {
	return len(am.amount.Currency()) < 1 || !am.amount.Big().OverNil() || len(am.contractid) < 1
}

func (am AmountValue) IsValid([]byte) error {
	if err := isvalid.Check(nil, false,
		am.BaseHinter,
		am.amount,
		am.contractid,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid AmountValue: %w", err)
	}

	return nil
}

func (am AmountValue) Amount() currency.Amount {
	return am.amount
}

func (am AmountValue) Add(b currency.Big) (AmountValue, error) {
	v := currency.NewAmount(am.amount.Big().Add(b), am.amount.Currency())
	amv := NewAmountValuefromAmount(v, am.contractid)
	return amv, nil
}

func (am AmountValue) Sub(b currency.Big) (AmountValue, error) {
	if !(am.amount.Big().Sub(b)).OverNil() {
		return AmountValue{}, errors.Errorf("under zero amount after substraction, %v", am.amount.Big().Sub(b))
	}
	v := currency.NewAmount(am.amount.Big().Sub(b), am.amount.Currency())
	amv := NewAmountValuefromAmount(v, am.contractid)
	return amv, nil
}

/*
func (am AmountValue) Big() currency.Big {
	return am.amount.Big()
}

func (am AmountValue) Currency() currency.CurrencyID {
	return am.amount.Currency()
}
*/

func (am AmountValue) ID() ContractID {
	return am.contractid
}

func (am AmountValue) String() string {
	return fmt.Sprintf("%s-%s", am.amount.String(), am.contractid)
}

func (am AmountValue) Equal(b AmountValue) bool {
	switch {
	case !am.amount.Equal(b.amount):
		return false
	case !(am.contractid != b.contractid):
		return false
	default:
		return true
	}
}

func (am AmountValue) WithBig(big currency.Big) AmountValue {
	a := am.amount.WithBig(big)
	amv := NewAmountValue(a.Big(), a.Currency(), am.contractid)

	return amv
}
