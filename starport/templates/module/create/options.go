package modulecreate

import (
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/templates/field"
)

// CreateOptions represents the options to scaffold a Cosmos SDK module
type CreateOptions struct {
	ModuleName string
	ModulePath string
	AppName    string
	AppPath    string
	OwnerName  string
	Params     field.Fields

	// True if the module should implement the IBC module interface
	IsIBC bool

	// Channel ordering of the IBC module: ordered, unordered or none
	IBCOrdering string

	// Dependencies of the module
	Dependencies []Dependency
}

// MsgServerOptions defines options to add MsgServer
type MsgServerOptions struct {
	ModuleName string
	ModulePath string
	AppName    string
	AppPath    string
	OwnerName  string
}

// Validate that options are usable
func (opts *CreateOptions) Validate() error {
	return nil
}

// Dependency represents a module dependency of a module
type Dependency struct {
	Name       string
	KeeperName string // KeeperName represents the name of the keeper for the module in app.go
}

// NewDependency returns a new dependency object
func NewDependency(name, keeperName string) Dependency {
	// Default keeper name
	if keeperName == "" {
		keeperName = fmt.Sprintf("%sKeeper", strings.Title(name))
	}
	return Dependency{
		name,
		keeperName,
	}
}
