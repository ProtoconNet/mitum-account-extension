package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type CreateContractAccountsFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	OW base.Address                 `json:"sender"`
	IT []CreateContractAccountsItem `json:"items"`
}

func (fact CreateContractAccountsFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateContractAccountsFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		OW:                    fact.sender,
		IT:                    fact.items,
	})
}

type CreateContractAccountsFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	OW string          `json:"sender"`
	IT json.RawMessage `json:"items"`
}

func (fact *CreateContractAccountsFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CreateContractAccountsFact")

	var uf CreateContractAccountsFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.OW, uf.IT)
}

func (op *CreateContractAccounts) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
