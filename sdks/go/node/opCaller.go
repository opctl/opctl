package node

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/internal/goopctx"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/inputs"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/outputs"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/sensitivity"
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
		scope map[string]*model.Value,
	) (
		map[string]*model.Value,
		error,
	)
}

func newOpCaller(
	caller caller,
	dataDirPath string,
	caches ...*goBootstrapCache,
) opCaller {
	return _opCaller{
		caller:       caller,
		bootstrapper: newBootstrapper(caller, dataDirPath, caches...),
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
	scope map[string]*model.Value,
) (
	map[string]*model.Value,
	error,
) {
	if opCall.RepositoryRoot != "" {
		ctx = goopctx.WithRepositoryRoot(ctx, opCall.RepositoryRoot)
	}
	var err error
	outboundScope := map[string]*model.Value{}

	opDir, resolveErr := fs.New(filepath.Join(oc.dataDirPath, "ops")).TryResolve(ctx, opCall.OpPath)
	var isGoOp bool
	var desc *model.OpSpec
	var outputBindings map[string]string
	effectiveOutputs := map[string]string{}
	for name, ref := range opCallSpec.Outputs {
		effectiveOutputs[name] = ref
	}
	childSpec := opCall.ChildCallCallSpec
	if resolveErr == nil {
		isGoOp, err = opfile.IsGoOp(ctx, opDir)
		if err != nil {
			return outboundScope, err
		}
		if isGoOp {
			var session *goBootstrapSession
			ctx, session = withGoBootstrapSession(ctx, opCall.OpPath)
			defer session.close()
			desc, err = oc.bootstrapper.Describe(ctx, opCall, rootCallID)
			if err != nil {
				return outboundScope, err
			}
			// Keep the published CallStarted snapshot free of resolved inputs.
			resolved := *opCall
			opCall = &resolved
			protected := sensitivity.Inputs(opCallSpec.Inputs, desc.Inputs, scope)
			opCall.Inputs, err = inputs.Interpret(opCallSpec.Inputs, desc.Inputs, opCall.OpPath, scope, filepath.Join(oc.dataDirPath, "dcg", opCall.OpID), protected)
			if err != nil {
				return outboundScope, err
			}
			for name, param := range desc.Inputs {
				value := opCall.Inputs[name]
				if value != nil && value.Sensitive && !secretParam(param) && param.Dir == nil && param.File == nil && param.Socket == nil {
					return outboundScope, fmt.Errorf("Go input %q must be declared secret when bound to a sensitive value", name)
				}
			}
			public, err := goop.PlanInputs(desc.Inputs, opCall.Inputs)
			if err != nil {
				return outboundScope, err
			}
			plan, err := oc.bootstrapper.Plan(ctx, opCall, rootCallID, public)
			if err != nil {
				return outboundScope, err
			}
			if err := plan.ValidateOutputs(desc.Outputs); err != nil {
				return outboundScope, err
			}
			outputBindings = plan.OutputBindings
			childSpec = plan.Run
			ctx = withPrivateGoInputs(ctx, desc.Inputs)
			// Output declarations protect their producer before execution.
			sources := map[string]*model.ParamSpec{}
			for name, expr := range outputBindings {
				if secretParam(desc.Outputs[name]) {
					sources[opspec.RefToName(expr)] = desc.Outputs[name]
				}
			}
			ctx = withProtectedSources(ctx, sources)
			if parentCallID == nil && opCall.OpID == rootCallID {
				for name := range desc.Outputs {
					if _, exists := effectiveOutputs[name]; !exists {
						effectiveOutputs[name] = "$(" + name + ")"
					}
				}
			}
		}
	}

	if !isGoOp && privateResultsActive(ctx) {
		ctx = withPrivateGoInputs(ctx, nil)
		ctx = withRuntimeScope(ctx, opCall.Inputs)
		ctx = withProtectedSources(ctx, opCall.OutputParams)
	}

	// form scope for op call by combining defined inputs & op dir
	opCallScope := map[string]*model.Value{}
	for varName, varData := range opCall.Inputs {
		opCallScope[varName] = varData
	}
	if isGoOp {
		empty := ""
		opCallScope["opctl_empty_string"] = &model.Value{String: &empty}
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

	opCallScope = privateScope(ctx, opCallScope)
	opOutputs, err := oc.caller.Call(
		ctx,
		opCall.ChildCallID,
		opCallScope,
		childSpec,
		opCall.OpPath,
		&opCall.OpID,
		rootCallID,
	)
	releaseCompletion(ctx, opCall.ChildCallID)
	if err != nil {
		return outboundScope, err
	}

	// Preserve the YAML op's child-call error precedence when resolution fails.
	if resolveErr != nil {
		return nil, resolveErr
	}
	opFile := desc
	if opFile == nil {
		opFile = &model.OpSpec{}
	}
	if !isGoOp {
		opFile, err = opfile.Get(ctx, opDir)
		if err != nil {
			return outboundScope, err
		}
	}

	if isGoOp {
		final := map[string]*model.Value{}
		for name, value := range opCallScope {
			final[name] = value
		}
		for name, value := range opOutputs {
			final[name] = value
		}
		bound := map[string]*model.Value{}
		names := make([]string, 0, len(outputBindings))
		for name := range outputBindings {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			value, bindErr := reference.Interpret(outputBindings[name], final, nil)
			if bindErr != nil {
				return outboundScope, fmt.Errorf("unable to bind Go output %q", name)
			}
			bound[name] = sensitivity.Mark(value, secretParam(desc.Outputs[name]) || sensitivity.References(outputBindings[name], final))
		}
		opOutputs = bound
	} else if privateResultsActive(ctx) {
		protected := map[string]*model.Value{}
		for name, value := range opOutputs {
			protected[name] = sensitivity.Mark(value, secretParam(opFile.Outputs[name]))
		}
		opOutputs = protected
	}
	opOutputs, err = outputs.Interpret(
		opOutputs,
		opFile.Outputs,
		effectiveOutputs,
		opCall.OpPath,
		filepath.Join(oc.dataDirPath, "call", opCall.OpID),
	)

	// filter op outboundScope to bound call outboundScope
	if err != nil {
		if privateResultsActive(ctx) {
			for name, param := range opFile.Outputs {
				value := opOutputs[name]
				if secretParam(param) || (value != nil && value.Sensitive) {
					return outboundScope, fmt.Errorf("sensitive output %q failed validation", name)
				}
			}
		}
		return outboundScope, err
	}
	for boundName, boundValue := range effectiveOutputs {
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
				outboundScope[boundValue] = sensitivity.Mark(opOutputValue, privateResultsActive(ctx) && secretParam(opFile.Outputs[boundName]))
			}
		}
	}

	return outboundScope, err
}
