package chain

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/pkg/archive"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/checksum"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/exec"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/gocmd"
)

const releaseDir = "release"
const checksumTxt = "checksum.txt"

// Build builds and installs app binaries.
func (c *Chain) Build(ctx context.Context) (binaryName string, err error) {
	if err := c.setup(ctx); err != nil {
		return "", err
	}

	if err := c.build(ctx); err != nil {
		return "", err
	}

	return c.Binary()
}

func (c *Chain) build(ctx context.Context) (err error) {
	defer func() {
		var exitErr *exec.ExitError

		if errors.As(err, &exitErr) {
			err = &CannotBuildAppError{err}
		}
	}()

	if err := c.generateAll(ctx); err != nil {
		return err
	}

	buildFlags, err := c.preBuild(ctx)
	if err != nil {
		return err
	}

	return gocmd.InstallAll(ctx, c.app.Path, buildFlags)
}

// BuildRelease builds binaries for a release. targets is a list
// of GOOS:GOARCH when provided. It defaults to your system when no targets provided.
// prefix is used as prefix to tarballs containing each target.
func (c *Chain) BuildRelease(ctx context.Context, prefix string, targets ...string) (releasePath string, err error) {
	if prefix == "" {
		prefix = c.app.Name
	}
	if len(targets) == 0 {
		targets = []string{gocmd.BuildTarget(runtime.GOOS, runtime.GOARCH)}
	}

	// prepare for build.
	if err := c.setup(ctx); err != nil {
		return "", err
	}

	buildFlags, err := c.preBuild(ctx)
	if err != nil {
		return "", err
	}

	releasePath = filepath.Join(c.app.Path, releaseDir)

	// reset the release dir.
	if err := os.RemoveAll(releasePath); err != nil {
		return "", err
	}
	if err := os.Mkdir(releasePath, 0755); err != nil {
		return "", err
	}

	for _, t := range targets {
		// build binary for a target, tarball it and save it under the release dir.
		goos, goarch, err := gocmd.ParseTarget(t)
		if err != nil {
			return "", err
		}

		out, err := os.MkdirTemp("", "")
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(out)

		buildOptions := []exec.Option{
			exec.StepOption(step.Env(
				cmdrunner.Env(gocmd.EnvGOOS, goos),
				cmdrunner.Env(gocmd.EnvGOARCH, goarch),
			)),
		}

		if err := gocmd.BuildAll(ctx, out, c.app.Path, buildFlags, buildOptions...); err != nil {
			return "", err
		}

		tarr, err := archive.Tar(out, archive.Gzip)
		if err != nil {
			return "", err
		}

		tarName := fmt.Sprintf("%s_%s_%s.tar.gz", prefix, goos, goarch)
		tarPath := filepath.Join(releasePath, tarName)

		tarf, err := os.Create(tarPath)
		if err != nil {
			return "", err
		}
		defer tarf.Close()

		if _, err := io.Copy(tarf, tarr); err != nil {
			return "", err
		}
		tarf.Close()
	}

	checksumPath := filepath.Join(releasePath, checksumTxt)

	// create a checksum.txt and return with the path to release dir.
	return releasePath, checksum.Sum(releasePath, checksumPath)
}

func (c *Chain) preBuild(ctx context.Context) (buildFlags []string, err error) {
	chainID, err := c.ID()
	if err != nil {
		return nil, err
	}

	binary, err := c.Binary()
	if err != nil {
		return nil, err
	}

	ldflags := gocmd.Ldflags(
		fmt.Sprintf("-X github.com/cosmos/cosmos-sdk/version.Name=%s", strings.Title(c.app.Name)),
		fmt.Sprintf("-X github.com/cosmos/cosmos-sdk/version.AppName=%sd", c.app.Name),
		fmt.Sprintf("-X github.com/cosmos/cosmos-sdk/version.Version=%s", c.sourceVersion.tag),
		fmt.Sprintf("-X github.com/cosmos/cosmos-sdk/version.Commit=%s", c.sourceVersion.hash),
		fmt.Sprintf("-X %s/cmd/%s/cmd.ChainID=%s", c.app.ImportPath, binary, chainID),
	)
	buildFlags = []string{
		gocmd.FlagMod, gocmd.FlagModValueReadOnly,
		gocmd.FlagLdflags, ldflags,
	}

	fmt.Fprintln(c.stdLog().out, "📦 Installing dependencies...")

	if err := gocmd.ModTidy(ctx, c.app.Path); err != nil {
		return nil, err
	}
	if err := gocmd.ModVerify(ctx, c.app.Path); err != nil {
		return nil, err
	}

	fmt.Fprintln(c.stdLog().out, "🛠️  Building the blockchain...")

	return buildFlags, nil
}
