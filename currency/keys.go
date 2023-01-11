package currency

import (
	"bytes"
	"sort"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
)

var ContractAccountKeysHint = hint.MustNewHint("mitum-currency-contract-account-keys-v0.0.1")

type ContractAccountKeys struct {
	hint.BaseHinter
	h         util.Hash
	keys      []currency.AccountKey
	threshold uint
}

func NewContractAccountKeys() ContractAccountKeys {
	ks := ContractAccountKeys{BaseHinter: hint.NewBaseHinter(ContractAccountKeysHint), keys: []currency.AccountKey{}, threshold: 100}

	h, err := ks.GenerateHash()
	if err != nil {
		return ContractAccountKeys{}
	}
	ks.h = h

	return ks
}

func (ks ContractAccountKeys) Hash() util.Hash {
	return ks.h
}

func (ks ContractAccountKeys) GenerateHash() (util.Hash, error) {
	return valuehash.NewSHA256(ks.Bytes()), nil
}

func (ks ContractAccountKeys) Bytes() []byte {
	return util.UintToBytes(ks.threshold)
}

func (ks ContractAccountKeys) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false, ks.h); err != nil {
		return err
	}

	if len(ks.keys) > 0 {
		return util.ErrInvalid.Errorf("keys of contract account exist")
	}

	if h, err := ks.GenerateHash(); err != nil {
		return err
	} else if !ks.h.Equal(h) {
		return util.ErrInvalid.Errorf("hash not matched")
	}

	return nil
}

func (ks ContractAccountKeys) Threshold() uint {
	return ks.threshold
}

func (ks ContractAccountKeys) Keys() []currency.AccountKey {
	return ks.keys
}

func (ks ContractAccountKeys) Key(k base.Publickey) (currency.AccountKey, bool) {
	return currency.BaseAccountKey{}, false
}

func (ks ContractAccountKeys) Equal(b currency.AccountKeys) bool {
	if ks.threshold != b.Threshold() {
		return false
	}

	if len(ks.keys) != len(b.Keys()) {
		return false
	}

	sort.Slice(ks.keys, func(i, j int) bool {
		return bytes.Compare(ks.keys[i].Key().Bytes(), ks.keys[j].Key().Bytes()) < 0
	})

	bkeys := b.Keys()
	sort.Slice(bkeys, func(i, j int) bool {
		return bytes.Compare(bkeys[i].Key().Bytes(), bkeys[j].Key().Bytes()) < 0
	})

	for i := range ks.keys {
		if !ks.keys[i].Equal(bkeys[i]) {
			return false
		}
	}

	return true
}

func checkThreshold(fs []base.Sign, keys currency.AccountKeys) error {
	var sum uint
	for i := range fs {
		ky, found := keys.Key(fs[i].Signer())
		if !found {
			return errors.Errorf("unknown key found, %s", fs[i].Signer())
		}
		sum += ky.Weight()
	}

	if sum < keys.Threshold() {
		return errors.Errorf("not passed threshold, sum=%d < threshold=%d", sum, keys.Threshold())
	}

	return nil
}
