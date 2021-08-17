package starportcmd

import (
	"github.com/spf13/cobra"
)

// NewRelayer returns a new relayer command.
func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:   "relayer",
		Short: "Connect blockchains by using IBC protocol",
	}

	c.AddCommand(NewRelayerConfigure())
	c.AddCommand(NewRelayerConnect())

	return c
}
