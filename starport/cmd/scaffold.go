package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewScaffold returns a command that groups scaffolding related sub commands.
func NewScaffold() *cobra.Command {
	c := &cobra.Command{
		Use:     "scaffold [command]",
		Short:   "Scaffold a new blockchain or scaffold components inside an existing one",
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(NewScaffoldChain())
	c.AddCommand(NewScaffoldList())
	c.AddCommand(NewScaffoldMap())
	c.AddCommand(NewScaffoldVue())

	return c
}

func scaffoldType(kind, module, typeName string, typeFields []string, opts scaffolder.AddTypeOption) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	sm, err := sc.AddType(placeholder.New(), opts, module, typeName, typeFields...)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Println(sourceModificationToString(sm))
	fmt.Printf("\n🎉 Created a %s `%s`.\n\n", kind, typeName)

	return nil
}

func flagSetScaffoldType() *flag.FlagSet {
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.String(flagModule, "", "Module to add into. Default is app's main module")
	f.Bool(flagNoMessage, false, "Disable CRUD interaction messages scaffolding")
	return f
}

func flagGetModule(cmd *cobra.Command) string {
	module, _ := cmd.Flags().GetString(flagModule)
	return module
}

func flagGetNoMessage(cmd *cobra.Command) bool {
	noMessage, _ := cmd.Flags().GetBool(flagNoMessage)
	return noMessage
}
