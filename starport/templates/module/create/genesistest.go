package modulecreate

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

func genesisTestCtx(appName, modulePath, moduleName string) *genny.Generator {
	g := genny.New()
	ctx := plush.NewContext()
	ctx.Set("moduleName", moduleName)
	ctx.Set("modulePath", modulePath)
	ctx.Set("appName", appName)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", moduleName))
	return g
}

// AddGenesisModuleTest returns the generator to generate genesis_test.go
func AddGenesisModuleTest(appName, modulePath, moduleName string) (*genny.Generator, error) {
	g := genesisTestCtx(appName, modulePath, moduleName)
	return g, g.Box(genesisModuleTestTemplate)
}

// AddGenesisTypesTest returns the generator to generate types/genesis_test.go
func AddGenesisTypesTest(appName, modulePath, moduleName string) (*genny.Generator, error) {
	g := genesisTestCtx(appName, modulePath, moduleName)
	return g, g.Box(genesisTypesTestTemplate)
}
