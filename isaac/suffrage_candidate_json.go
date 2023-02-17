package isaacoperation

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type suffrageCandidateFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Address   base.Address   `json:"address"`
	Publickey base.Publickey `json:"publickey"`
}

func (fact SuffrageCandidateFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(suffrageCandidateFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Address:               fact.address,
		Publickey:             fact.publickey,
	})
}

type suffrageCandidateFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Address   string `json:"address"`
	Publickey string `json:"publickey"`
}

func (fact *SuffrageCandidateFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of SuffrageCandidateFact")

	var uf suffrageCandidateFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Address, uf.Publickey)
}
