package example

/*
import (
	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	ConfigExampleType   = hint.Type("mitum-currency-config-example")
	ConfigExampleHint   = hint.NewHint(ConfigExampleType, "v0.0.1")
	ConfigExampleHinter = configExample{BaseHinter: hint.NewBaseHinter(ConfigExampleHint)}
)

func NewBaseConfigData(id string, address base.Address) (extensioncurrency.BaseConfigData, error) {
	cfg := configExample{
		BaseHinter: hint.NewBaseHinter(ConfigExampleHint),
		id:         id,
		address:    address,
	}
	bcfg, err := extensioncurrency.NewBaseConfigData(cfg)
	if err != nil {
		return extensioncurrency.BaseConfigData{}, err
	}
	return bcfg, nil
}

type configExample struct {
	hint.BaseHinter
	id      string
	address base.Address
}

func (cfg configExample) ID() string {
	return cfg.id
}
func (cfg configExample) ConfigType() hint.Type {
	return cfg.Hint().Type()
}
func (cfg configExample) Hint() hint.Hint {
	return cfg.Hint()
}
func (cfg configExample) Bytes() []byte {
	return nil
}
func (cfg configExample) Hash() valuehash.Hash {
	return cfg.GenerateHash()
}
func (cfg configExample) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(cfg.Bytes())
}
func (cfg configExample) IsValid([]byte) error {
	return nil
}
func (cfg configExample) Address() base.Address {
	return cfg.address
}

func (cfg configExample) SetStateValue(st state.State) (state.State, error) {
	st, err := setStateConfigExampleValue(st, cfg)
	if err != nil {
		return nil, err
	}
	return st, nil
}
*/
