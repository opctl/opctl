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
	//
	// Examples of valid file expressions:
	// scope ref: $(scope-ref)
	// scope ref w/ path expansion: $(scope-ref/file.txt)
	// scope ref w/ deprecated path expansion: $(scope-ref)/file.txt
	// pkg fs ref: $(/pkg-fs-ref)
	// pkg fs ref w/ path expansion: $(/pkg-fs-ref/file.txt)
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

func (etf _evalToFile) EvalToFile(
	scope map[string]*model.Value,
	expression interface{},
	pkgHandle model.PkgHandle,
	scratchDir string,
) (*model.Value, error) {
	var value *model.Value

	switch expression := expression.(type) {
	case float64:
		value = &model.Value{Number: &expression}
	case map[string]interface{}:
		value = &model.Value{Object: expression}
	case string:
		if ref, ok := tryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := etf.interpolater.Interpolate(
				expression,
				scope,
				pkgHandle,
			)
			if nil != err {
				return nil, err
			}

			value = &model.Value{String: &stringValue}
		}
	default:
		return nil, fmt.Errorf("unable to evaluate %+v to file; unsupported type", expression)
	}

	return etf.data.CoerceToFile(value, scratchDir)
}
