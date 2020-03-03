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
	// Interpret interprets an SCG into a DCG
	Interpret(
		scope map[string]*model.Value,
		scg *model.SCG,
		id string,
		opHandle model.DataHandle,
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
	scg *model.SCG,
	id string,
	opHandle model.DataHandle,
	parentID *string,
	rootOpID string,
) (*model.DCG, error) {
	dcg := &model.DCG{
		Id:       id,
		ParentID: parentID,
	}
	var err error

	if nil != scg.If {
		dcgIf, err := itp.predicatesInterpreter.Interpret(
			opHandle,
			*scg.If,
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
	case nil != scg.Container:
		dcg.Container, err = itp.containerCallInterpreter.Interpret(
			scope,
			scg.Container,
			id,
			rootOpID,
			opHandle,
		)
		return dcg, err
	case nil != scg.Op:
		dcg.Op, err = itp.opCallInterpreter.Interpret(
			scope,
			scg.Op,
			id,
			opHandle,
			rootOpID,
		)
		return dcg, err
	case nil != scg.Parallel:
		dcg.Parallel = scg.Parallel
		return dcg, nil
	case nil != scg.ParallelLoop:
		dcg.ParallelLoop, err = itp.parallelLoopInterpreter.Interpret(
			opHandle,
			*scg.ParallelLoop,
			scope,
		)
		return dcg, err
	case nil != scg.Serial:
		dcg.Serial = scg.Serial
		return dcg, nil
	case nil != scg.SerialLoop:
		dcg.SerialLoop, err = itp.serialLoopInterpreter.Interpret(
			opHandle,
			*scg.SerialLoop,
			scope,
		)
		return dcg, err
	default:
		return nil, fmt.Errorf("Invalid call graph %+v\n", scg)
	}
}
