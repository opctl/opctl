package number

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name fakeInterpreter ./ interpreter

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

type interpreter interface {
	// interprets an expression to a string
	Interpret(
		scope map[string]*model.Value,
		expression string,
	) (float64, error)
}

func newInterpreter() interpreter {
	return _interpreter{
		data:         data.New(),
		interpolater: interpolater.New(),
	}
}

type _interpreter struct {
	data         data.Data
	interpolater interpolater.Interpolater
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	expression string,
) (float64, error) {
	value, err := itp.interpolater.Interpolate(
		expression,
		scope,
	)

	if nil != err {
		return 0, err
	}

	return itp.data.CoerceToNumber(&model.Value{String: &value})
}
