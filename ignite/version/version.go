package version

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"text/tabwriter"

	"github.com/blang/semver"
	"github.com/google/go-github/v37/github"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gitpod"
	"github.com/ignite/cli/ignite/pkg/xexec"
)

const (
	versionDev     = "development"
	versionNightly = "nightly"
)

// Version is the semantic version of Ignite CLI.
var Version = versionDev

// CheckNext checks whether there is a new version of Ignite CLI.
func CheckNext(ctx context.Context) (isAvailable bool, version string, err error) {
	if Version == versionDev || Version == versionNightly {
		return false, "", nil
	}

	latest, _, err := github.
		NewClient(nil).
		Repositories.
		GetLatestRelease(ctx, "ignite", "cli")
	if err != nil {
		return false, "", err
	}

	if latest.TagName == nil {
		return false, "", nil
	}

	currentVersion, err := semver.ParseTolerant(Version)
	if err != nil {
		return false, "", err
	}

	latestVersion, err := semver.ParseTolerant(*latest.TagName)
	if err != nil {
		return false, "", err
	}

	isAvailable = latestVersion.GT(currentVersion)

	return isAvailable, *latest.TagName, nil
}

// Long generates a detailed version info.
func Long(ctx context.Context) string {
	var (
		w          = &tabwriter.Writer{}
		b          = &bytes.Buffer{}
		date       = "undefined"
		head       = "undefined"
		modified   bool
		sdkVersion = "undefined"
	)
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/cosmos/cosmos-sdk" {
				sdkVersion = dep.Version
				break
			}
		}
		for _, kv := range info.Settings {
			switch kv.Key {
			case "vcs.revision":
				head = kv.Value
			case "vcs.time":
				date = kv.Value
			case "vcs.modified":
				modified = kv.Value == "true"
			}
		}
		if modified {
			// add * suffix to head to indicate the sources have been modified.
			head += "*"
		}
	}

	write := func(k string, v interface{}) {
		fmt.Fprintf(w, "%s:\t%v\n", k, v)
	}

	w.Init(b, 0, 8, 0, '\t', 0)

	write("Ignite CLI version", Version)
	write("Ignite CLI build date", date)
	write("Ignite CLI source hash", head)
	write("Cosmos SDK version", sdkVersion)

	write("Your OS", runtime.GOOS)
	write("Your arch", runtime.GOARCH)

	cmdOut := &bytes.Buffer{}

	nodeJSCmd := "node"
	if xexec.IsCommandAvailable(nodeJSCmd) {
		cmdOut.Reset()

		err := exec.Exec(ctx, []string{nodeJSCmd, "-v"}, exec.StepOption(step.Stdout(cmdOut)))
		if err == nil {
			write("Your Node.js version", strings.TrimSpace(cmdOut.String()))
		}
	}

	cmdOut.Reset()
	err := exec.Exec(ctx, []string{"go", "version"}, exec.StepOption(step.Stdout(cmdOut)))
	if err != nil {
		panic(err)
	}
	write("Your go version", strings.TrimSpace(cmdOut.String()))

	unameCmd := "uname"
	if xexec.IsCommandAvailable(unameCmd) {
		cmdOut.Reset()

		err := exec.Exec(ctx, []string{unameCmd, "-a"}, exec.StepOption(step.Stdout(cmdOut)))
		if err == nil {
			write("Your uname -a", strings.TrimSpace(cmdOut.String()))
		}
	}

	if cwd, err := os.Getwd(); err == nil {
		write("Your cwd", cwd)
	}

	write("Is on Gitpod", gitpod.IsOnGitpod())

	w.Flush()

	return b.String()
}
