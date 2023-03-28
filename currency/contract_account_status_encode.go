package currency // nolint: dupl, revive

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
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
