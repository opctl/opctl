package call

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/parallelloop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/serialloop"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	// Interpret interprets an CallSpec into a Call
	Interpret(
		scope map[string]*model.Value,
		callSpec *model.CallSpec,
		id string,
		opPath string,
		parentID *string,
		rootCallID string,
	) (*model.Call, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter(
	containerCallInterpreter container.Interpreter,
	dataDirPath string,
) Interpreter {
	return _interpreter{
		containerCallInterpreter: containerCallInterpreter,
		opCallInterpreter:        op.NewInterpreter(dataDirPath),
		predicatesInterpreter:    predicates.NewInterpreter(),
		parallelLoopInterpreter:  parallelloop.NewInterpreter(),
		serialLoopInterpreter:    serialloop.NewInterpreter(),
	}
}

type _interpreter struct {
	containerCallInterpreter container.Interpreter
	opCallInterpreter        op.Interpreter
	parallelLoopInterpreter  parallelloop.Interpreter
	predicatesInterpreter    predicates.Interpreter
	serialLoopInterpreter    serialloop.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	callSpec *model.CallSpec,
	id string,
	opPath string,
	parentID *string,
	rootCallID string,
) (*model.Call, error) {
	call := &model.Call{
		ID:       id,
		Name:     callSpec.Name,
		Needs:    callSpec.Needs,
		ParentID: parentID,
		RootID:   rootCallID,
	}
	var err error

	if nil != callSpec.If {
		callIf, err := itp.predicatesInterpreter.Interpret(
			*callSpec.If,
			scope,
		)
		if nil != err {
			return nil, err
		}

		call.If = &callIf

		if !callIf {
			// end interpretation early since call will be skipped
			return call, err
		}
	}

	switch {
	case nil != callSpec.Container:
		call.Container, err = itp.containerCallInterpreter.Interpret(
			scope,
			callSpec.Container,
			id,
			opPath,
		)
		return call, err
	case nil != callSpec.Op:
		call.Op, err = itp.opCallInterpreter.Interpret(
			scope,
			callSpec.Op,
			id,
			opPath,
		)
		return call, err
	case nil != callSpec.Parallel:
		call.Parallel = *callSpec.Parallel
		return call, nil
	case nil != callSpec.ParallelLoop:
		call.ParallelLoop, err = itp.parallelLoopInterpreter.Interpret(
			*callSpec.ParallelLoop,
			scope,
		)
		return call, err
	case nil != callSpec.Serial:
		call.Serial = *callSpec.Serial
		return call, nil
	case nil != callSpec.SerialLoop:
		call.SerialLoop, err = itp.serialLoopInterpreter.Interpret(
			*callSpec.SerialLoop,
			scope,
		)
		return call, err
	default:
		return nil, fmt.Errorf("Invalid call graph %+v\n", callSpec)
	}
}
