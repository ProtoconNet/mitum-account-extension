package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type WithdrawsItemJSONPacker struct {
	hint.BaseHinter
	TG base.Address      `json:"target"`
	AM []currency.Amount `json:"amounts"`
}

func (it BaseWithdrawsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(WithdrawsItemJSONPacker{
		BaseHinter: it.BaseHinter,
		TG:         it.target,
		AM:         it.amounts,
	})
}

type BaseWithdrawsItemJSONUnpacker struct {
	HT hint.Hint       `json:"_hint"`
	TG string          `json:"target"`
	AM json.RawMessage `json:"amounts"`
}

func (it *BaseWithdrawsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of BaseWithdrawsItem")

	var uit BaseWithdrawsItemJSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.HT, uit.TG, uit.AM)
}
