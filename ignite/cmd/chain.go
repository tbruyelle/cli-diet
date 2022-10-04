package ignitecmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
)

var (
	msgMigration       = "Migrating blockchain config file from v%d to v%d..."
	msgMigrationCancel = "Stopping because config version v%d is required to run the command"
	msgMigrationPrompt = "Your blockchain config version is v%[1]d and the latest is v%[2]d. Would you like to upgrade your config file to v%[2]d?"
)

// NewChain returns a command that groups sub commands related to compiling, serving
// blockchains and so on.
func NewChain() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain [command]",
		Short: "Build, initialize and start a blockchain node or perform other actions on the blockchain",
		Long: `Commands in this namespace let you to build, initialize, and start your
blockchain node locally for development purposes.

To run these commands you should be inside the project's directory so that
Ignite can find the source code. To ensure that you are, run "ls", you should
see the following files in the output: "go.mod", "x", "proto", "app", etc.

By default the "build" command will identify the "main" package of the project,
install dependencies if necessary, set build flags, compile the project into a
binary and install the binary. The "build" command is useful if you just want
the compiled binary, for example, to initialize and start the chain manually. It
can also be used to release your chain's binaries automatically as part of
continuous integration workflow.

The "init" command will build the chain's binary and use it to initialize a
local validator node. By default the validator node will be initialized in your
$HOME directory in a hidden directory that matches the name of your project.
This directory is called a data directory and contains a chain's genesis file
and a validator key. This command is useful if you want to quickly build and
initialize the data directory and use the chain's binary to manually start the
blockchain. The "init" command is meant only for development purposes, not
production.

The "serve" command builds, initializes, and starts your blockchain locally with
a single validator node for development purposes. "serve" also watches the
source code directory for file changes and intelligently
re-builds/initializes/starts the chain, essentially providing "code-reloading".
The "serve" command is meant only for development purposes, not production.

To distinguish between production and development consider the following.

In production, blockchains often run the same software on many validator nodes
that are run by different people and entities. To launch a blockchain in
production, the validator entities coordinate the launch process to start their
nodes simultaneously.

During development, a blockchain can be started locally on a single validator
node. This convenient process lets you restart a chain quickly and iterate
faster. Starting a chain on a single node in development is similar to starting
a traditional web application on a local server.

The "faucet" command lets you send tokens to an address from the "faucet"
account defined in "config.yml". Alternatively, you can use the chain's binary
to send token from any other account that exists on chain.

The "simulate" command helps you start a simulation testing process for your
chain.
`,
		Aliases:           []string{"c"},
		Args:              cobra.ExactArgs(1),
		PersistentPreRunE: configMigrationPreRunHandler,
	}

	// Add flags required for the configMigrationPreRunHandler
	c.PersistentFlags().AddFlagSet(flagSetConfig())
	c.PersistentFlags().AddFlagSet(flagSetYes())

	c.AddCommand(NewChainServe())
	c.AddCommand(NewChainBuild())
	c.AddCommand(NewChainInit())
	c.AddCommand(NewChainFaucet())
	c.AddCommand(NewChainSimulate())

	return c
}

func configMigrationPreRunHandler(cmd *cobra.Command, args []string) (err error) {
	session := cliui.New()
	defer session.Cleanup()

	appPath := flagGetPath(cmd)
	configPath := getConfig(cmd)
	if configPath == "" {
		if configPath, err = chainconfig.LocateDefault(appPath); err != nil {
			return err
		}
	}

	rawCfg, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	version, err := chainconfig.ReadConfigVersion(bytes.NewReader(rawCfg))
	if err != nil {
		return err
	}

	// Config files with older versions must be migrated to the latest before executing the command
	if version != chainconfig.LatestVersion {
		if !getYes(cmd) {
			// Confirm before overwritting the config file
			question := fmt.Sprintf(msgMigrationPrompt, version, chainconfig.LatestVersion)
			if err := session.AskConfirm(question); err != nil {
				if errors.Is(err, promptui.ErrAbort) {
					return fmt.Errorf(msgMigrationCancel, chainconfig.LatestVersion)
				}

				return err
			}

			// Confirm before migrating the config if there are uncommitted changes
			if err := confirmWhenUncommittedChanges(session, appPath); err != nil {
				return err
			}
		} else {
			session.Printf("%s %s\n", icons.Info, colors.Infof(msgMigration, version, chainconfig.LatestVersion))
		}

		file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0o755)
		if err != nil {
			return err
		}

		defer file.Close()

		// Convert the current config to the latest version and update the YAML file
		return chainconfig.MigrateLatest(bytes.NewReader(rawCfg), file)
	}

	return nil
}
