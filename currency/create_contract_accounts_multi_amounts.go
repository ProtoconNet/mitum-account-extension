package currency

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

var maxCurrenciesCreateContractAccountsItemMultiAmounts = 10

var (
	CreateContractAccountsItemMultiAmountsHint = hint.MustNewHint("mitum-currency-create-contract-accounts-multiple-amounts-v0.0.1")
)

type CreateContractAccountsItemMultiAmounts struct {
	BaseCreateContractAccountsItem
}

func NewCreateContractAccountsItemMultiAmounts(keys currency.AccountKeys, amounts []currency.Amount) CreateContractAccountsItemMultiAmounts {
	return CreateContractAccountsItemMultiAmounts{
		BaseCreateContractAccountsItem: NewBaseCreateContractAccountsItem(CreateContractAccountsItemMultiAmountsHint, keys, amounts),
	}
}

func (it CreateContractAccountsItemMultiAmounts) IsValid([]byte) error {
	if err := it.BaseCreateContractAccountsItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.amounts); n > maxCurrenciesCreateContractAccountsItemMultiAmounts {
		return util.ErrInvalid.Errorf("amounts over allowed; %d > %d", n, maxCurrenciesCreateContractAccountsItemMultiAmounts)
	}

	return nil
}

func (it CreateContractAccountsItemMultiAmounts) Rebuild() CreateContractAccountsItem {
	it.BaseCreateContractAccountsItem = it.BaseCreateContractAccountsItem.Rebuild().(BaseCreateContractAccountsItem)

	return it
}
