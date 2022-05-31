package currency

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (de *CurrencyDesign) unpack(enc encoder.Encoder, am currency.Amount, ga base.AddressDecoder, bpo []byte, ag currency.Big) error {
	de.Amount = am

	a, err := ga.Encode(enc)
	if err != nil {
		return err
	}
	de.genesisAccount = a

	if err := encoder.Decode(bpo, enc, &de.policy); err != nil {
		return err
	}

	de.aggregate = ag

	return nil
}
