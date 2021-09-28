// Package field provides methods to parse a field provided in a command with the format name:type
package field

import (
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

const (
	TypeCustom = "custom"
	TypeString = "string"
	TypeBool   = "bool"
	TypeInt    = "int"
	TypeInt32  = "int32"
	TypeUint   = "uint"
	TypeUint64 = "uint64"

	TypeSeparator = ":"
)

var (
	StaticDataTypes = map[string]string{
		TypeString: TypeString,
		TypeBool:   TypeBool,
		TypeInt:    TypeInt32,
		TypeUint:   TypeUint64,
	}
)

type (
	// Field represents a field inside a structure for a component
	// it can be a field contained in a type or inside the response of a query, etc...
	Field struct {
		Name         multiformatname.Name
		Datatype     string
		DatatypeName string
		Nested       Fields
	}

	// Fields represents a Field slice
	Fields []Field
)

// validateField validate the field name and type, and run the forbidden method
func validateField(field string, isForbiddenField func(string) error) (multiformatname.Name, string, error) {
	fieldSplit := strings.Split(field, TypeSeparator)
	if len(fieldSplit) > 2 {
		return multiformatname.Name{}, "", fmt.Errorf("invalid field format: %s, should be 'name' or 'name:type'", field)
	}

	name, err := multiformatname.NewName(fieldSplit[0])
	if err != nil {
		return name, "", err

	}

	// Ensure the field name is not a Go reserved name, it would generate an incorrect code
	if err := isForbiddenField(name.LowerCamel); err != nil {
		return name, "", fmt.Errorf("%s can't be used as a field name: %s", name, err.Error())
	}

	// Check if the object has an explicit type. The default is a string
	dataTypeName := TypeString
	isTypeSpecified := len(fieldSplit) == 2
	if isTypeSpecified {
		dataTypeName = fieldSplit[1]
	}
	return name, dataTypeName, nil
}

// ParseFields parses the provided fields, analyses the types
// and checks there is no duplicated field
func ParseFields(
	fields []string,
	isForbiddenField func(string) error,
) (Fields, error) {
	// Used to check duplicated field
	existingFields := make(map[string]bool)

	var parsedFields Fields
	for _, field := range fields {
		name, datatypeName, err := validateField(field, isForbiddenField)
		if err != nil {
			return parsedFields, err
		}

		// Ensure the field is not duplicated
		if _, exists := existingFields[name.LowerCamel]; exists {
			return parsedFields, fmt.Errorf("the field %s is duplicated", name.Original)
		}
		existingFields[name.LowerCamel] = true

		// Check if is a static type
		if datatype, ok := StaticDataTypes[datatypeName]; ok {
			parsedFields = append(parsedFields, Field{
				Name:         name,
				Datatype:     datatype,
				DatatypeName: datatypeName,
			})
			continue
		}

		parsedFields = append(parsedFields, Field{
			Name:         name,
			Datatype:     datatypeName,
			DatatypeName: TypeCustom,
		})
	}
	return parsedFields, nil
}

// GetDatatype return the Datatype based in the DatatypeName
func (f Field) GetDatatype() string {
	switch f.DatatypeName {
	case TypeString, TypeBool, TypeInt, TypeUint:
		return f.Datatype
	case TypeCustom:
		return fmt.Sprintf("*%s", f.Datatype)
	default:
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
}

// GetProtoDatatype return the proto Datatype based in the DatatypeName
func (f Field) GetProtoDatatype() string {
	switch f.DatatypeName {
	case TypeString, TypeBool, TypeInt, TypeUint, TypeCustom:
		return f.Datatype
	default:
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
}

// NeedCastImport return true if the field slice
// needs import the cast library
func (f Fields) NeedCastImport() bool {
	for _, field := range f {
		if field.DatatypeName != TypeString &&
			field.DatatypeName != TypeCustom {
			return true
		}
	}
	return false
}

// IsComplex return true if the field slice
// needs import the json library
func (f Fields) IsComplex() bool {
	for _, field := range f {
		if field.DatatypeName == TypeCustom {
			return true
		}
	}
	return false
}

// String return all inline fields args for command usage
func (f Fields) String() string {
	args := ""
	for _, field := range f {
		args += fmt.Sprintf(" [%s]", field.Name.Kebab)
	}
	return args
}

// Custom return a list of custom fields
func (f Fields) Custom() []string {
	fields := make([]string, 0)
	for _, field := range f {
		if field.DatatypeName == TypeCustom {
			dataType, err := multiformatname.NewName(field.Datatype)
			if err != nil {
				panic(err)
			}
			fields = append(fields, dataType.Snake)
		}
	}
	return fields
}
