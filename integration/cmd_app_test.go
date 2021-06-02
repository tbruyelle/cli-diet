// +build !relayer

package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnApp(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	_, statErr := os.Stat(filepath.Join(path, "x", "blog"))
	require.False(t, os.IsNotExist(statErr), "the default module should be scaffolded")

	env.EnsureAppIsSteady(path)
}

func TestGenerateAnAppWithNoDefaultModule(t *testing.T) {
	var (
		env  = newEnv(t)
		appName = "blog"
	)

	root := env.TmpDir()
	env.Exec("scaffold an app",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				"app",
				fmt.Sprintf("github.com/test/%s", appName),
				"--no-default-module",
			),
			step.Workdir(root),
		)),
	)

	// Cleanup the home directory of the app
	env.t.Cleanup(func() {
		os.RemoveAll(filepath.Join(env.Home(), fmt.Sprintf(".%s", appName)))
	})

	path := filepath.Join(root, appName)

	_, statErr := os.Stat(filepath.Join(path, "x", "blog"))
	require.True(t, os.IsNotExist(statErr), "the default module should not be scaffolded")

	env.EnsureAppIsSteady(path)
}


func TestGenerateAnAppWithWasm(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("add Wasm module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "import", "wasm"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should not add Wasm module second time",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "import", "wasm"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}

func TestGenerateAStargateAppWithEmptyModule(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "example", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an existing module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "example", "--require-registration"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}
