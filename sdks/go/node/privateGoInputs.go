package node

import (
	"context"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/goopctx"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/sensitivity"
)

type privateGoInputsKey struct{}
type pendingPrivateNamesKey struct{}

func withPrivateGoInputs(ctx context.Context, params map[string]*model.ParamSpec) context.Context {
	ctx = goopctx.WithActive(withPrivateCompletions(ctx))
	activatePrivateResults(ctx)
	names := []string{}
	for name, param := range params {
		if secretParam(param) {
			names = append(names, name)
		}
	}
	ctx = context.WithValue(ctx, pendingPrivateNamesKey{}, names)
	return context.WithValue(ctx, privateGoInputsKey{}, names)
}

func callForEvent(ctx context.Context, call *model.Call, spec *model.CallSpec) (snapshot model.Call) {
	snapshot = *call
	defer func() {
		if snapshot.Op != nil {
			op := *snapshot.Op
			op.InputParams = nil
			op.OutputParams = nil
			snapshot.Op = &op
		}
	}()
	names, _ := ctx.Value(privateGoInputsKey{}).([]string)
	if privateResultsActive(ctx) {
		if spec != nil && spec.Name != nil && sensitiveExpression(ctx, *spec.Name) {
			name := *spec.Name
			snapshot.Name = &name
		}
		snapshot.Serial = safeSpecs(call.Serial)
		snapshot.Parallel = safeSpecs(call.Parallel)
		if call.SerialLoop != nil && spec.SerialLoop != nil {
			loop := *call.SerialLoop
			if sensitiveExpression(ctx, spec.SerialLoop.Range) {
				loop.Range = nil
			}
			loop.Run = callForEvent(ctx, &loop.Run, &spec.SerialLoop.Run)
			snapshot.SerialLoop = &loop
		}
		if call.ParallelLoop != nil && spec.ParallelLoop != nil {
			loop := *call.ParallelLoop
			if sensitiveExpression(ctx, spec.ParallelLoop.Range) {
				loop.Range = nil
			}
			loop.Run = callForEvent(ctx, &loop.Run, &spec.ParallelLoop.Run)
			snapshot.ParallelLoop = &loop
		}

		if call.Op != nil {
			op := *call.Op
			op.ChildCallCallSpec = safeSpec(op.ChildCallCallSpec)
			op.Inputs = map[string]*model.Value{}
			for name, value := range call.Op.Inputs {
				if !value.Sensitive && !sensitiveExpression(ctx, spec.Op.Inputs[name]) && !secretParam(call.Op.InputParams[name]) {
					op.Inputs[name] = value
				}
			}
			snapshot.Op = &op
		}
	}
	if (!privateResultsActive(ctx) && len(names) == 0) || call.Container == nil || spec.Container == nil {
		return snapshot
	}
	container := *call.Container
	if spec.Container.Name != nil && sensitiveExpression(ctx, *spec.Container.Name) {
		name := *spec.Container.Name
		container.Name = &name
	}
	if container.Image != nil {
		image := *container.Image
		image.PullCreds = nil
		if spec.Container.Image != nil && sensitiveExpression(ctx, spec.Container.Image.Ref) {
			ref := spec.Container.Image.Ref
			image.Ref = &ref
			image.Src = nil
		}
		container.Image = &image
	}
	bindings := func(actual map[string]string, raw interface{}) map[string]string {
		expressions := map[string]string{}
		switch values := raw.(type) {
		case string:
			if sensitiveExpression(ctx, values) {
				return nil
			}
		case map[string]interface{}:
			for key, value := range values {
				if sensitiveExpression(ctx, key) {
					return nil
				}
				expressions[key], _ = value.(string)
			}
		case map[string]string:
			for key := range values {
				if sensitiveExpression(ctx, key) {
					return nil
				}
			}
			expressions = values
		}
		safe := map[string]string{}
		for key, value := range actual {
			safe[key] = value
			expression := expressions[key]
			if sensitiveExpression(ctx, expression) {
				safe[key] = expression
			}
		}
		return safe
	}
	if command, ok := spec.Container.Cmd.([]interface{}); ok && len(command) == len(container.Cmd) {
		container.Cmd = append([]string(nil), container.Cmd...)
		for i, expr := range command {
			if sensitiveExpression(ctx, expr) {
				container.Cmd[i], _ = expr.(string)
			}
		}
	}
	if sensitiveExpression(ctx, spec.Container.WorkDir) {
		container.WorkDir = spec.Container.WorkDir
	}
	container.EnvVars = bindings(container.EnvVars, spec.Container.EnvVars)
	container.Dirs = bindings(container.Dirs, spec.Container.Dirs)
	container.Files = bindings(container.Files, spec.Container.Files)
	container.Sockets = bindings(container.Sockets, spec.Container.Sockets)
	snapshot.Container = &container
	return snapshot
}

