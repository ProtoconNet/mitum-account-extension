package currency

import (
	"github.com/spikeekips/mitum-currency/currency"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

type AmountValueBSONPacker struct {
	AM currency.Amount `bson:"amount"`
	CI ContractID      `bson:"contractid"`
}

func (am AmountValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(am.Hint()),
		bson.M{
			"amount":     am.amount,
			"contractid": am.contractid,
		}),
	)
}

type AmountValueBSONUnpacker struct {
	HT hint.Hint `bson:"_hint"`
	AM bson.Raw  `bson:"amount"`
	CI string    `bson:"contractid"`
}

func (am *AmountValue) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uam AmountValueBSONUnpacker
	if err := bsonenc.Unmarshal(b, &uam); err != nil {
		return err
	}
	return am.unpack(enc, uam.AM, uam.CI)
}
