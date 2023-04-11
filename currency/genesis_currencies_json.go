package currency

import (
	"encoding/json"

	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type GenesisCurrenciesFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	GenesisNodeKey base.Publickey            `json:"genesis_node_key"`
	Keys           mitumcurrency.AccountKeys `json:"keys"`
	Currencies     []CurrencyDesign          `json:"currencies"`
}

func (fact GenesisCurrenciesFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(GenesisCurrenciesFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		GenesisNodeKey:        fact.genesisNodeKey,
		Keys:                  fact.keys,
		Currencies:            fact.cs,
	})
}

type GenesisCurrenciesFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	GenesisNodeKey string          `json:"genesis_node_key"`
	Keys           json.RawMessage `json:"keys"`
	Currencies     json.RawMessage `json:"currencies"`
}

func (fact *GenesisCurrenciesFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of GenesisCurrenciesFact")

	var uf GenesisCurrenciesFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc, uf.GenesisNodeKey, uf.Keys, uf.Currencies)
}

type genesisCurrenciesMarshaler struct {
	mitumcurrency.BaseOperationJSONMarshaler
}

func (op GenesisCurrencies) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(genesisCurrenciesMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *GenesisCurrencies) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of GenesisCurrencies")

	var ubo mitumcurrency.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
