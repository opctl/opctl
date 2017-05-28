// Package opcall implements usecases surrounding op calls
package opcall

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ OpCall

import (
	"github.com/golang-interfaces/satori-go.uuid"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/inputs"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"path/filepath"
)

type OpCall interface {
	// Interpret interprets an SCGOpCall into a DCGOpCall
	Interpret(
		scope map[string]*model.Data,
		scgOpCall *model.SCGOpCall,
		opId string,
		pkgBasePath string,
		rootOpId string,
	) (*model.DCGOpCall, error)
}

func New(
	rootFSPath string,
) OpCall {
	return _OpCall{
		interpolater: interpolater.New(),
		manifest:     manifest.New(),
		pkg:          pkg.New(),
		pkgCachePath: filepath.Join(rootFSPath, "pkgs"),
		uuid:         iuuid.New(),
		inputs:       inputs.New(),
	}
}

type _OpCall struct {
	interpolater interpolater.Interpolater
	manifest     manifest.Manifest
	pkg          pkg.Pkg
	pkgCachePath string
	uuid         iuuid.IUUID
	inputs       inputs.Inputs
}
