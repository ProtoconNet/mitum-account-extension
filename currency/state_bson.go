package currency

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s CurrencyDesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":          s.Hint().String(),
			"currencydesign": s.CurrencyDesign,
		},
	)
}

type CurrencyDesignStateValueBSONUnmarshaler struct {
	Hint           string   `bson:"_hint"`
	CurrencyDesign bson.Raw `bson:"currencydesign"`
}

func (s *CurrencyDesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CurrencyDesignStateValue")

	var u CurrencyDesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var cd CurrencyDesign
	if err := cd.DecodeBSON(u.CurrencyDesign, enc); err != nil {
		return e(err, "")
	}

	s.CurrencyDesign = cd

	return nil
}

func (s ContractAccountStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":           s.Hint().String(),
			"contractaccount": s.account,
		},
	)

}

type ContractAccountStateValueBSONUnmarshaler struct {
	Hint            string   `bson:"_hint"`
	ContractAccount bson.Raw `bson:"contractaccount"`
}

func (s *ContractAccountStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of ContractAccountStateValue")

	var u ContractAccountStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var ca ContractAccount
	if err := ca.DecodeBSON(u.ContractAccount, enc); err != nil {
		return e(err, "")
	}

	s.account = ca

	return nil
}
