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
