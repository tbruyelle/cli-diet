package starportcmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/numbers"
	"github.com/tendermint/starport/starport/services/network"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

// NewNetworkRequestVerify verify the request and simulate the chain.
func NewNetworkRequestVerify() *cobra.Command {
	c := &cobra.Command{
		Use:   "verify [launch-id] [number<,...>]",
		Short: "Verify the request and simulate the chain genesis from them",
		RunE:  networkRequestVerifyHandler,
		Args:  cobra.ExactArgs(2),
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkRequestVerifyHandler(cmd *cobra.Command, args []string) error {
	// initialize network common methods
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	// get the list of request ids
	ids, err := numbers.ParseList(args[1])
	if err != nil {
		return err
	}

	// verify the requests
	if err := verifyRequest(cmd.Context(), nb, launchID, ids...); err != nil {
		fmt.Printf("%s Request(s) %s not valid\n", clispinner.NotOK, numbers.List(ids, "#"))
		return err
	}

	nb.Spinner.Stop()
	fmt.Printf("%s Request(s) %s verified\n", clispinner.OK, numbers.List(ids, "#"))
	return nil
}

// verifyRequest initialize the chain from the launch ID in a temporary directory
// and simulate the launch of the chain from genesis with the request IDs
func verifyRequest(
	ctx context.Context,
	nb NetworkBuilder,
	launchID uint64,
	requestIDs ...uint64,
) error {
	n, err := nb.Network()
	if err != nil {
		return err
	}

	// initialize the chain with a temporary dir
	chainLaunch, err := n.ChainLaunch(ctx, launchID)
	if err != nil {
		return err
	}

	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(homeDir)

	c, err := nb.Chain(
		networkchain.SourceLaunch(chainLaunch),
		networkchain.WithHome(homeDir),
		networkchain.WithKeyringBackend(chaincmd.KeyringBackendTest),
	)
	if err != nil {
		return err
	}

	// fetch the current genesis information and the requests for the chain for simulation
	genesisInformation, err := n.GenesisInformation(ctx, launchID)
	if err != nil {
		return err
	}

	requests, err := n.RequestFromIDs(ctx, launchID, requestIDs...)
	if err != nil {
		return err
	}

	return c.SimulateRequests(ctx, genesisInformation, requests)
}
