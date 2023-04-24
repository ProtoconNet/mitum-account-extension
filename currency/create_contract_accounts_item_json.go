package currency

import (
	"encoding/json"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CreateContractAccountsItemJSONMarshaler struct {
	hint.BaseHinter
	Keys     mitumcurrency.AccountKeys `json:"keys"`
	Amounts  []mitumcurrency.Amount    `json:"amounts"`
	AddrType hint.Type                 `json:"addrtype"`
}

func (it BaseCreateContractAccountsItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateContractAccountsItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Keys:       it.keys,
		Amounts:    it.amounts,
		AddrType:   it.addressType,
	})
}

type CreateContractAccountsItemJSONUnMarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Keys     json.RawMessage `json:"keys"`
	Amounts  json.RawMessage `json:"amounts"`
	AddrType string          `json:"addrtype"`
}

func (it *BaseCreateContractAccountsItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of BaseCreateContractAccountsItem")

	var uit CreateContractAccountsItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return e(err, "")
	}

	return it.unpack(enc, uit.Hint, uit.Keys, uit.Amounts, uit.AddrType)
}
