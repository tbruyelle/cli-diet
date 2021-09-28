package plushhelpers

import (
	"strings"

	"github.com/gobuffalo/plush"
)

// ExtendPlushContext sets available helpers on the provided context.
func ExtendPlushContext(ctx *plush.Context) {
	ctx.Set("castArg", castArg)
	ctx.Set("castToBytes", CastToBytes)
	ctx.Set("castToString", CastToString)
	ctx.Set("genValidArg", GenerateValidArg)
	ctx.Set("genUniqueArg", GenerateUniqueArg)
	ctx.Set("genValidIndex", GenerateValidIndex)
	ctx.Set("genNotFoundIndex", GenerateNotFoundIndex)
	ctx.Set("title", strings.Title)
}
