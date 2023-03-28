package currency // nolint:dupl

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it BaseWithdrawsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":   it.Hint().String(),
			"target":  it.target,
			"amounts": it.amounts,
		},
	)
}

type WithdrawsItemBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	Target  string   `bson:"target"`
	Amounts bson.Raw `bson:"amounts"`
}

func (it *BaseWithdrawsItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of BaseWithdrawsItem")

	var uit WithdrawsItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return e(err, "")
	}

	return it.unpack(enc, ht, uit.Target, uit.Amounts)
}
