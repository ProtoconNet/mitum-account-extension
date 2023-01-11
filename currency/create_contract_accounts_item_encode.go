package currency

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (it *BaseCreateContractAccountsItem) unpack(enc encoder.Encoder, ht hint.Hint, bks []byte, bam []byte) error {
	e := util.StringErrorFunc("failed to unmarshal BaseCreateContractAccountsItem")

	it.BaseHinter = hint.NewBaseHinter(ht)

	if hinter, err := enc.Decode(bks); err != nil {
		return e(err, "")
	} else if k, ok := hinter.(currency.AccountKeys); !ok {
		return util.ErrWrongType.Errorf("expected AccountsKeys, not %T", hinter)
	} else {
		it.keys = k
	}

	ham, err := enc.DecodeSlice(bam)
	if err != nil {
		return e(err, "")
	}

	amounts := make([]currency.Amount, len(ham))
	for i := range ham {
		j, ok := ham[i].(currency.Amount)
		if !ok {
			return util.ErrWrongType.Errorf("expected Amount, not %T", ham[i])
		}

		amounts[i] = j
	}

	it.amounts = amounts

	return nil
}
