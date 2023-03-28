package isaacoperation

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type suffrageJoinFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Candidate base.Address `json:"candidate"`
	Start     base.Height  `json:"start_height"`
}

func (fact SuffrageJoinFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(suffrageJoinFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Candidate:             fact.candidate,
		Start:                 fact.start,
	})
}

type suffrageJoinFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Candidate string      `json:"candidate"`
	Start     base.Height `json:"start_height"`
}

func (fact *SuffrageJoinFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode SuffrageJoinFact")

	var uf suffrageJoinFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.Candidate, uf.Start)
}

type suffrageGenesisJoinFactJSONMarshaler struct {
	Nodes []base.Node `json:"nodes"`
	base.BaseFactJSONMarshaler
}

func (fact SuffrageGenesisJoinFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(suffrageGenesisJoinFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Nodes:                 fact.nodes,
	})
}

type suffrageGenesisJoinFactJSONUnmarshaler struct {
	Nodes []json.RawMessage `json:"nodes"`
	base.BaseFactJSONUnmarshaler
}

func (fact *SuffrageGenesisJoinFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode SuffrageGenesisJoinFact")

	var uf suffrageGenesisJoinFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	fact.nodes = make([]base.Node, len(uf.Nodes))

	for i := range uf.Nodes {
		if err := encoder.Decode(enc, uf.Nodes[i], &fact.nodes[i]); err != nil {
			return e(err, "")
		}
	}

	return nil
}
