package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type CurrencyDesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	CurrencyDesign CurrencyDesign `json:"currencydesign"`
}

func (s CurrencyDesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		CurrencyDesignStateValueJSONMarshaler(s),
	)
}

type CurrencyDesignStateValueJSONUnmarshaler struct {
	Hint           hint.Hint       `json:"_hint"`
	CurrencyDesign json.RawMessage `json:"currencydesign"`
}

func (s *CurrencyDesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CurrencyDesignStateValue")

	var u CurrencyDesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var cd CurrencyDesign
	if err := cd.DecodeJSON(u.CurrencyDesign, enc); err != nil {
		return e(err, "")
	}
	s.CurrencyDesign = cd

	return nil
}

type ContractAccountStateValueJSONMarshaler struct {
	hint.BaseHinter
	ContractAccount ContractAccount `json:"contractaccount"`
}

func (s ContractAccountStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ContractAccountStateValueJSONMarshaler{
		BaseHinter:      s.BaseHinter,
		ContractAccount: s.account,
	})
}

type ContractAccountStateValueJSONUnmarshaler struct {
	Hint            hint.Hint       `json:"_hint"`
	ContractAccount json.RawMessage `json:"contractaccount"`
}

func (s *ContractAccountStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ContractAccountStateValue")

	var u ContractAccountStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var ca ContractAccount
	if err := ca.DecodeJSON(u.ContractAccount, enc); err != nil {
		return e(err, "")
	}
	s.account = ca

	return nil
}
