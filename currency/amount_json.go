package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type AmountValueJSONPacker struct {
	jsonenc.HintedHead
	AM currency.Amount `json:"amount"`
	CI ContractID      `json:"contractid"`
}

func (am AmountValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(AmountValueJSONPacker{
		HintedHead: jsonenc.NewHintedHead(am.Hint()),
		AM:         am.amount,
		CI:         am.contractid,
	})
}

type AmountValueJSONUnpacker struct {
	HT hint.Hint       `json:"_hint"`
	AM json.RawMessage `json:"amount"`
	CI string          `json:"contractid"`
}

func (am *AmountValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uam AmountValueJSONUnpacker
	if err := enc.Unmarshal(b, &uam); err != nil {
		return err
	}

	return am.unpack(enc, uam.AM, uam.CI)
}
