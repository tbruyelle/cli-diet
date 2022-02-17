package scaffolder

import (
	"context"
	"fmt"

	"github.com/gobuffalo/genny"

	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/gocmd"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/ibc"
)

const (
	bandImport  = "github.com/bandprotocol/bandchain-packet"
	bandVersion = "v0.0.2"
)

// OracleOption configures options for AddOracle.
type OracleOption func(*oracleOptions)

type oracleOptions struct {
	signer string
}

// newOracleOptions returns a oracleOptions with default options
func newOracleOptions() oracleOptions {
	return oracleOptions{
		signer: "creator",
	}
}

// OracleWithSigner provides a custom signer name for the message
func OracleWithSigner(signer string) OracleOption {
	return func(m *oracleOptions) {
		m.signer = signer
	}
}

// AddOracle adds a new BandChain oracle envtest.
func (s *Scaffolder) AddOracle(
	tracer *placeholder.Tracer,
	moduleName,
	queryName string,
	options ...OracleOption,
) (sm xgenny.SourceModification, err error) {
	if err := s.installBandPacket(); err != nil {
		return sm, err
	}

	o := newOracleOptions()
	for _, apply := range options {
		apply(&o)
	}

	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.LowerCase

	name, err := multiformatname.NewName(queryName)
	if err != nil {
		return sm, err
	}

	if err := checkComponentValidity(s.path, moduleName, name, false); err != nil {
		return sm, err
	}

	mfSigner, err := multiformatname.NewName(o.signer, checkForbiddenOracleFieldName)
	if err != nil {
		return sm, err
	}

	// Module must implement IBC
	ok, err := isIBCModule(s.path, moduleName)
	if err != nil {
		return sm, err
	}
	if !ok {
		return sm, fmt.Errorf("the module %s doesn't implement IBC module interface", moduleName)
	}

	// Generate the packet
	var (
		g    *genny.Generator
		opts = &ibc.OracleOptions{
			AppName:    s.modpath.Package,
			AppPath:    s.path,
			ModulePath: s.modpath.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(s.modpath.RawPath),
			QueryName:  name,
			MsgSigner:  mfSigner,
		}
	)
	g, err = ibc.NewOracle(tracer, opts)
	if err != nil {
		return sm, err
	}
	sm, err = xgenny.RunWithValidation(tracer, g)
	if err != nil {
		return sm, err
	}
	return sm, finish(opts.AppPath, s.modpath.RawPath)
}

func (s Scaffolder) installBandPacket() error {
	return cmdrunner.New().
		Run(context.Background(),
			step.New(step.Exec(gocmd.Name(), "get", gocmd.PackageLiteral(bandImport, bandVersion))),
		)
}
