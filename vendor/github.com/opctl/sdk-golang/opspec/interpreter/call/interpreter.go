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
	rootOpID string,
) (*model.DCG, error) {
	var dcgIf *bool
	if len(scg.If) > 0 {
		dcgIfBool, err := itp.predicatesInterpreter.Interpret(
			opHandle,
			scg.If,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgIf = &dcgIfBool
	}

	var dcgLoop *model.DCGLoop
	if nil != scg.Loop {
		var err error
		dcgLoop, err = itp.loopInterpreter.Interpret(
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
		dcgContainerCall, err := itp.containerCallInterpreter.Interpret(
			scope,
			scg.Container,
			id,
			rootOpID,
			opHandle,
		)
		if nil != err {
			return nil, err
		}

		return &model.DCG{
			Container: dcgContainerCall,
			If:        dcgIf,
			Loop:      dcgLoop,
		}, nil
	case nil != scg.Op:
		dcgOpCall, err := itp.opCallInterpreter.Interpret(
			scope,
			scg.Op,
			id,
			opHandle,
			rootOpID,
		)
		if nil != err {
			return nil, err
		}

		return &model.DCG{
			Op:   dcgOpCall,
			If:   dcgIf,
			Loop: dcgLoop,
		}, nil
	case len(scg.Parallel) > 0:
		return &model.DCG{
			Parallel: scg.Parallel,
			If:       dcgIf,
			Loop:     dcgLoop,
		}, nil
	case len(scg.Serial) > 0:
		return &model.DCG{
			Serial: scg.Serial,
			If:     dcgIf,
			Loop:   dcgLoop,
		}, nil
	default:
		return nil, fmt.Errorf("Invalid call graph %+v\n", scg)
	}
}
