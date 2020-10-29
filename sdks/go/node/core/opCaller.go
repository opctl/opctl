package core

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/outputs"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/opCaller.go . opCaller
type opCaller interface {
	// Executes an op call
	Call(
		ctx context.Context,
		opCall *model.OpCall,
		inboundScope map[string]*model.Value,
		parentCallID *string,
		opCallSpec *model.OpCallSpec,
	)
}

func newOpCaller(
	stateStore stateStore,
	pubSub pubsub.PubSub,
	caller caller,
	dataDirPath string,
) opCaller {
	return _opCaller{
		caller:             caller,
		stateStore:         stateStore,
		callScratchDir:     filepath.Join(dataDirPath, "call"),
		outputsInterpreter: outputs.NewInterpreter(),
		opFileGetter:       opfile.NewGetter(),
		pubSub:             pubSub,
	}
}

type _opCaller struct {
	stateStore         stateStore
	caller             caller
	callScratchDir     string
	outputsInterpreter outputs.Interpreter
	opFileGetter       opfile.Getter
	pubSub             pubsub.PubSub
}

func (oc _opCaller) Call(
	ctx context.Context,
	opCall *model.OpCall,
	inboundScope map[string]*model.Value,
	parentCallID *string,
	opCallSpec *model.OpCallSpec,
) {
	var err error
	outboundScope := map[string]*model.Value{}

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		var outcome string
		if call := oc.stateStore.TryGet(opCall.RootCallID); nil != call && call.IsKilled {
			outcome = model.OpOutcomeKilled
		} else if nil != err {
			outcome = model.OpOutcomeFailed
		} else {
			outcome = model.OpOutcomeSucceeded
		}

		event := model.Event{
			Timestamp: time.Now().UTC(),
			CallEnded: &model.CallEnded{
				CallID:     opCall.OpID,
				CallType:   model.CallTypeOp,
				Ref:        opCall.OpPath,
				Outcome:    outcome,
				RootCallID: opCall.RootCallID,
				Outputs:    outboundScope,
			},
		}

		if outcome == model.OpOutcomeFailed {
			event.CallEnded.Error = &model.CallEndedError{
				Message: err.Error(),
			}
		}

		oc.pubSub.Publish(event)

	}()

	callStartedTime := time.Now().UTC()

	// form scope for op call by combining defined inputs & op dir
	opCallScope := map[string]*model.Value{}
	for varName, varData := range opCall.Inputs {
		opCallScope[varName] = varData
	}
	opCallScope["/"] = &model.Value{
		Dir: &opCall.OpPath,
	}

	oc.caller.Call(
		ctx,
		opCall.ChildCallID,
		opCallScope,
		opCall.ChildCallCallSpec,
		opCall.OpPath,
		&opCall.OpID,
		opCall.RootCallID,
	)

	// subscribe to events
	eventChannel, _ := oc.pubSub.Subscribe(
		ctx,
		model.EventFilter{
			Roots: []string{opCall.RootCallID},
			Since: &callStartedTime,
		},
	)

	opOutputs := map[string]*model.Value{}

eventLoop:
	for event := range eventChannel {
		switch {
		case nil != event.CallEnded && event.CallEnded.CallID == opCall.OpID:
			// parent ended prematurely
			return
		case nil != event.CallEnded && event.CallEnded.CallID == opCall.ChildCallID:
			if nil != event.CallEnded.Error {
				err = errors.New(event.CallEnded.Error.Message)
				// end on any error
				return
			}
			for name, value := range event.CallEnded.Outputs {
				opOutputs[name] = value
			}
			break eventLoop
		}
	}

	var opFile *model.OpSpec
	opFile, err = oc.opFileGetter.Get(
		ctx,
		opCall.OpPath,
	)
	if nil != err {
		return
	}
	opOutputs, err = oc.outputsInterpreter.Interpret(
		opOutputs,
		opFile.Outputs,
		opCall.OpPath,
		filepath.Join(oc.callScratchDir, opCall.OpID),
	)

	// filter op outboundScope to bound call outboundScope
	for boundName, boundValue := range opCallSpec.Outputs {
		// return bound outboundScope
		if "" == boundValue {
			// implicit value
			boundValue = boundName
		}
		for opOutputName, opOutputValue := range opOutputs {
			if boundValue == opOutputName {
				outboundScope[boundName] = opOutputValue
			}
		}
	}
}
