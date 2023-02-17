package currency // nolint:dupl

import (
	bsonenc "github.com/spikeekips/mitum-currency/digest/util/bson"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it BaseCreateContractAccountsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":   it.Hint().String(),
			"keys":    it.keys,
			"amounts": it.amounts,
		},
	)
}

type CreateContractAccountsItemBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	Keys    bson.Raw `bson:"keys"`
	Amounts bson.Raw `bson:"amounts"`
}

func (it *BaseCreateContractAccountsItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of BaseCreateContractAccountsItem")

	var uit CreateContractAccountsItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unpack(enc, ht, uit.Keys, uit.Amounts)
}
