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
		caller:       caller,
		bootstrapper: newBootstrapper(caller, dataDirPath),
		dataDirPath:  dataDirPath,
	}
}

type _opCaller struct {
	bootstrapper bootstrapper
	caller       caller
	dataDirPath  string
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

	opDir, resolveErr := fs.New(filepath.Join(oc.dataDirPath, "ops")).TryResolve(ctx, opCall.OpPath)
	var isGoOp bool
	childSpec := opCall.ChildCallCallSpec
	if resolveErr == nil {
		isGoOp, err = opfile.IsGoOp(ctx, opDir)
		if err != nil {
			return outboundScope, err
		}
		if isGoOp {
			plan, err := oc.bootstrapper.Plan(ctx, opCall, rootCallID)
			if err != nil {
				return outboundScope, err
			}
			childSpec = plan.Run
		}
	}

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
		childSpec,
		opCall.OpPath,
		&opCall.OpID,
		rootCallID,
	)
	if err != nil {
		return outboundScope, err
	}

	// Preserve the YAML op's child-call error precedence when resolution fails.
	if resolveErr != nil {
		return nil, resolveErr
	}
	opFile := &model.OpSpec{}
	if !isGoOp {
		opFile, err = opfile.Get(ctx, opDir)
		if err != nil {
			return outboundScope, err
		}
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
