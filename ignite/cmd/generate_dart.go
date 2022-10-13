package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateDart() *cobra.Command {
	c := &cobra.Command{
		Use:     "dart",
		Short:   "Generate a Dart client",
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    generateDartHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func generateDartHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	session.StartSpinner("Generating...")

	c, err := newChainWithHomeFlags(
		cmd,
		chain.EnableThirdPartyModuleCodegen(),
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
	)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GenerateDart()); err != nil {
		return err
	}

	return session.Println("⛏️  Generated Dart client.")
}
