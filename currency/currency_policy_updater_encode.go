package currency

import (
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CurrencyPolicyUpdaterFact) unpack(enc encoder.Encoder, cid string, bpo []byte) error {
	e := util.StringErrorFunc("failed to unmarshal CurrencyPolicyUpdaterFact")

	if hinter, err := enc.Decode(bpo); err != nil {
		return e(err, "")
	} else if po, ok := hinter.(CurrencyPolicy); !ok {
		return e(util.ErrWrongType.Errorf("expected CurrencyPolicy, not %T", hinter), "")
	} else {
		fact.policy = po
	}

	fact.currency = mitumcurrency.CurrencyID(cid)

	return nil
}
