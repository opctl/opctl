package expression

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type evalToFile interface {
	// EvalToFile evaluates an expression to a file value
	// expression must be a type supported by data.CoerceToFile
	// scratchDir will be used as the containing dir if file creation necessary
	EvalToFile(
		scope map[string]*model.Value,
		expression interface{},
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
	expression interface{},
	pkgHandle model.PkgHandle,
	scratchDir string,
) (*model.Value, error) {
	var value *model.Value
	var err error

	switch expression := expression.(type) {
	case float64:
		value = &model.Value{Number: &expression}
	case map[string]interface{}:
		value = &model.Value{Object: expression}
	case string:
		if value, err = itp.interpolater.Interpolate(
			expression,
			scope,
			pkgHandle,
		); nil != err {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unable to evaluate %+v to file; unsupported type", expression)
	}

	return itp.data.CoerceToFile(value, scratchDir)
}
