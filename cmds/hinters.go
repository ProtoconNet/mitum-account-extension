package cmds

import (
	"github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	mitumcurrency "github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum-currency/digest"
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/util/encoder"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: mitumcurrency.AccountHint, Instance: mitumcurrency.Account{}},
	{Hint: mitumcurrency.AddressHint, Instance: mitumcurrency.Address{}},
	{Hint: mitumcurrency.AmountHint, Instance: mitumcurrency.Amount{}},
	{Hint: mitumcurrency.CreateAccountsItemMultiAmountsHint, Instance: mitumcurrency.CreateAccountsItemMultiAmounts{}},
	{Hint: mitumcurrency.CreateAccountsItemSingleAmountHint, Instance: mitumcurrency.CreateAccountsItemSingleAmount{}},
	{Hint: mitumcurrency.CreateAccountsHint, Instance: mitumcurrency.CreateAccounts{}},
	{Hint: mitumcurrency.KeyUpdaterHint, Instance: mitumcurrency.KeyUpdater{}},
	{Hint: mitumcurrency.TransfersItemMultiAmountsHint, Instance: mitumcurrency.TransfersItemMultiAmounts{}},
	{Hint: mitumcurrency.TransfersItemSingleAmountHint, Instance: mitumcurrency.TransfersItemSingleAmount{}},
	{Hint: mitumcurrency.TransfersHint, Instance: mitumcurrency.Transfers{}},
	{Hint: currency.CurrencyDesignHint, Instance: currency.CurrencyDesign{}},
	{Hint: currency.CurrencyPolicyHint, Instance: currency.CurrencyPolicy{}},
	// {Hint: mitumcurrency.CurrencyRegisterHint, Instance: mitumcurrency.CurrencyRegister{}},
	// {Hint: mitumcurrency.CurrencyPolicyUpdaterHint, Instance: mitumcurrency.CurrencyPolicyUpdater{}},
	// {Hint: mitumcurrency.SuffrageInflationHint, Instance: mitumcurrency.SuffrageInflation{}},
	{Hint: currency.ContractAccountKeysHint, Instance: currency.ContractAccountKeys{}},
	{Hint: currency.CreateContractAccountsItemMultiAmountsHint, Instance: currency.CreateContractAccountsItemMultiAmounts{}},
	{Hint: currency.CreateContractAccountsItemSingleAmountHint, Instance: currency.CreateContractAccountsItemSingleAmount{}},
	{Hint: currency.CreateContractAccountsHint, Instance: currency.CreateContractAccounts{}},
	{Hint: currency.WithdrawsItemMultiAmountsHint, Instance: currency.WithdrawsItemMultiAmounts{}},
	{Hint: currency.WithdrawsItemSingleAmountHint, Instance: currency.WithdrawsItemSingleAmount{}},
	{Hint: currency.WithdrawsHint, Instance: currency.Withdraws{}},
	// {Hint: mitumcurrency.FeeOperationFactHint, Instance: mitumcurrency.FeeOperationFact{}},
	// {Hint: mitumcurrency.FeeOperationHint, Instance: mitumcurrency.FeeOperation{}},
	{Hint: currency.GenesisCurrenciesFactHint, Instance: currency.GenesisCurrenciesFact{}},
	{Hint: currency.GenesisCurrenciesHint, Instance: currency.GenesisCurrencies{}},
	{Hint: mitumcurrency.AccountKeysHint, Instance: mitumcurrency.BaseAccountKeys{}},
	{Hint: mitumcurrency.AccountKeyHint, Instance: mitumcurrency.BaseAccountKey{}},
	{Hint: currency.NilFeeerHint, Instance: currency.NilFeeer{}},
	{Hint: currency.FixedFeeerHint, Instance: currency.FixedFeeer{}},
	{Hint: currency.RatioFeeerHint, Instance: currency.RatioFeeer{}},
	{Hint: mitumcurrency.AccountStateValueHint, Instance: mitumcurrency.AccountStateValue{}},
	{Hint: mitumcurrency.BalanceStateValueHint, Instance: mitumcurrency.BalanceStateValue{}},
	{Hint: currency.ContractAccountStateValueHint, Instance: currency.ContractAccountStateValue{}},
	{Hint: currency.CurrencyDesignStateValueHint, Instance: currency.CurrencyDesignStateValue{}},
	{Hint: digest.AccountValueHint, Instance: digest.AccountValue{}},
	{Hint: digest.OperationValueHint, Instance: digest.OperationValue{}},
}

var supportedProposalOperationFactHinters = []encoder.DecodeDetail{
	{Hint: mitumcurrency.CreateAccountsFactHint, Instance: mitumcurrency.CreateAccountsFact{}},
	{Hint: mitumcurrency.KeyUpdaterFactHint, Instance: mitumcurrency.KeyUpdaterFact{}},
	{Hint: mitumcurrency.TransfersFactHint, Instance: mitumcurrency.TransfersFact{}},
	// {Hint: mitumcurrency.CurrencyRegisterFactHint, Instance: mitumcurrency.CurrencyRegisterFact{}},
	// {Hint: mitumcurrency.CurrencyPolicyUpdaterFactHint, Instance: mitumcurrency.CurrencyPolicyUpdaterFact{}},
	// {Hint: mitumcurrency.SuffrageInflationFactHint, Instance: mitumcurrency.SuffrageInflationFact{}},
	{Hint: currency.CreateContractAccountsFactHint, Instance: currency.CreateContractAccountsFact{}},
	{Hint: currency.WithdrawsFactHint, Instance: currency.WithdrawsFact{}},
}

func init() {
	Hinters = make([]encoder.DecodeDetail, len(launch.Hinters)+len(hinters))
	copy(Hinters, launch.Hinters)
	copy(Hinters[len(launch.Hinters):], hinters)

	SupportedProposalOperationFactHinters = make([]encoder.DecodeDetail, len(launch.SupportedProposalOperationFactHinters)+len(supportedProposalOperationFactHinters))
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[len(launch.SupportedProposalOperationFactHinters):], supportedProposalOperationFactHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for i := range Hinters {
		if err := enc.Add(Hinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := enc.Add(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
