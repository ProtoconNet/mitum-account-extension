package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type CurrencyDesignJSONMarshaler struct {
	hint.BaseHinter
	Amount         currency.Amount `json:"amount"`
	GenesisAccount base.Address    `json:"genesis_account"`
	Policy         CurrencyPolicy  `json:"policy"`
	Aggregate      string          `json:"aggregate"`
}

func (de CurrencyDesign) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CurrencyDesignJSONMarshaler{
		BaseHinter:     de.BaseHinter,
		Amount:         de.amount,
		GenesisAccount: de.genesisAccount,
		Policy:         de.policy,
		Aggregate:      de.aggregate.String(),
	})
}

type CurrencyDesignJSONUnmarshaler struct {
	Hint           hint.Hint       `json:"_hint"`
	Amount         json.RawMessage `json:"amount"`
	GenesisAccount string          `json:"genesis_account"`
	Policy         json.RawMessage `json:"policy"`
	Aggregate      string          `json:"aggregate"`
}

func (de *CurrencyDesign) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CurrencyDesign")

	var ude CurrencyDesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &ude); err != nil {
		return e(err, "")
	}

	return de.unpack(enc, ude.Hint, ude.Amount, ude.GenesisAccount, ude.Policy, ude.Aggregate)
}
