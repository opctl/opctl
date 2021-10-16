package core

import (
	"context"
	"path/filepath"
	"regexp"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/outputs"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

//counterfeiter:generate -o internal/fakes/opCaller.go . opCaller
type opCaller interface {
	// Executes an op call
	Call(
		ctx context.Context,
		opCall *model.OpCall,
		inboundScope map[string]*model.Value,
		parentCallID *string,
		rootCallID string,
		opCallSpec *model.OpCallSpec,
	) (
		map[string]*model.Value,
		error,
	)
}

func newOpCaller(
	caller caller,
	dataDirPath string,
) opCaller {
	return _opCaller{
		caller:         caller,
		callScratchDir: filepath.Join(dataDirPath, "call"),
	}
}

type _opCaller struct {
	caller         caller
	callScratchDir string
}

func (oc _opCaller) Call(
	ctx context.Context,
	opCall *model.OpCall,
	inboundScope map[string]*model.Value,
	parentCallID *string,
	rootCallID string,
	opCallSpec *model.OpCallSpec,
) (
	map[string]*model.Value,
	error,
) {
	// form scope for op call by combining defined inputs & op dir
	opCallScope := map[string]*model.Value{}
	for varName, varData := range opCall.Inputs {
		opCallScope[varName] = varData
	}
	// add deprecated absolute path to scope
	opCallScope["/"] = &model.Value{
		Dir: &opCall.OpPath,
	}
	// add current directory to scope
	opCallScope["./"] = &model.Value{
		Dir: &opCall.OpPath,
	}
	// add parent directory to scope
	parentDirPath := filepath.Dir(opCall.OpPath)
	opCallScope["../"] = &model.Value{
		Dir: &parentDirPath,
	}

	opOutputs, err := oc.caller.Call(
		ctx,
		opCall.ChildCallID,
		opCallScope,
		opCall.ChildCallCallSpec,
		opCall.OpPath,
		&opCall.OpID,
		rootCallID,
	)
	if err != nil {
		return nil, err
	}

	opFile, err := opfile.Get(ctx, opCall.OpPath)
	if err != nil {
		return nil, err
	}

	opOutputs, err = outputs.Interpret(
		opOutputs,
		opFile.Outputs,
		opCallSpec.Outputs,
		opCall.OpPath,
		filepath.Join(oc.callScratchDir, opCall.OpID),
	)
	if nil != err {
		return nil, err
	}

	outboundScope := map[string]*model.Value{}

	// filter op outboundScope to bound call outboundScope
	for boundName, boundValue := range opCallSpec.Outputs {
		// return bound outboundScope
		if boundValue == "" {
			// implicit value
			boundValue = boundName
		} else if !regexp.MustCompile(`^\$\(.+\)$`).MatchString(boundValue) {
			// handle obsolete syntax by swapping order
			prevBoundName := boundName
			boundName = boundValue
			boundValue = prevBoundName
		} else {
			boundValue = opspec.RefToName(boundValue)
		}
		for opOutputName, opOutputValue := range opOutputs {
			if boundName == opOutputName {
				outboundScope[boundValue] = opOutputValue
			}
		}
	}

	return outboundScope, nil
}
