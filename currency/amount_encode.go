package currency // nolint: dupl, revive

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (am *AmountValue) unpack(
	enc encoder.Encoder,
	bam []byte,
	sci string,
) error {
	h, err := enc.Decode(bam)
	if err != nil {
		return err
	}
	v, ok := h.(currency.Amount)
	if !ok {
		return util.WrongTypeError.Errorf("expected Amount, not %T", h)
	}
	am.amount = v
	am.contractid = ContractID(sci)

	return nil
}
