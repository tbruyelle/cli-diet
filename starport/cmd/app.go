package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagAddressPrefix   = "address-prefix"
	flagNoDefaultModule = "no-default-module"
)

// NewApp creates new command named `app` to create Cosmos scaffolds customized
// by the user given options.
func NewApp() *cobra.Command {
	c := &cobra.Command{
		Use:   "app [github.com/org/repo]",
		Short: "Scaffold a new blockchain",
		Long:  "Scaffold a new Cosmos SDK blockchain with a default directory structure",
		Args:  cobra.ExactArgs(1),
		RunE:  appHandler,
	}
	c.Flags().String(flagAddressPrefix, "cosmos", "Address prefix")
	c.Flags().Bool(flagNoDefaultModule, false, "Prevent scaffolding a default module in the app")
	return c
}

func appHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	var (
		name               = args[0]
		addressPrefix, _   = cmd.Flags().GetString(flagAddressPrefix)
		noDefaultModule, _ = cmd.Flags().GetBool(flagNoDefaultModule)
	)

	sc, err := scaffolder.New("",
		scaffolder.AddressPrefix(addressPrefix),
	)
	if err != nil {
		return err
	}

	appdir, err := sc.Init(placeholder.New(), name, noDefaultModule)
	if err != nil {
		return err
	}

	s.Stop()

	message := `
⭐️ Successfully created a new blockchain '%[1]v'.
👉 Get started with the following commands:

 %% cd %[1]v
 %% starport serve

Documentation: https://docs.starport.network
`
	fmt.Printf(message, appdir)

	return nil
}
