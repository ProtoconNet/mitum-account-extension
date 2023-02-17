package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type WithdrawsItemJSONMarshaler struct {
	hint.BaseHinter
	Target  base.Address      `json:"target"`
	Amounts []currency.Amount `json:"amounts"`
}

func (it BaseWithdrawsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(WithdrawsItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Target:     it.target,
		Amounts:    it.amounts,
	})
}

type BaseWithdrawsItemJSONUnpacker struct {
	Hint    hint.Hint       `json:"_hint"`
	Target  string          `json:"target"`
	Amounts json.RawMessage `json:"amounts"`
}

func (it *BaseWithdrawsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of BaseWithdrawsItem")

	var uit BaseWithdrawsItemJSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Target, uit.Amounts)
}
