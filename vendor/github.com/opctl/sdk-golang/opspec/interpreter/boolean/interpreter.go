package boolean

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
)

type Interpreter interface {
	// Interpret interprets an expression to an boolean value.
	// Expression must be a type supported by coerce.ToBoolean
	Interpret(
		scope map[string]*model.Value,
		expression interface{},
		opHandle model.DataHandle,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		coerce:       coerce.New(),
		interpolater: interpolater.New(),
	}
}

type _interpreter struct {
	coerce       coerce.Coerce
	interpolater interpolater.Interpolater
}

func (ea _interpreter) Interpret(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case bool:
		return &model.Value{Boolean: &expression}, nil
	case string:
		var value *model.Value
		if ref, ok := interpreter.TryResolveExplicitRef(expression, scope); ok {
			value = ref
		} else {
			stringValue, err := ea.interpolater.Interpolate(
				expression,
				scope,
				opHandle,
			)
			if nil != err {
				return nil, err
			}
			value = &model.Value{String: &stringValue}
		}
		return ea.coerce.ToBoolean(value)
	}
	return nil, fmt.Errorf("unable to interpretuate %+v to boolean; unsupported type", expression)
}
