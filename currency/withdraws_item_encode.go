package currency

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (it *BaseWithdrawsItem) unpack(enc encoder.Encoder, ht hint.Hint, tg string, bam []byte) error {
	e := util.StringErrorFunc("failed to unmarshal BaseWithdrawsItem")

	it.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(tg, enc); {
	case err != nil:
		return e(err, "")
	default:
		it.target = a
	}

	ham, err := enc.DecodeSlice(bam)
	if err != nil {
		return e(err, "")
	}

	amounts := make([]currency.Amount, len(ham))
	for i := range ham {
		j, ok := ham[i].(currency.Amount)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected Amount, not %T", ham[i]), "")
		}

		amounts[i] = j
	}

	it.amounts = amounts

	return nil
}
