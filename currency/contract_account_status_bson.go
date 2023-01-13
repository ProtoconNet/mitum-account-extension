package currency // nolint: dupl, revive

import (
	"github.com/spikeekips/mitum/util"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (cs ContractAccount) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(cs.Hint()),
		bson.M{
			"isactive": cs.isActive,
			"owner":    cs.owner,
		}),
	)
}

type ContractAccountBSONUnmarshaler struct {
	HT hint.Hint `json:"_hint"`
	IA bool      `bson:"isactive"`
	OW string    `bson:"owner"`
}

func (cs *ContractAccount) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of ContractAccount")

	var ucs ContractAccountBSONUnmarshaler
	if err := bsonenc.Unmarshal(b, &ucs); err != nil {
		return e(err, "")
	}

	return cs.unpack(enc, ucs.HT, ucs.IA, ucs.OW)
}
