package currency // nolint: dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum-currency/currency"
	bsonenc "github.com/spikeekips/mitum-currency/digest/util/bson"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact WithdrawsFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  fact.Hint().String(),
			"sender": fact.sender,
			"items":  fact.items,
			"hash":   fact.BaseFact.Hash().String(),
			"token":  fact.BaseFact.Token(),
		},
	)
}

type WithdrawsFactBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Sender string   `bson:"sender"`
	Items  bson.Raw `bson:"items"`
}

func (fact *WithdrawsFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of WithdrawsFact")

	var ubf currency.BaseFactBSONUnmarshaler
	if err := enc.Unmarshal(b, &ubf); err != nil {
		return err
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf WithdrawsFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}

	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc, uf.Sender, uf.Items)
}

func (op Withdraws) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(op.BaseOperation)
}

func (op *Withdraws) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Withdraw")

	var ubo currency.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
