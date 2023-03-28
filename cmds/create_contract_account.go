package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
)

type CreateContractAccountCommand struct {
	baseCommand
	OperationFlags
	Sender    AddressFlag          `arg:"" name:"sender" help:"sender address" required:"true"`
	Threshold uint                 `help:"threshold for keys (default: ${create_contract_account_threshold})" default:"${create_contract_account_threshold}"` // nolint
	Keys      []KeyFlag            `name:"key" help:"key for new account (ex: \"<public key>,<weight>\")" sep:"@"`
	Amounts   []CurrencyAmountFlag `arg:"" name:"currency-amount" help:"amount (ex: \"<currency>,<amount>\")"`
	sender    base.Address
	keys      mitumcurrency.BaseAccountKeys
}

func NewCreateContractAccountCommand() CreateContractAccountCommand {
	cmd := NewbaseCommand()
	return CreateContractAccountCommand{
		baseCommand: *cmd,
	}
}

func (cmd *CreateContractAccountCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.encs
	enc = cmd.enc

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *CreateContractAccountCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = a

	if len(cmd.Keys) < 1 {
		return errors.Errorf("--key must be given at least one")
	}

	if len(cmd.Amounts) < 1 {
		return errors.Errorf("empty currency-amount, must be given at least one")
	}

	{
		ks := make([]mitumcurrency.AccountKey, len(cmd.Keys))
		for i := range cmd.Keys {
			ks[i] = cmd.Keys[i].Key
		}

		if kys, err := mitumcurrency.NewBaseAccountKeys(ks, cmd.Threshold); err != nil {
			return err
		} else if err := kys.IsValid(nil); err != nil {
			return err
		} else {
			cmd.keys = kys
		}
	}

	return nil
}

func (cmd *CreateContractAccountCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	var items []currency.CreateContractAccountsItem

	ams := make([]mitumcurrency.Amount, len(cmd.Amounts))
	for i := range cmd.Amounts {
		a := cmd.Amounts[i]
		am := mitumcurrency.NewAmount(a.Big, a.CID)
		if err := am.IsValid(nil); err != nil {
			return nil, err
		}

		ams[i] = am
	}

	item := currency.NewCreateContractAccountsItemMultiAmounts(cmd.keys, ams)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := currency.NewCreateContractAccountsFact([]byte(cmd.Token), cmd.sender, items)

	op, err := currency.NewCreateContractAccounts(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create create-contract-account operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create create-contract-account operation")
	}

	return op, nil
}