func outputsForEvent(ctx context.Context, outputs map[string]*model.Value) map[string]*model.Value {
	if outputs == nil {
		return outputs
	}
	safe := make(map[string]*model.Value, len(outputs))
	for name, value := range outputs {
		if value == nil || !value.Sensitive {
			safe[name] = value
		}
	}
	return safe
}

func secretParam(param *model.ParamSpec) bool { return sensitivity.Secret(param) }
func sensitiveExpression(ctx context.Context, expr interface{}) bool {
	names, _ := ctx.Value(privateGoInputsKey{}).([]string)
	scope := map[string]*model.Value{}
	for _, name := range names {
		scope[name] = &model.Value{Sensitive: true}
	}
	return sensitivity.References(expr, scope)
}
func safeSpecs(specs []*model.CallSpec) []*model.CallSpec {
	if specs == nil {
		return nil
	}
	copies := make([]*model.CallSpec, len(specs))
	for i, spec := range specs {
		copies[i] = safeSpec(spec)
	}
	return copies
}
func safeSpec(spec *model.CallSpec) *model.CallSpec {
	if spec == nil {
		return nil
	}
	snapshot := *spec
	if spec.Op != nil {
		op := *spec.Op
		snapshot.Op = &op
		if op.PullCreds != nil {
			creds := *op.PullCreds
			op.PullCreds = &creds
			if !wholeRef(creds.Username) {
				creds.Username = ""
			}
			if !wholeRef(creds.Password) {
				creds.Password = ""
			}
		}
	}
	if spec.Container != nil && spec.Container.Image != nil {
		container := *spec.Container
		image := *container.Image
		if image.PullCreds != nil {
			creds := *image.PullCreds
			if !wholeRef(creds.Username) {
				creds.Username = ""
			}
			if !wholeRef(creds.Password) {
				creds.Password = ""
			}
			image.PullCreds = &creds
		}
		container.Image = &image
		snapshot.Container = &container
	}
	if spec.Serial != nil {
		children := safeSpecs(*spec.Serial)
		snapshot.Serial = &children
	}
	if spec.Parallel != nil {
		children := safeSpecs(*spec.Parallel)
		snapshot.Parallel = &children
	}
	if spec.SerialLoop != nil {
		loop := *spec.SerialLoop
		loop.Run = *safeSpec(&loop.Run)
		snapshot.SerialLoop = &loop
	}
	if spec.ParallelLoop != nil {
		loop := *spec.ParallelLoop
		loop.Run = *safeSpec(&loop.Run)
		snapshot.ParallelLoop = &loop
	}
	return &snapshot
}
func wholeRef(expr string) bool { return strings.HasPrefix(expr, "$(") && strings.HasSuffix(expr, ")") }

type protectedSourcesKey struct{}

func withProtectedSources(ctx context.Context, params map[string]*model.ParamSpec) context.Context {
	sources := map[string]bool{}
	for name, param := range params {
		if secretParam(param) {
			sources[name] = true
		}
	}
	return context.WithValue(ctx, protectedSourcesKey{}, sources)
}
func privateScope(ctx context.Context, scope map[string]*model.Value) map[string]*model.Value {
	names, _ := ctx.Value(pendingPrivateNamesKey{}).([]string)
	sources, _ := ctx.Value(protectedSourcesKey{}).(map[string]bool)
	result := map[string]*model.Value{}
	for name, value := range scope {
		result[name] = sensitivity.Mark(value, sources[name])
	}
	for _, name := range names {
		if value, ok := result[name]; ok {
			result[name] = sensitivity.Mark(value, true)
		}
	}
	return result
}
func withRuntimeScope(ctx context.Context, scope map[string]*model.Value) context.Context {
	names := []string{}
	for name, value := range scope {
		if value != nil && value.Sensitive {
			names = append(names, name)
		}
	}
	sources, _ := ctx.Value(protectedSourcesKey{}).(map[string]bool)
	for name := range sources {
		names = append(names, name)
	}
	ctx = context.WithValue(ctx, pendingPrivateNamesKey{}, []string(nil))
	return context.WithValue(ctx, privateGoInputsKey{}, names)
}
func protectProduced(ctx context.Context, outputs map[string]*model.Value) {
	sources, _ := ctx.Value(protectedSourcesKey{}).(map[string]bool)
	for name, value := range outputs {
		outputs[name] = sensitivity.Mark(value, sources[name])
	}
}

func protectLoopScope(scope map[string]*model.Value, inbound map[string]*model.Value, expr interface{}, vars *model.LoopVarsSpec) map[string]*model.Value {
	if vars == nil || !sensitivity.References(expr, inbound) {
		return scope
	}
	result := map[string]*model.Value{}
	for name, value := range scope {
		result[name] = value
	}
	for _, ref := range []*string{vars.Index, vars.Key, vars.Value} {
		if ref != nil {
			name := strings.TrimSuffix(strings.TrimPrefix(*ref, "$("), ")")
			result[name] = sensitivity.Mark(result[name], true)
		}
	}
	return result
}
