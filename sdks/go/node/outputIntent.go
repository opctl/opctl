package node

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/sensitivity"
)

type outputIntentKey struct{}

// Keep declarations private even before a nested call activates result protection.
func withOutputIntent(ctx context.Context, params map[string]*model.ParamSpec) context.Context {
	return context.WithValue(ctx, outputIntentKey{}, params)
}
func outputIntent(ctx context.Context) map[string]*model.ParamSpec {
	params, _ := ctx.Value(outputIntentKey{}).(map[string]*model.ParamSpec)
	return params
}
func secretDestination(ctx context.Context, scope map[string]*model.Value, ref string) bool {
	sources, _ := ctx.Value(protectedSourcesKey{}).(map[string]bool)
	return sources[opspec.RefToName(ref)] || sensitiveExpression(ctx, ref) || sensitivity.References(ref, scope) || secretParam(outputIntent(ctx)[opspec.RefToName(ref)])
}
func outputDestinations(ctx context.Context, params map[string]*model.ParamSpec, bindings map[string]string, scope map[string]*model.Value) map[string]*model.ParamSpec {
	result := make(map[string]*model.ParamSpec, len(params))
	for name, param := range params {
		result[name] = param
	}
	for name, ref := range bindings {
		if ref == "" {
			ref = "$(" + name + ")"
		} else if !wholeRef(ref) {
			name, ref = ref, "$("+name+")"
		}
		if secretDestination(ctx, scope, ref) {
			result[name] = &model.ParamSpec{String: &model.StringParamSpec{IsSecret: true}}
		}
	}
	return result
}
