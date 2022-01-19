//go:build !relayer
// +build !relayer

package map_test

import (
	"path/filepath"
	"testing"

	envtest "github.com/tendermint/starport/integration"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestCreateMapWithStargate(t *testing.T) {
	var (
		env  = envtest.New(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a map",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "user", "user-id", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a map with custom path",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "appPath", "email", "--path", filepath.Join(path, "app")),
			step.Workdir(filepath.Dir(path)),
		)),
	))

	env.Must(env.Exec("create a map with no message",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "nomessage", "email", "--no-message"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "module", "example", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				"s",
				"list",
				"user",
				"email",
				"--module",
				"example",
				"--no-simulation",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a map with a typename that already exist",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a map in a custom module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "mapUser", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a map with a custom field type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "mapDetail", "user:MapUser", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a map with Coin and []Coin",
		step.NewSteps(step.New(
			step.Exec("starport",
				"s",
				"map",
				"salary",
				"numInt:int",
				"numsInt:array.int",
				"numsIntAlias:ints",
				"numUint:uint",
				"numsUint:array.uint",
				"numsUintAlias:uints",
				"textString:string",
				"textStrings:array.string",
				"textStringsAlias:strings",
				"textCoin:coin",
				"textCoins:array.coin",
				"textCoinsAlias:coins",
				"--module",
				"example",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a map with index",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				"s",
				"map",
				"map_with_index",
				"email",
				"emailIds:ints",
				"--index",
				"foo:string,bar:int,foobar:uint,barFoo:bool",
				"--module",
				"example",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a map with invalid index",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				"s",
				"map",
				"map_with_invalid_index",
				"email",
				"--index",
				"foo:strings,bar:ints",
				"--module",
				"example",
			),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a message and a map with no-message flag to check conflicts",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "message", "create-scavenge", "description"),
			step.Exec("starport", "s", "map", "scavenge", "description", "--no-message"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a map with duplicated indexes",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "map_with_duplicated_index", "email", "--index", "foo,foo"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a map with an index present in fields",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "map_with_invalid_index", "email", "--index", "email"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}
