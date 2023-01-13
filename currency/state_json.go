package currency

import (
	"encoding/json"

	"github.com/spikeekips/mitum/util"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type CurrencyDesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	CD CurrencyDesign `json:"currencydesign"`
}

func (s CurrencyDesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CurrencyDesignStateValueJSONMarshaler{
		BaseHinter: s.BaseHinter,
		CD:         s.CurrencyDesign,
	})
}

type CurrencyDesignStateValueJSONUnmarshaler struct {
	CD json.RawMessage `json:"currencydesign"`
}

func (s *CurrencyDesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode CurrencyDesignStateValue")

	var u CurrencyDesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	var cd CurrencyDesign

	if err := cd.DecodeJSON(u.CD, enc); err != nil {
		return e(err, "")
	}

	s.CurrencyDesign = cd

	return nil
}

type ContractAccountStateValueJSONMarshaler struct {
	hint.BaseHinter
	CA ContractAccount `json:"contractaccount"`
}

func (s ContractAccountStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ContractAccountStateValueJSONMarshaler{
		BaseHinter: s.BaseHinter,
		CA:         s.account,
	})
}

type ContractAccountStateValueJSONUnmarshaler struct {
	CA json.RawMessage `json:"contractaccount"`
}

func (s *ContractAccountStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode ContractAccountStateValue")

	var u ContractAccountStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	var ca ContractAccount

	if err := ca.DecodeJSON(u.CA, enc); err != nil {
		return e(err, "")
	}

	s.account = ca

	return nil
}
