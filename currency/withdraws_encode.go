package currency

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (fact *WithdrawsFact) unpack(enc encoder.Encoder, sd string, bit []byte) error {
	e := util.StringErrorFunc("failed to unmarshal WithdrawsFact")

	switch a, err := base.DecodeAddress(sd, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.sender = a
	}

	hit, err := enc.DecodeSlice(bit)
	if err != nil {
		return e(err, "")
	}

	items := make([]WithdrawsItem, len(hit))
	for i := range hit {
		j, ok := hit[i].(WithdrawsItem)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected WithdrawsItem, not %T", hit[i]), "")
		}

		items[i] = j
	}
	fact.items = items

	return nil
}
