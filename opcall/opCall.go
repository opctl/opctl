// Package opcall implements usecases surrounding op calls
package opcall

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ OpCall

import (
	"github.com/golang-interfaces/satori-go.uuid"
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/inputs"
	"github.com/opspec-io/sdk-golang/pkg"
	"path/filepath"
)

type OpCall interface {
	// Interpret interprets an SCGOpCall into a DCGOpCall
	Interpret(
		scope map[string]*model.Value,
		scgOpCall *model.SCGOpCall,
		opId string,
		parentPkgHandle model.PkgHandle,
		rootOpId string,
	) (*model.DCGOpCall, error)
}

func New(
	rootFSPath string,
) OpCall {
	pkgCachePath := filepath.Join(rootFSPath, "pkgs")
	return _OpCall{
		expression:   expression.New(),
		pkg:          pkg.New(),
		pkgCachePath: pkgCachePath,
		uuid:         iuuid.New(),
		inputs:       inputs.New(),
	}
}

type _OpCall struct {
	expression   expression.Expression
	pkg          pkg.Pkg
	pkgCachePath string
	uuid         iuuid.IUUID
	inputs       inputs.Inputs
}
