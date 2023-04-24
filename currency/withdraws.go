package currency

import (
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	WithdrawsFactHint = hint.MustNewHint("mitum-currency-contract-account-withdraws-operation-fact-v0.0.1")
	WithdrawsHint     = hint.MustNewHint("mitum-currency-contract-account-withdraws-operation-v0.0.1")
)

var MaxWithdrawsItems uint = 10

type WithdrawsItem interface {
	hint.Hinter
	util.IsValider
	mitumcurrency.AmountsItem
	Bytes() []byte
	Target() base.Address
	Rebuild() WithdrawsItem
}

type WithdrawsFact struct {
	base.BaseFact
	sender base.Address
	items  []WithdrawsItem
}

func NewWithdrawsFact(token []byte, sender base.Address, items []WithdrawsItem) WithdrawsFact {
	bf := base.NewBaseFact(WithdrawsFactHint, token)
	fact := WithdrawsFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact WithdrawsFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact WithdrawsFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact WithdrawsFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact WithdrawsFact) Bytes() []byte {
	its := make([][]byte, len(fact.items))
	for i := range fact.items {
		its[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(its...),
	)
}

func (fact WithdrawsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := mitumcurrency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	} else if n > int(MaxWithdrawsItems) {
		return util.ErrInvalid.Errorf("items, %d over max, %d", n, MaxWithdrawsItems)
	}

	if err := util.CheckIsValiders(nil, false, fact.sender); err != nil {
		return err
	}

	foundTargets := map[string]struct{}{}
	for i := range fact.items {
		it := fact.items[i]
		if err := util.CheckIsValiders(nil, false, it); err != nil {
			return err
		}

		k := it.Target().String()
		switch _, found := foundTargets[k]; {
		case found:
			return util.ErrInvalid.Errorf("duplicate target found, %s", it.Target())
		case fact.sender.Equal(it.Target()):
			return util.ErrInvalid.Errorf("target is same with sender, %q", fact.sender)
		default:
			foundTargets[k] = struct{}{}
		}
	}

	return nil
}

func (fact WithdrawsFact) Sender() base.Address {
	return fact.sender
}

func (fact WithdrawsFact) Items() []WithdrawsItem {
	return fact.items
}

func (fact WithdrawsFact) Rebuild() WithdrawsFact {
	items := make([]WithdrawsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact WithdrawsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)+1)
	for i := range fact.items {
		as[i] = fact.items[i].Target()
	}

	as[len(fact.items)] = fact.Sender()

	return as, nil
}

type Withdraws struct {
	mitumcurrency.BaseOperation
}

func NewWithdraws(fact WithdrawsFact) (Withdraws, error) {
	return Withdraws{BaseOperation: mitumcurrency.NewBaseOperation(WithdrawsHint, fact)}, nil
}

func (op *Withdraws) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
