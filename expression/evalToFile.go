package expression

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalToFile interface {
	// EvalToFile evaluates an expression to a file value
	// scratchDir will be used as the containing dir if file creation necessary
	EvalToFile(
		scope map[string]*model.Value,
		expression string,
		pkgHandle model.PkgHandle,
		scratchDir string,
	) (*model.Value, error)
}

func newEvalToFile() evalToFile {
	return _evalToFile{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalToFile struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (itp _evalToFile) EvalToFile(
	scope map[string]*model.Value,
	expression string,
	pkgHandle model.PkgHandle,
	scratchDir string,
) (*model.Value, error) {
	value, err := itp.interpolater.Interpolate(
		expression,
		scope,
		pkgHandle,
	)

	if nil != err {
		return nil, err
	}

	return itp.data.CoerceToFile(value, scratchDir)
}
