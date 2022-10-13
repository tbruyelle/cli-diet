package networkchain

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ignite/cli/ignite/pkg/cache"
	cosmosgenesis "github.com/ignite/cli/ignite/pkg/cosmosutil/genesis"
	"github.com/ignite/cli/ignite/pkg/events"
)

// Init initializes blockchain by building the binaries and running the init command and
// create the initial genesis of the chain, and set up a validator key
func (c *Chain) Init(ctx context.Context, cacheStorage cache.Storage) error {
	chainHome, err := c.chain.Home()
	if err != nil {
		return err
	}

	// cleanup home dir of app if exists.
	if err = os.RemoveAll(chainHome); err != nil {
		return err
	}

	// build the chain and initialize it with a new validator key
	if _, err := c.Build(ctx, cacheStorage); err != nil {
		return err
	}

	c.ev.Send("Initializing the blockchain", events.ProgressStarted())

	if err = c.chain.Init(ctx, false); err != nil {
		return err
	}

	c.ev.Send("Blockchain initialized", events.ProgressFinished())

	// initialize and verify the genesis
	if err = c.initGenesis(ctx); err != nil {
		return err
	}

	c.isInitialized = true

	return nil
}

// initGenesis creates the initial genesis of the genesis depending on the initial genesis type (default, url, ...)
func (c *Chain) initGenesis(ctx context.Context) error {
	c.ev.Send("Computing the Genesis", events.ProgressStarted())

	genesisPath, err := c.chain.GenesisPath()
	if err != nil {
		return err
	}

	// remove existing genesis
	if err := os.RemoveAll(genesisPath); err != nil {
		return err
	}

	// if the blockchain has a genesis URL, the initial genesis is fetched from the URL
	// otherwise, the default genesis is used, which requires no action since the default genesis is generated from the init command
	if c.genesisURL != "" {
		c.ev.Send("Fetching custom Genesis from URL", events.ProgressStarted())
		genesis, err := cosmosgenesis.FromURL(ctx, c.genesisURL, genesisPath)
		if err != nil {
			return err
		}

		if genesis.TarballPath() != "" {
			c.ev.Send(
				fmt.Sprintf("Extracted custom Genesis from tarball at %s", genesis.TarballPath()),
				events.ProgressFinished(),
			)
		} else {
			c.ev.Send("Custom Genesis JSON from URL fetched", events.ProgressFinished())
		}

		hash, err := genesis.Hash()
		if err != nil {
			return err
		}

		// if the blockchain has been initialized with no genesis hash, we assign the fetched hash to it
		// otherwise we check the genesis integrity with the existing hash
		if c.genesisHash == "" {
			c.genesisHash = hash
		} else if hash != c.genesisHash {
			return fmt.Errorf("genesis from URL %s is invalid. expected hash %s, actual hash %s", c.genesisURL, c.genesisHash, hash)
		}

		genBytes, err := genesis.Bytes()
		if err != nil {
			return err
		}

		// replace the default genesis with the fetched genesis
		if err := os.WriteFile(genesisPath, genBytes, 0o644); err != nil {
			return err
		}
	} else {
		// default genesis is used, init CLI command is used to generate it
		cmd, err := c.chain.Commands(ctx)
		if err != nil {
			return err
		}

		// TODO: use validator moniker https://github.com/ignite/cli/issues/1834
		if err := cmd.Init(ctx, "moniker"); err != nil {
			return err
		}

	}

	// check the initial genesis is valid
	if err := c.checkInitialGenesis(ctx); err != nil {
		return err
	}

	c.ev.Send("Genesis initialized", events.ProgressFinished())
	return nil
}

// checkGenesis checks the stored genesis is valid
func (c *Chain) checkInitialGenesis(ctx context.Context) error {
	// perform static analysis of the chain with the validate-genesis command.
	chainCmd, err := c.chain.Commands(ctx)
	if err != nil {
		return err
	}

	// the chain initial genesis should not contain gentx, gentxs should be added through requests
	genesisPath, err := c.chain.GenesisPath()
	if err != nil {
		return err
	}
	chainGenesis, err := cosmosgenesis.FromPath(genesisPath)
	if err != nil {
		return err
	}
	gentxCount, err := chainGenesis.GentxCount()
	if err != nil {
		return err
	}
	if gentxCount > 0 {
		return errors.New("the initial genesis for the chain should not contain gentx")
	}

	return chainCmd.ValidateGenesis(ctx)

	// TODO: static analysis of the genesis with validate-genesis doesn't check the full validity of the genesis
	// example: gentxs formats are not checked
	// to perform a full validity check of the genesis we must try to start the chain with sample accounts
}
