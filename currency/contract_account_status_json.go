package currency

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type ContractAccountJSONMarshaler struct {
	hint.BaseHinter
	IsActive bool         `json:"isactive"`
	Owner    base.Address `json:"owner"`
}

func (cs ContractAccount) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ContractAccountJSONMarshaler{
		BaseHinter: cs.BaseHinter,
		IsActive:   cs.isActive,
		Owner:      cs.owner,
	})
}

type ContractAccountJSONUnmarshaler struct {
	Hint     hint.Hint `json:"_hint"`
	IsActive bool      `json:"isactive"`
	Owner    string    `json:"owner"`
}

func (ca *ContractAccount) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ContractAccount")

	var ucs ContractAccountJSONUnmarshaler
	if err := enc.Unmarshal(b, &ucs); err != nil {
		return e(err, "")
	}

	return ca.unpack(enc, ucs.Hint, ucs.IsActive, ucs.Owner)
}
