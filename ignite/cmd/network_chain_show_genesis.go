package ignitecmd

import (
	"os"
	"path/filepath"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
	"github.com/ignite-hq/cli/ignite/services/network/networkchain"
	"github.com/spf13/cobra"
)

func newNetworkChainShowGenesis() *cobra.Command {
	c := &cobra.Command{
		Use:   "genesis [launch-id]",
		Short: "Show the chain genesis file",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainShowGenesisHandler,
	}

	c.Flags().String(flagOut, "./genesis.json", "Path to output Genesis file")

	return c
}

func networkChainShowGenesisHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	out, _ := cmd.Flags().GetString(flagOut)

	nb, launchID, err := networkChainLaunch(cmd, args, session)
	if err != nil {
		return err
	}
	n, err := nb.Network()
	if err != nil {
		return err
	}

	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch))
	if err != nil {
		return err
	}

	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}

	// check if the genesis already exists
	if _, err = os.Stat(genesisPath); os.IsNotExist(err) {
		// fetch the information to construct genesis
		genesisInformation, err := n.GenesisInformation(cmd.Context(), launchID)
		if err != nil {
			return err
		}

		// create the chain in a temp dir
		home := filepath.Join(os.TempDir(), "spn/temp", chainLaunch.ChainID)
		defer os.RemoveAll(home)

		c.SetHome(home)

		err = c.Prepare(cmd.Context(), genesisInformation)
		if err != nil {
			return err
		}

		// get the new genesis path
		genesisPath, err = c.GenesisPath()
		if err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(out), 0744); err != nil {
		return err
	}

	if err := os.Rename(genesisPath, out); err != nil {
		return err
	}

	session.StopSpinner()

	return session.Printf("%s Genesis generated: %s\n", icons.Bullet, out)
}
