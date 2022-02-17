package scaffolder

import (
	"context"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gobuffalo/genny"
	"github.com/tendermint/flutter/v2"
	"github.com/tendermint/vue"

	"github.com/tendermint/starport/starport/pkg/giturl"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/localfs"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/templates/app"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
)

var (
	commitMessage = "Initialized with Starport"
	devXAuthor    = &object.Signature{
		Name:  "Developer Experience team at Tendermint",
		Email: "hello@tendermint.com",
		When:  time.Now(),
	}
)

// Init initializes a new app with name and given options.
func Init(tracer *placeholder.Tracer, root, name, addressPrefix string, noDefaultModule bool) (path string, err error) {
	if root, err = filepath.Abs(root); err != nil {
		return "", err
	}

	pathInfo, err := gomodulepath.Parse(name)
	if err != nil {
		return "", err
	}

	path = filepath.Join(root, pathInfo.Root)

	// create the project
	if err := generate(tracer, pathInfo, addressPrefix, path, noDefaultModule); err != nil {
		return "", err
	}

	if err := finish(path, pathInfo.RawPath); err != nil {
		return "", err
	}

	// initialize git repository and perform the first commit
	if err := initGit(path); err != nil {
		return "", err
	}

	return path, nil
}

//nolint:interfacer
func generate(
	tracer *placeholder.Tracer,
	pathInfo gomodulepath.Path,
	addressPrefix,
	absRoot string,
	noDefaultModule bool,
) error {
	gu, err := giturl.Parse(pathInfo.RawPath)
	if err != nil {
		return err
	}

	g, err := app.New(&app.Options{
		// generate application template
		ModulePath:       pathInfo.RawPath,
		AppName:          pathInfo.Package,
		AppPath:          absRoot,
		OwnerName:        owner(pathInfo.RawPath),
		OwnerAndRepoName: gu.UserAndRepo(),
		BinaryNamePrefix: pathInfo.Root,
		AddressPrefix:    addressPrefix,
	})
	if err != nil {
		return err
	}

	run := func(runner *genny.Runner, gen *genny.Generator) error {
		runner.With(gen)
		runner.Root = absRoot
		return runner.Run()
	}
	if err := run(genny.WetRunner(context.Background()), g); err != nil {
		return err
	}

	// generate module template
	if !noDefaultModule {
		opts := &modulecreate.CreateOptions{
			ModuleName: pathInfo.Package, // App name
			ModulePath: pathInfo.RawPath,
			AppName:    pathInfo.Package,
			AppPath:    absRoot,
			OwnerName:  owner(pathInfo.RawPath),
			IsIBC:      false,
		}
		g, err = modulecreate.NewStargate(opts)
		if err != nil {
			return err
		}
		if err := run(genny.WetRunner(context.Background()), g); err != nil {
			return err
		}
		g = modulecreate.NewStargateAppModify(tracer, opts)
		if err := run(genny.WetRunner(context.Background()), g); err != nil {
			return err
		}

	}

	// generate the vue app.
	return Vue(filepath.Join(absRoot, "vue"))
}

// Vue scaffolds a Vue.js app for a chain.
func Vue(path string) error {
	return localfs.Save(vue.Boilerplate(), path)
}

// Flutter scaffolds a Flutter app for a chain.
func Flutter(path string) error {
	return localfs.Save(flutter.Boilerplate(), path)
}

func initGit(path string) error {
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	if _, err := wt.Add("."); err != nil {
		return err
	}
	_, err = wt.Commit(commitMessage, &git.CommitOptions{
		All:    true,
		Author: devXAuthor,
	})
	return err
}
