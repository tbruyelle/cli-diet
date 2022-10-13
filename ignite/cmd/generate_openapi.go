package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateOpenAPI() *cobra.Command {
	c := &cobra.Command{
		Use:     "openapi",
		Short:   "Generate generates an OpenAPI spec for your chain from your config.yml",
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    generateOpenAPIHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func generateOpenAPIHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	session.StartSpinner("Generating...")

	c, err := newChainWithHomeFlags(
		cmd,
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

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GenerateOpenAPI()); err != nil {
		return err
	}

	return session.Println("⛏️  Generated OpenAPI spec.")
}
