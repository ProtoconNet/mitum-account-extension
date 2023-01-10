package currency

import (
	"github.com/spikeekips/mitum/util"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s CurrencyDesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(
			bsonenc.NewHintedDoc(s.Hint()),
			bson.M{
				"currencydesign": s.CurrencyDesign,
			},
		))

}

type CurrencyDesignStateValueBSONUnmarshaler struct {
	HT hint.Hint `bson:"_hint"`
	CD bson.Raw  `bson:"currencydesign"`
}

func (s *CurrencyDesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode CurrencyDesignStateValue")

	var u CurrencyDesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.HT)

	var cd CurrencyDesign
	if err := cd.DecodeBSON(u.CD, enc); err != nil {
		return e(err, "")
	}

	s.CurrencyDesign = cd

	return nil
}
