package currency // nolint:dupl

import (
	"github.com/spikeekips/mitum/util"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it BaseWithdrawsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"target":  it.target,
				"amounts": it.amounts,
			}),
	)
}

type WithdrawsItemBSONUnmarshaler struct {
	HT hint.Hint `bson:"_hint"`
	TG string    `bson:"target"`
	AM bson.Raw  `bson:"amounts"`
}

func (it *BaseWithdrawsItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of BaseWithdrawsItem")

	var uit WithdrawsItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.HT, uit.TG, uit.AM)
}
