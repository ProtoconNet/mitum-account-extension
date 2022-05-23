package example

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	ConfigContractAccountFactType   = hint.Type("mitum-currency-config-contract-account-operation-fact")
	ConfigContractAccountFactHint   = hint.NewHint(ConfigContractAccountFactType, "v0.0.1")
	ConfigContractAccountFactHinter = ConfigContractAccountFact{BaseHinter: hint.NewBaseHinter(ConfigContractAccountFactHint)}
	ConfigContractAccountType       = hint.Type("mitum-currency-config-contract-account-operation")
	ConfigContractAccountHint       = hint.NewHint(ConfigContractAccountType, "v0.0.1")
	ConfigContractAccountHinter     = ConfigContractAccount{BaseOperation: operationHinter(ConfigContractAccountHint)}
)

type ConfigContractAccountFact struct {
	hint.BaseHinter
	h        valuehash.Hash
	token    []byte
	sender   base.Address
	target   base.Address
	currency currency.CurrencyID
}

func NewConfigContractAccountFact(token []byte, sender, target base.Address, currency currency.CurrencyID) ConfigContractAccountFact {
	fact := ConfigContractAccountFact{
		BaseHinter: hint.NewBaseHinter(ConfigContractAccountFactHint),
		token:      token,
		sender:     sender,
		target:     target,
		currency:   currency,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact ConfigContractAccountFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact ConfigContractAccountFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ConfigContractAccountFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		fact.target.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact ConfigContractAccountFact) IsValid(b []byte) error {
	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return isvalid.Check(nil, false,
		fact.sender,
		fact.target,
		fact.currency,
	)
}

func (fact ConfigContractAccountFact) Token() []byte {
	return fact.token
}

func (fact ConfigContractAccountFact) Sender() base.Address {
	return fact.sender
}

func (fact ConfigContractAccountFact) Target() base.Address {
	return fact.target
}

func (fact ConfigContractAccountFact) Currency() currency.CurrencyID {
	return fact.currency
}

func (fact ConfigContractAccountFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.sender, fact.target}, nil
}

type ConfigContractAccount struct {
	currency.BaseOperation
}

func NewConfigContractAccount(fact ConfigContractAccountFact, fs []base.FactSign, memo string) (ConfigContractAccount, error) {
	bo, err := currency.NewBaseOperationFromFact(ConfigContractAccountHint, fact, fs, memo)
	if err != nil {
		return ConfigContractAccount{}, err
	}

	return ConfigContractAccount{BaseOperation: bo}, nil
}
