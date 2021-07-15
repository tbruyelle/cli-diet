package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewScaffoldMap returns a new command to scaffold a map.
func NewScaffoldMap() *cobra.Command {
	c := &cobra.Command{
		Use:   "map NAME [field]...",
		Short: "CRUD for data stored as key-value pairs",
		Args:  cobra.MinimumNArgs(1),
		RunE:  scaffoldMapHandler,
	}

	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().AddFlagSet(flagSetScaffoldType())

	return c
}

func scaffoldMapHandler(cmd *cobra.Command, args []string) error {
	opts := scaffolder.AddTypeOption{
		NoMessage: flagGetNoMessage(cmd),
		Model:     scaffolder.Map,
	}

	return scaffoldType(flagGetModule(cmd), args[0], args[1:], opts)
}
