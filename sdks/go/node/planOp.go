package node

import (
	"context"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/inputs"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/opctl/opctl/sdks/go/opspec/snapshot"
	"path/filepath"
)

func (c core) PlanOp(ctx context.Context, req model.PlanOpReq) (*model.OpPlan, error) {
	dir, err := c.ResolveData(ctx, req.Op.Ref, req.Op.PullCreds)
	if err != nil {
		return nil, err
	}
	spec, err := describeOpDir(ctx, dir, req.Op.PullCreds, "", "", c.bootstrapper, c.dataCachePath)
	if err != nil {
		return nil, err
	}
	planParams := goop.PlanParams(spec.Inputs)
	args := map[string]any{}
	for name, v := range req.Args {
		if planParams[name] == nil {
			return nil, &OpDeclarationError{Err: fmt.Errorf("input %s isn't passed to plan", name)}
		}
		if v == nil {
			return nil, &OpDeclarationError{Err: fmt.Errorf("input %s has no value", name)}
		}
		args[name] = *v
	}
	path := dir.Ref()
	if !filepath.IsAbs(path) {
		path = filepath.Join(c.dataCachePath, path)
	}
	resolved, err := inputs.Interpret(args, planParams, path, nil, "")
	if err != nil {
		return nil, &OpDeclarationError{Err: err}
	}
	public, err := goop.PlanInputs(planParams, resolved)
	if err != nil {
		return nil, &OpDeclarationError{Err: err}
	}
	isGo, err := opfile.IsGoOp(ctx, dir)
	if err != nil {
		return nil, err
	}
	if !isGo {
		expanded, err := snapshot.Expand(spec, public)
		if err != nil {
			return nil, &OpDeclarationError{Err: err}
		}
		return &model.OpPlan{OpSpec: expanded}, nil
	}
	call := &model.OpCall{BaseCall: model.BaseCall{OpPath: path}}
	if fetched, ok := dir.(interface{ RepositoryRoot() string }); ok {
		call.RepositoryRoot = fetched.RepositoryRoot()
	}
	if source, ok := dir.(interface{ SourceResource() string }); ok {
		call.BootstrapSource, call.BootstrapCreds = source.SourceResource(), req.Op.PullCreds
	}
	plan, err := c.bootstrapper.Plan(ctx, call, "", public)
	if err != nil {
		return nil, err
	}
	if err = plan.ValidateOutputs(spec.Outputs); err != nil {
		return nil, &OpDeclarationError{Err: err}
	}
	result := *spec
	result.Run = plan.Run
	return &model.OpPlan{OpSpec: &result, OutputBindings: plan.OutputBindings}, nil
}
