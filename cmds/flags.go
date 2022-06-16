package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

type AddressFlag struct {
	s string
}

func (v *AddressFlag) UnmarshalText(b []byte) error {
	v.s = string(b)

	return nil
}

func (v *AddressFlag) String() string {
	return v.s
}

func (v *AddressFlag) Encode(enc encoder.Encoder) (base.Address, error) {
	return base.DecodeAddressFromString(v.s, enc)
}

type ContractIDFlag struct {
	ID extensioncurrency.ContractID
}

func (v *ContractIDFlag) UnmarshalText(b []byte) error {
	cid := extensioncurrency.ContractID(string(b))
	if err := cid.IsValid(nil); err != nil {
		return err
	}
	v.ID = cid

	return nil
}

func (v *ContractIDFlag) String() string {
	return v.ID.String()
}
