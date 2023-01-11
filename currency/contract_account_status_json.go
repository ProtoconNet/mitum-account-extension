package currency

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type ContractAccountJSONMarshaler struct {
	hint.BaseHinter
	IA bool         `json:"isactive"`
	OW base.Address `json:"owner"`
}

func (cs ContractAccount) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ContractAccountJSONMarshaler{
		BaseHinter: cs.BaseHinter,
		IA:         cs.isActive,
		OW:         cs.owner,
	})
}

type ContractAccountJSONUnmarshaler struct {
	HT hint.Hint `json:"_hint"`
	IA bool      `json:"isactive"`
	OW string    `json:"owner"`
}

func (ca *ContractAccount) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ContractAccount")

	var ucs ContractAccountJSONUnmarshaler
	if err := enc.Unmarshal(b, &ucs); err != nil {
		return e(err, "")
	}

	return ca.unpack(enc, ucs.HT, ucs.IA, ucs.OW)
}
