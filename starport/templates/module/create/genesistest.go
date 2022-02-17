package modulecreate

import (
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"

	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/field/plushhelpers"
)

// AddGenesisTest returns the generator to generate genesis_test.go files
func AddGenesisTest(appPath, appName, modulePath, moduleName string, isIBC bool) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(fsGenesisTest, "genesistest/", appPath)
	)

	ctx := plush.NewContext()
	ctx.Set("moduleName", moduleName)
	ctx.Set("modulePath", modulePath)
	ctx.Set("appName", appName)
	ctx.Set("isIBC", isIBC)
	ctx.Set("title", strings.Title)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", moduleName))

	if err := xgenny.Box(g, template); err != nil {
		return nil, err
	}

	return g, nil
}
