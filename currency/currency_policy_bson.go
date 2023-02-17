package currency

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/spikeekips/mitum-currency/digest/util/bson"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

func (po CurrencyPolicy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":                   po.Hint().String(),
			"new_account_min_balance": po.newAccountMinBalance.String(),
			"feeer":                   po.feeer,
		})
}

type CurrencyPolicyBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	MinBalance string   `bson:"new_account_min_balance"`
	Feeer      bson.Raw `bson:"feeer"`
}

func (po *CurrencyPolicy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CurrencyPolicy")

	var upo CurrencyPolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(upo.Hint)
	if err != nil {
		return e(err, "")
	}

	return po.unpack(enc, ht, upo.MinBalance, upo.Feeer)
}
