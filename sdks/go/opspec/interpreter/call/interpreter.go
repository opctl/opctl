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
	// Interpret interprets an CallSpec into a DCG
	Interpret(
		scope map[string]*model.Value,
		callSpec *model.CallSpec,
		id string,
		opPath string,
		parentID *string,
		rootOpID string,
	) (*model.DCG, error)
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
	rootOpID string,
) (*model.DCG, error) {
	dcg := &model.DCG{
		Id:       id,
		Name:     callSpec.Name,
		Needs:    callSpec.Needs,
		ParentID: parentID,
	}
	var err error

	if nil != callSpec.If {
		dcgIf, err := itp.predicatesInterpreter.Interpret(
			*callSpec.If,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcg.If = &dcgIf

		if !dcgIf {
			// end interpretation early since call will be skipped
			return dcg, err
		}
	}

	switch {
	case nil != callSpec.Container:
		dcg.Container, err = itp.containerCallInterpreter.Interpret(
			scope,
			callSpec.Container,
			id,
			rootOpID,
			opPath,
		)
		return dcg, err
	case nil != callSpec.Op:
		dcg.Op, err = itp.opCallInterpreter.Interpret(
			scope,
			callSpec.Op,
			id,
			opPath,
			rootOpID,
		)
		return dcg, err
	case nil != callSpec.Parallel:
		dcg.Parallel = *callSpec.Parallel
		return dcg, nil
	case nil != callSpec.ParallelLoop:
		dcg.ParallelLoop, err = itp.parallelLoopInterpreter.Interpret(
			*callSpec.ParallelLoop,
			scope,
		)
		return dcg, err
	case nil != callSpec.Serial:
		dcg.Serial = *callSpec.Serial
		return dcg, nil
	case nil != callSpec.SerialLoop:
		dcg.SerialLoop, err = itp.serialLoopInterpreter.Interpret(
			*callSpec.SerialLoop,
			scope,
		)
		return dcg, err
	default:
		return nil, fmt.Errorf("Invalid call graph %+v\n", callSpec)
	}
}
