package cmds

import (
	"os"

	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/logging"
)

type VerifyKeyCommand struct {
	printCommand
	Key    StringLoad `arg:"" name:"key" help:"key" required:""`
	Quite  bool       `name:"quite" short:"q" help:"keep silence"`
	JSON   bool       `name:"json" help:"json output format (default: false)" optional:"" default:"false"`
	Pretty bool       `name:"pretty" help:"pretty format"`
}

func (cmd *VerifyKeyCommand) Run(log logging.Logger) error {
	var pk key.Key
	if k, err := loadKey(cmd.Key.Bytes()); err != nil {
		if cmd.Quite {
			os.Exit(1)
		}

		return err
	} else {
		pk = k
	}

	log.Debug().Interface("key", pk).Msg("key parsed")

	if cmd.Quite {
		return nil
	}

	var priv key.Privatekey
	var pub key.Publickey
	switch t := pk.(type) {
	case key.Privatekey:
		priv = t
		pub = t.Publickey()
	case key.Publickey:
		pub = t
	}

	if cmd.JSON {
		m := map[string]interface{}{
			"publickey": map[string]interface{}{
				"hint": pub.Hint(),
				"key":  pub.String(),
			},
		}

		if priv != nil {
			m["privtekey"] = map[string]interface{}{
				"hint": priv.Hint(),
				"key":  priv.String(),
			}
		}

		cmd.pretty(cmd.Pretty, m)

		return nil
	}

	if priv != nil {
		cmd.print("privatekey hint: %s", priv.Hint().Verbose())
		cmd.print("     privatekey: %s", hint.HintedString(priv.Hint(), priv.String()))
	}

	cmd.print(" publickey hint: %s", pub.Hint().Verbose())
	cmd.print("      publickey: %s", hint.HintedString(pub.Hint(), pub.String()))

	return nil
}
