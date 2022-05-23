package example

/*
import (
	"sync"

	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var configContractAccountProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ConfigContractAccountProcessor)
	},
}

func (ConfigContractAccount) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type ConfigContractAccountProcessor struct {
	cp *currency.CurrencyPool
	ConfigContractAccount
	cs  state.State          // contract account status state
	sb  currency.AmountState // sender amount state
	fee currency.Big
	as  extension.ContractAccount // contract account status value
}

func NewConfigContractAccountProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(ConfigContractAccount)
		if !ok {
			return nil, errors.Errorf("not ConfigContractAccount, %T", op)
		}

		opp := configContractAccountProcessorPool.Get().(*ConfigContractAccountProcessor)

		opp.cp = cp
		opp.ConfigContractAccount = i
		opp.cs = nil
		opp.sb = currency.AmountState{}
		opp.fee = currency.ZeroBig
		opp.as = extension.ContractAccount{}

		return opp, nil
	}
}

func (opp *ConfigContractAccountProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(ConfigContractAccountFact)

	// check existence of target account state
	// keep target account state
	st, err := existsState(currency.StateKeyAccount(fact.target), "target keys", getState)
	if err != nil {
		return nil, err
	}
	/*
		ks, err := currency.StateKeysValue(st)
		if err != nil {
			return nil, err
		}
		k, ok := ks.(ContractAccountKeys)
		if !ok {
			return nil, errors.Errorf("contract account keys is not type of ContractAccountKeys")
		}
		emptykeys := NewContractAccountKeys()
		if !k.Equal(emptykeys) {
			return nil, errors.Errorf("not contract account, contract account keys is not empty contract account keys")
		}
*/
/*
	// check not existence of contract account status state
	// check sender matched with contract account owner
	st, err = existsState(extension.StateKeyContractAccount(fact.target), "contract account status", getState)
	if err != nil {
		return nil, err
	}

	v, err := extension.StateContractAccountValue(st)
	if err != nil {
		return nil, err
	}
	if !v.Owner().Equal(fact.sender) {
		return nil, errors.Errorf("contract account owner, %q is not matched with %q", v.Owner(), fact.sender)
	}
	if !v.IsActive() {
		return nil, errors.Errorf("contract account is already deactivated, %q", fact.target)
	}
	opp.cs = st
	opp.as = v

	// check sender has amount of currency
	// keep amount state of sender
	st, err = existsState(currency.StateKeyBalance(fact.sender, fact.currency), "balance of target", getState)
	if err != nil {
		return nil, err
	}
	opp.sb = currency.NewAmountState(st, fact.currency)

	// check fact sign
	if err = checkFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
		return nil, errors.Wrap(err, "invalid signing")
	}

	// check feeer
	feeer, found := opp.cp.Feeer(fact.currency)
	if !found {
		return nil, operation.NewBaseReasonError("currency, %q not found of Deactivate", fact.currency)
	}

	// get fee value
	// keep fee value
	fee, err := feeer.Fee(currency.ZeroBig)
	if err != nil {
		return nil, operation.NewBaseReasonErrorFromError(err)
	}
	switch b, err := currency.StateBalanceValue(opp.sb); {
	case err != nil:
		return nil, operation.NewBaseReasonErrorFromError(err)
	case b.Big().Compare(fee) < 0:
		return nil, operation.NewBaseReasonError("insufficient balance with fee")
	default:
		opp.fee = fee
	}

	return opp, nil
}

func (opp *ConfigContractAccountProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(ConfigContractAccountFact)

	opp.sb = opp.sb.Sub(opp.fee).AddFee(opp.fee)
	v := configExample{}
	st, err := setStateConfigExampleValue(opp.cs, v)
	if err != nil {
		return operation.NewBaseReasonErrorFromError(err)
	}
	return setState(fact.Hash(), st, opp.sb)
}

func (opp *ConfigContractAccountProcessor) Close() error {
	opp.cp = nil
	opp.ConfigContractAccount = ConfigContractAccount{}
	opp.cs = nil
	opp.sb = currency.AmountState{}
	opp.fee = currency.ZeroBig

	configContractAccountProcessorPool.Put(opp)

	return nil
}
*/
