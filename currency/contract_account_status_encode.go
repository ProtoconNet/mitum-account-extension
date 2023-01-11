package currency // nolint: dupl, revive

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (cs *ContractAccount) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	ia bool,
	ow string,
) error {
	e := util.StringErrorFunc("failed to unmarshal ContractAccount")

	cs.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(ow, enc); {
	case err != nil:
		return e(err, "failed to decode address")
	default:
		cs.owner = a
	}

	cs.isActive = ia

	return nil
}
