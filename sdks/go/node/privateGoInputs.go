package node

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/goopctx"
	"github.com/opctl/opctl/sdks/go/model"
)

type privateGoInputsKey struct{}

func withPrivateGoInputs(ctx context.Context, params map[string]*model.ParamSpec) context.Context {
	ctx = goopctx.WithActive(withPrivateCompletions(ctx))
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
	if privateResultsActive(ctx) {
		snapshot.Serial = safeSpecs(call.Serial)
		snapshot.Parallel = safeSpecs(call.Parallel)
		if call.Op != nil {
			op := *call.Op
			op.ChildCallCallSpec = safeSpec(op.ChildCallCallSpec)
			op.Inputs = map[string]*model.Value{}
			for name, value := range call.Op.Inputs {
				if !sensitiveExpression(ctx, spec.Op.Inputs[name]) && !secretParam(call.Op.InputParams[name]) {
					op.Inputs[name] = value
				}
			}
			snapshot.Op = &op
		}
	}
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

func secretParam(param *model.ParamSpec) bool {
	return param != nil && ((param.String != nil && param.String.IsSecret) || (param.Number != nil && param.Number.IsSecret) || (param.Array != nil && param.Array.IsSecret) || (param.Object != nil && param.Object.IsSecret) || (param.Dir != nil && param.Dir.IsSecret) || (param.File != nil && param.File.IsSecret) || (param.Socket != nil && param.Socket.IsSecret))
}
func sensitiveExpression(ctx context.Context, expr interface{}) bool {
	names, _ := ctx.Value(privateGoInputsKey{}).([]string)
	var visit func(reflect.Value) bool
	visit = func(value reflect.Value) bool {
		if !value.IsValid() {
			return false
		}
		switch value.Kind() {
		case reflect.Interface, reflect.Pointer:
			if !value.IsNil() {
				return visit(value.Elem())
			}
		case reflect.String:
			for _, name := range names {
				if strings.Contains(value.String(), "$("+name+")") {
					return true
				}
			}
		case reflect.Map:
			for _, key := range value.MapKeys() {
				if visit(key) || visit(value.MapIndex(key)) {
					return true
				}
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < value.Len(); i++ {
				if visit(value.Index(i)) {
					return true
				}
			}
		}
		return false
	}
	return visit(reflect.ValueOf(expr))
}
func guardSecretInterop(ctx context.Context, call *model.OpCall, spec *model.OpCallSpec) error {
	if !privateResultsActive(ctx) || call == nil {
		return nil
	}
	for name, expr := range spec.Inputs {
		if sensitiveExpression(ctx, expr) || secretParam(call.InputParams[name]) {
			return fmt.Errorf("secret interop input %q is not supported until private output collection is available", name)
		}
	}
	for name, param := range call.OutputParams {
		if secretParam(param) {
			return fmt.Errorf("secret interop output %q is not supported until private output collection is available", name)
		}
	}
	return nil
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
