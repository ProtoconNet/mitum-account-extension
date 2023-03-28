package currency

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de *CurrencyDesign) unpack(enc encoder.Encoder, ht hint.Hint, bam []byte, ga string, bpo []byte, ag string) error {
	e := util.StringErrorFunc("failed to unmarshal CurrencyDesign")

	de.BaseHinter = hint.NewBaseHinter(ht)

	var am currency.Amount
	if err := encoder.Decode(enc, bam, &am); err != nil {
		return e(err, "failed to decode amount")
	}

	de.amount = am

	switch a, err := base.DecodeAddress(ga, enc); {
	case err != nil:
		return e(err, "failed to decode address")
	default:
		de.genesisAccount = a
	}

	var policy CurrencyPolicy

	if err := encoder.Decode(enc, bpo, &policy); err != nil {
		return e(err, "failed to decode currency policy")
	}

	de.policy = policy

	if big, err := currency.NewBigFromString(ag); err != nil {
		return e(err, "")
	} else {
		de.aggregate = big
	}

	return nil
}
