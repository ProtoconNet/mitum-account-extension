package currency

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	CreateContractAccountsItemSingleAmountHint = hint.MustNewHint("mitum-currency-create-contract-accounts-single-amount-v0.0.1")
)

type CreateContractAccountsItemSingleAmount struct {
	BaseCreateContractAccountsItem
}

func NewCreateContractAccountsItemSingleAmount(keys currency.AccountKeys, amount currency.Amount) CreateContractAccountsItemSingleAmount {
	return CreateContractAccountsItemSingleAmount{
		BaseCreateContractAccountsItem: NewBaseCreateContractAccountsItem(CreateContractAccountsItemSingleAmountHint, keys, []currency.Amount{amount}),
	}
}

func (it CreateContractAccountsItemSingleAmount) IsValid([]byte) error {
	if err := it.BaseCreateContractAccountsItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.amounts); n != 1 {
		return util.ErrInvalid.Errorf("only one amount allowed; %d", n)
	}

	return nil
}

func (it CreateContractAccountsItemSingleAmount) Rebuild() CreateContractAccountsItem {
	it.BaseCreateContractAccountsItem = it.BaseCreateContractAccountsItem.Rebuild().(BaseCreateContractAccountsItem)

	return it
}
