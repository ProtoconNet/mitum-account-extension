package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type TransferFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	SD base.Address    `json:"sender"`
	IT []WithdrawsItem `json:"items"`
}

func (fact WithdrawsFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		SD:                    fact.sender,
		IT:                    fact.items,
	})
}

type WithdrawsFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	SD string          `json:"sender"`
	IT json.RawMessage `json:"items"`
}

func (fact *WithdrawsFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of WithdrawsFact")

	var uf WithdrawsFactJSONUnmarshaler

	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.SD, uf.IT)
}
