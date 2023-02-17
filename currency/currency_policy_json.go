package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type CurrencyPolicyJSONMarshaler struct {
	hint.BaseHinter
	MinBalance string `json:"new_account_min_balance"`
	Feeer      Feeer  `json:"feeer"`
}

func (po CurrencyPolicy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CurrencyPolicyJSONMarshaler{
		BaseHinter: po.BaseHinter,
		MinBalance: po.newAccountMinBalance.String(),
		Feeer:      po.feeer,
	})
}

type CurrencyPolicyJSONUnmarshaler struct {
	Hint       hint.Hint       `json:"_hint"`
	MinBalance string          `json:"new_account_min_balance"`
	Feeer      json.RawMessage `json:"feeer"`
}

func (po *CurrencyPolicy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CurrencyPolicy")

	var upo CurrencyPolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e(err, "")
	}

	return po.unpack(enc, upo.Hint, upo.MinBalance, upo.Feeer)
}
