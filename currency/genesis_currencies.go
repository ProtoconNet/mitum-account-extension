package currency

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	GenesisCurrenciesFactHint = hint.MustNewHint("mitum-currency-genesis-currencies-operation-fact-v0.0.1")
	GenesisCurrenciesHint     = hint.MustNewHint("mitum-currency-genesis-currencies-operation-v0.0.1")
)

type GenesisCurrenciesFact struct {
	base.BaseFact
	genesisNodeKey base.Publickey
	keys           currency.AccountKeys
	cs             []CurrencyDesign
}

func NewGenesisCurrenciesFact(
	token []byte,
	genesisNodeKey base.Publickey,
	keys currency.AccountKeys,
	cs []CurrencyDesign,
) GenesisCurrenciesFact {
	fact := GenesisCurrenciesFact{
		BaseFact:       base.NewBaseFact(GenesisCurrenciesFactHint, token),
		genesisNodeKey: genesisNodeKey,
		keys:           keys,
		cs:             cs,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact GenesisCurrenciesFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact GenesisCurrenciesFact) Bytes() []byte {
	bs := make([][]byte, len(fact.cs)+3)
	bs[0] = fact.Token()
	bs[1] = []byte(fact.genesisNodeKey.String())
	bs[2] = fact.keys.Bytes()

	for i := range fact.cs {
		bs[i+3] = fact.cs[i].Bytes()
	}

	return util.ConcatBytesSlice(bs...)
}

func (fact GenesisCurrenciesFact) IsValid(b []byte) error {
	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if len(fact.cs) < 1 {
		return util.ErrInvalid.Errorf("empty GenesisCurrency for GenesisCurrenciesFact")
	}

	if err := util.CheckIsValiders(nil, false, fact.genesisNodeKey, fact.keys); err != nil {
		return util.ErrInvalid.Errorf("invalid fact: %w", err)
	}

	founds := map[currency.CurrencyID]struct{}{}
	for i := range fact.cs {
		c := fact.cs[i]
		if err := c.IsValid(nil); err != nil {
			return err
		} else if _, found := founds[c.Currency()]; found {
			return util.ErrInvalid.Errorf("duplicate currency id found, %q", c.Currency())
		} else {
			founds[c.Currency()] = struct{}{}
		}
	}

	return nil
}

func (fact GenesisCurrenciesFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact GenesisCurrenciesFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact GenesisCurrenciesFact) GenesisNodeKey() base.Publickey {
	return fact.genesisNodeKey
}

func (fact GenesisCurrenciesFact) Keys() currency.AccountKeys {
	return fact.keys
}

func (fact GenesisCurrenciesFact) Address() (base.Address, error) {
	return currency.NewAddressFromKeys(fact.keys)
}

func (fact GenesisCurrenciesFact) Currencies() []CurrencyDesign {
	return fact.cs
}

type GenesisCurrencies struct {
	currency.BaseOperation
}

func NewGenesisCurrencies(fact GenesisCurrenciesFact) GenesisCurrencies {
	return GenesisCurrencies{BaseOperation: currency.NewBaseOperation(GenesisCurrenciesHint, fact)}
}

func (op GenesisCurrencies) IsValid(networkID []byte) error {
	if err := op.BaseOperation.IsValid(networkID); err != nil {
		return err
	}

	if len(op.Signs()) != 1 {
		return util.ErrInvalid.Errorf("genesis currencies should be signed only by genesis node key")
	}

	fact, ok := op.Fact().(GenesisCurrenciesFact)
	if !ok {
		return errors.Errorf("expected GenesisCurrenciesFact, not %T", op.Fact())
	}

	if !fact.genesisNodeKey.Equal(op.Signs()[0].Signer()) {
		return util.ErrInvalid.Errorf("not signed by genesis node key")
	}

	return nil
}
