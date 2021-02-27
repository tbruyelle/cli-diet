package scaffolder

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/tendermint/starport/starport/templates/typed/indexed"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/templates/typed"
)

const (
	TypeString = "string"
	TypeBool   = "bool"
	TypeInt32  = "int32"
)

type AddTypeOption struct {
	Legacy  bool
	Indexed bool
}

// AddType adds a new type stype to scaffolded app by using optional type fields.
func (s *Scaffolder) AddType(addTypeOptions AddTypeOption, moduleName string, stype string, fields ...string) error {
	version, err := s.version()
	if err != nil {
		return err
	}
	majorVersion := version.Major()
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return err
	}

	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = path.Package
	}
	ok, err := moduleExists(s.path, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't exist", moduleName)
	}

	// Ensure the type name is not a Go reserved name, it would generate an incorrect code
	if isGoReservedWord(stype) {
		return fmt.Errorf("%s can't be used as a type name", stype)
	}

	// Check type is not already created
	ok, err = isTypeCreated(s.path, moduleName, stype)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("%s type is already added", stype)
	}

	// Parse provided field
	tFields, err := parseFields(fields)
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &typed.Options{
			AppName:    path.Package,
			ModulePath: path.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(path.RawPath),
			TypeName:   stype,
			Fields:     tFields,
			Legacy:     addTypeOptions.Legacy,
		}
	)
	// generate depending on the version
	if majorVersion == cosmosver.Launchpad {
		if addTypeOptions.Indexed {
			return errors.New("indexed types not supported on Launchpad")
		}

		g, err = typed.NewLaunchpad(opts)
	} else {
		// Check if indexed type
		if addTypeOptions.Indexed {
			g, err = indexed.NewStargate(opts)
		} else {
			// Scaffolding a type with ID

			// check if the msgServer convention is used
			var msgServerDefined bool
			msgServerDefined, err = isMsgServerDefined(s.path, moduleName)
			if err != nil {
				return err
			}
			if !msgServerDefined {
				opts.Legacy = true
			}
			g, err = typed.NewStargate(opts)
		}
	}
	if err != nil {
		return err
	}
	run := genny.WetRunner(context.Background())
	run.With(g)
	if err := run.Run(); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := s.protoc(pwd, path.RawPath, majorVersion); err != nil {
		return err
	}
	return fmtProject(pwd)
}

// parseFields parses the provided fields, analyses the types and checks there is no duplicated field
func parseFields(fields []string) ([]typed.Field, error) {
	// Used to check duplicated field
	existingFields := make(map[string]bool)

	var tFields []typed.Field
	for _, f := range fields {
		fs := strings.Split(f, ":")
		name := fs[0]

		// Ensure the field name is not a Go reserved name, it would generate an incorrect code
		if isGoReservedWord(name) {
			return tFields, fmt.Errorf("%s can't be used as a field name", name)
		}

		// Ensure the field is not duplicated
		if _, exists := existingFields[name]; exists {
			return tFields, fmt.Errorf("the field %s is duplicated", name)
		}
		existingFields[name] = true

		datatypeName, datatype := TypeString, TypeString
		acceptedTypes := map[string]string{
			"string": TypeString,
			"bool":   TypeBool,
			"int":    TypeInt32,
		}
		isTypeSpecified := len(fs) == 2
		if isTypeSpecified {
			if t, ok := acceptedTypes[fs[1]]; ok {
				datatype = t
				datatypeName = fs[1]
			} else {
				return tFields, fmt.Errorf("the field type %s doesn't exist", fs[1])
			}
		}
		tFields = append(tFields, typed.Field{
			Name:         name,
			Datatype:     datatype,
			DatatypeName: datatypeName,
		})
	}

	return tFields, nil
}

func isTypeCreated(appPath, moduleName, typeName string) (isCreated bool, err error) {
	abspath, err := filepath.Abs(filepath.Join(appPath, "x", moduleName, "types"))
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	all, err := parser.ParseDir(fset, abspath, func(os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return false, err
	}
	// To check if the file is created, we check if the message MsgCreate[TypeName] or Msg[TypeName] is defined
	for _, pkg := range all {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(x ast.Node) bool {
				typeSpec, ok := x.(*ast.TypeSpec)
				if !ok {
					return true
				}
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					return true
				}
				if ("MsgCreate"+strings.Title(typeName) != typeSpec.Name.Name) && ("Msg"+strings.Title(typeName) != typeSpec.Name.Name) {
					return true
				}
				isCreated = true
				return false
			})
		}
	}
	return
}

// isMsgServerDefined checks if the module uses the MsgServer convention for transactions
// this is checked by verifying the existence of the tx.proto file
func isMsgServerDefined(appPath, moduleName string) (bool, error) {
	txProto, err := filepath.Abs(filepath.Join(appPath, "proto", moduleName, "tx.proto"))
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(txProto); os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func isGoReservedWord(name string) bool {
	// Check keyword or literal
	if token.Lookup(name).IsKeyword() {
		return true
	}

	// Check with builtin identifier
	switch name {
	case
		"panic",
		"recover",
		"append",
		"bool",
		"byte",
		"cap",
		"close",
		"complex",
		"complex64",
		"complex128",
		"uint16",
		"copy",
		"false",
		"float32",
		"float64",
		"imag",
		"int",
		"int8",
		"int16",
		"uint32",
		"int32",
		"int64",
		"iota",
		"len",
		"make",
		"new",
		"nil",
		"uint64",
		"print",
		"println",
		"real",
		"string",
		"true",
		"uint",
		"uint8",
		"uintptr":
		return true
	}
	return false
}
