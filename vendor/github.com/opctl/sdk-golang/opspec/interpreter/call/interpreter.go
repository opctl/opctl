package call

import (
	"fmt"

	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/container"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/predicates"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

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
		loopInterpreter:          loop.NewInterpreter(),
		predicatesInterpreter:    predicates.NewInterpreter(),
		opCallInterpreter:        op.NewInterpreter(dataDirPath),
	}
}

type _interpreter struct {
	containerCallInterpreter container.Interpreter
	loopInterpreter          loop.Interpreter
	opCallInterpreter        op.Interpreter
	predicatesInterpreter    predicates.Interpreter
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

	if len(scg.If) > 0 {
		var dcgIf bool
		dcgIf, err = itp.predicatesInterpreter.Interpret(
			opHandle,
			scg.If,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcg.If = &dcgIf
	}

	if nil != scg.Loop {
		dcg.Loop, err = itp.loopInterpreter.Interpret(
			opHandle,
			scg.Loop,
			scope,
		)
		if nil != err {
			return nil, err
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
	case len(scg.Parallel) > 0:
		dcg.Parallel = scg.Parallel
		return dcg, nil
	case len(scg.Serial) > 0:
		dcg.Serial = scg.Serial
		return dcg, nil
	default:
		return nil, fmt.Errorf("Invalid call graph %+v\n", scg)
	}
}
