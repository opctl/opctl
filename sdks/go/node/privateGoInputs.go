package node

import (
	"context"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
)

type privateGoInputsKey struct{}

func withPrivateGoInputs(ctx context.Context, params map[string]*model.ParamSpec) context.Context {
	ctx = withPrivateCompletions(ctx)
	activatePrivateResults(ctx)
	names := []string{}
	for name, param := range params {
		if (param.String != nil && param.String.IsSecret) || (param.Number != nil && param.Number.IsSecret) || (param.Array != nil && param.Array.IsSecret) || (param.Object != nil && param.Object.IsSecret) || (param.Dir != nil && param.Dir.IsSecret) || (param.File != nil && param.File.IsSecret) || (param.Socket != nil && param.Socket.IsSecret) {
			names = append(names, name)
		}
	}
	return context.WithValue(ctx, privateGoInputsKey{}, names)
}

func callForEvent(ctx context.Context, call *model.Call, spec *model.CallSpec) model.Call {
	snapshot := *call
	names, _ := ctx.Value(privateGoInputsKey{}).([]string)
	if len(names) == 0 || call.Container == nil || spec.Container == nil {
		return snapshot
	}
	container := *call.Container
	bindings := func(actual map[string]string, raw interface{}) map[string]string {
		expressions := map[string]string{}
		switch values := raw.(type) {
		case map[string]interface{}:
			for key, value := range values {
				expressions[key], _ = value.(string)
			}
		case map[string]string:
			expressions = values
		}
		safe := map[string]string{}
		for key, value := range actual {
			safe[key] = value
			expression := expressions[key]
			for _, name := range names {
				if strings.Contains(expression, "$("+name+")") {
					safe[key] = expression
					break
				}
			}
		}
		return safe
	}
	container.EnvVars = bindings(container.EnvVars, spec.Container.EnvVars)
	container.Dirs = bindings(container.Dirs, spec.Container.Dirs)
	container.Files = bindings(container.Files, spec.Container.Files)
	container.Sockets = bindings(container.Sockets, spec.Container.Sockets)
	snapshot.Container = &container
	return snapshot
}

func outputsForEvent(ctx context.Context, outputs map[string]*model.Value) map[string]*model.Value {
	names, _ := ctx.Value(privateGoInputsKey{}).([]string)
	if len(names) == 0 || outputs == nil {
		return outputs
	}
	safe := make(map[string]*model.Value, len(outputs))
	for name, value := range outputs {
		safe[name] = value
	}
	for _, name := range names {
		delete(safe, name)
	}
	return safe
}
