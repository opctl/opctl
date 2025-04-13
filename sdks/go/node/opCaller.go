package node

import (
	"context"
	"path/filepath"
	"regexp"

	"github.com/opctl/opctl/sdks/go/data/fs"
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
		caller:      caller,
		dataDirPath: dataDirPath,
	}
}

type _opCaller struct {
	caller      caller
	dataDirPath string
}

func (oc _opCaller) Call(
	ctx context.Context,
	opCall *model.OpCall,
	parentCallID *string,
	rootCallID string,
	opCallSpec *model.OpCallSpec,
) (
	map[string]*model.Value,
	error,
) {
	var err error
	outboundScope := map[string]*model.Value{}

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
		return outboundScope, err
	}

	opDir, err := fs.New(filepath.Join(oc.dataDirPath, "ops")).TryResolve(ctx, opCall.OpPath)
	if err != nil {
		return nil, err
	}

	var opFile *model.OpSpec
	opFile, err = opfile.Get(
		ctx,
		opDir,
	)
	if err != nil {
		return outboundScope, err
	}
	opOutputs, err = outputs.Interpret(
		opOutputs,
		opFile.Outputs,
		opCallSpec.Outputs,
		opCall.OpPath,
		filepath.Join(oc.dataDirPath, "call", opCall.OpID),
	)

	// filter op outboundScope to bound call outboundScope
	for boundName, boundValue := range opCallSpec.Outputs {
		// return bound outboundScope
		if boundValue == "" {
			// implicit value
			boundValue = boundName
		} else if !regexp.MustCompile("^\\$\\(.+\\)$").MatchString(boundValue) {
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

	return outboundScope, err
}
