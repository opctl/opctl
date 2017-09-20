package expression

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

type evalToDir interface {
	// EvalToDir evaluates an expression to a dir value
	//
	// Examples of valid dir expressions:
	// scope ref: $(scope-ref)
	// scope ref w/ path expansion: $(scope-ref/sub-dir)
	// scope ref w/ deprecated path expansion: $(scope-ref)/sub-dir
	// pkg fs ref: $(/pkg-fs-ref)
	// pkg fs ref w/ path expansion: $(/pkg-fs-ref/sub-dir)
	EvalToDir(
		scope map[string]*model.Value,
		expression string,
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func newEvalToDir() evalToDir {
	return _evalToDir{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _evalToDir struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (etd _evalToDir) EvalToDir(
	scope map[string]*model.Value,
	expression string,
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	if ref, ok := tryResolveExplicitRef(expression, scope); ok {
		return ref, nil
	}

	possibleRefCloserIndex := strings.Index(expression, string(interpolater.RefCloser))
	var dir *model.Value

	if strings.HasPrefix(expression, string(interpolater.Operator+interpolater.RefOpener)) && possibleRefCloserIndex > 0 {
		possibleRef := expression[2:possibleRefCloserIndex]
		interpolatedPossibleRef, err := etd.interpolater.Interpolate(possibleRef, scope, pkgHandle)
		if nil != err {
			return nil, err
		}

		if dcgValue, ok := scope[interpolatedPossibleRef]; ok && nil != dcgValue.Dir {
			// scope ref w/ deprecated path expansion
			dir = dcgValue
			// trim initial dir ref & interpolate remaining expression
			expression = expression[possibleRefCloserIndex+1:]
		}
	}

	result, err := etd.interpolater.Interpolate(expression, scope, pkgHandle)
	if nil != err {
		return nil, err
	}

	expandedPath := filepath.Join(*dir.Dir, result)

	return &model.Value{Dir: &expandedPath}, nil

}
