package example

/*
import (
	"fmt"
	"strings"

	"github.com/ProtoconNet/mitum-account-extension/extension"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util"
)

var (
	StateKeyExampleConfigSuffix = ":exampleconfig"
)

func StateKeyExampleConfig(a base.Address) string {
	return fmt.Sprintf("%s%s", a.String(), StateKeyExampleConfigSuffix)
}

func IsStateExampleConfigKey(key string) bool {
	return strings.HasSuffix(key, StateKeyExampleConfigSuffix)
}

func StateExampleConfigValue(st state.State) (extension.BaseConfigData, error) {
	v := st.Value()
	if v == nil {
		return extension.BaseConfigData{}, util.NotFoundError.Errorf("not found BaseConfigData in State")
	}

	s, ok := v.Interface().(extension.BaseConfigData)
	if !ok {
		return extension.BaseConfigData{}, errors.Errorf("invalid BaseConfigData value found, %T", v.Interface())
	}
	return s, nil
}

func setStateConfigExampleValue(st state.State, v configExample) (state.State, error) {
	uv, err := state.NewHintedValue(v)
	if err != nil {
		return nil, err
	}
	return st.SetValue(uv)
}

func checkExistsState(
	key string,
	getState func(key string) (state.State, bool, error),
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return operation.NewBaseReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState func(key string) (state.State, bool, error),
) (state.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, operation.NewBaseReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func notExistsState(
	k,
	name string,
	getState func(key string) (state.State, bool, error),
) (state.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, operation.NewBaseReasonError("%s already exists", name)
	default:
		return st, nil
	}
}
*/
