package currency // nolint: dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/hint"
)

func (fact CreateContractAccountsFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(
			bsonenc.NewHintedDoc(fact.Hint()),
			bson.M{
				"sender": fact.sender,
				"items":  fact.items,
			},
			fact.BaseFact.BSONM(),
		))
}

type CreateContractAccountsFactBSONUnmarshaler struct {
	HT hint.Hint `bson:"_hint"`
	SD string    `bson:"sender"`
	IT bson.Raw  `bson:"items"`
}

func (fact *CreateContractAccountsFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CreateContractAccountsFact")

	var ubf base.BaseFact
	if err := ubf.DecodeBSON(b, enc); err != nil {
		return err
	}

	fact.BaseFact = ubf

	var uf CreateContractAccountsFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseHinter = hint.NewBaseHinter(uf.HT)

	return fact.unpack(enc, uf.SD, uf.IT)
}

func (op *CreateContractAccounts) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
