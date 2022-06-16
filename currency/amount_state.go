package currency

import (
	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	amountStateType = hint.Type("mitum-currency-amount-state")
	amountStateHint = hint.NewHint(amountStateType, "v0.0.1")
)

type AmountState struct {
	state.State
	contractid ContractID
	cid        currency.CurrencyID
	add        currency.Big
	fee        currency.Big
}

func NewAmountState(st state.State, cid currency.CurrencyID, id ContractID) AmountState {
	if sst, ok := st.(AmountState); ok {
		return sst
	}

	return AmountState{
		State:      st,
		contractid: id,
		cid:        cid,
		add:        currency.ZeroBig,
		fee:        currency.ZeroBig,
	}
}

func (AmountState) Hint() hint.Hint {
	return amountStateHint
}

func (st AmountState) IsValid(b []byte) error {
	if err := isvalid.Check(b, false, st.State); err != nil {
		return err
	}

	if !st.fee.OverNil() {
		return isvalid.InvalidError.Errorf("invalid fee; under zero, %v", st.fee)
	}

	return nil
}

func (st AmountState) Merge(b state.State) (state.State, error) {
	var am AmountValue
	if b, err := StateBalanceValue(b); err != nil {
		if !errors.Is(err, util.NotFoundError) {
			return nil, err
		}
		am = AmountValue{}
	} else if b.ID() != st.contractid {
		return nil, errors.Errorf("contractid is not matched with state to be merged, %v", b.ID())
	} else {
		am = b
	}
	return SetStateBalanceValue(
		st.AddFee(b.(AmountState).fee),
		am.WithBig(am.amount.Big().Add(st.add)),
	)
}

func (st AmountState) Currency() currency.CurrencyID {
	return st.cid
}

func (st AmountState) Fee() currency.Big {
	return st.fee
}

func (st AmountState) AddFee(fee currency.Big) AmountState {
	st.fee = st.fee.Add(fee)

	return st
}

func (st AmountState) Add(a currency.Big) AmountState {
	st.add = st.add.Add(a)

	return st
}

func (st AmountState) Sub(a currency.Big) AmountState {
	st.add = st.add.Sub(a)

	return st
}

func (st AmountState) SetContractID(cid ContractID) AmountState {
	st.contractid = cid

	return st
}

func (st AmountState) SetValue(v state.Value) (state.State, error) {
	s, err := st.State.SetValue(v)
	if err != nil {
		return nil, err
	}
	st.State = s

	return st, nil
}

func (st AmountState) SetHash(h valuehash.Hash) (state.State, error) {
	s, err := st.State.SetHash(h)
	if err != nil {
		return nil, err
	}
	st.State = s

	return st, nil
}

func (st AmountState) SetHeight(h base.Height) state.State {
	st.State = st.State.SetHeight(h)

	return st
}

func (st AmountState) SetPreviousHeight(h base.Height) (state.State, error) {
	s, err := st.State.SetPreviousHeight(h)
	if err != nil {
		return nil, err
	}
	st.State = s

	return st, nil
}

func (st AmountState) SetOperation(ops []valuehash.Hash) state.State {
	st.State = st.State.SetOperation(ops)

	return st
}

func (st AmountState) Clear() state.State {
	st.State = st.State.Clear()

	st.add = currency.ZeroBig
	st.fee = currency.ZeroBig

	return st
}
