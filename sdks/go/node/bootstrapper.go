package node

import (
	"context"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

//counterfeiter:generate -o internal/fakes/bootstrapper.go . bootstrapper
type bootstrapper interface {
	Describe(ctx context.Context, opCall *model.OpCall, rootCallID string) (*model.OpSpec, error)
	Plan(ctx context.Context, opCall *model.OpCall, rootCallID string, inputs map[string]interface{}) (*goop.Plan, error)
}

func newBootstrapper(caller caller, dataDirPath string, cache *goBootstrapCache) bootstrapper {
	return _bootstrapper{caller: caller, dataDirPath: dataDirPath, cache: cache}
}

type _bootstrapper struct {
	cache       *goBootstrapCache
	caller      caller
	dataDirPath string
}

func (bootstrap _bootstrapper) Describe(ctx context.Context, opCall *model.OpCall, rootCallID string) (*model.OpSpec, error) {
	data, err := bootstrap.runCached(ctx, opCall, rootCallID, "describe", nil)
	if err != nil {
		return nil, err
	}
	spec, err := goop.ParseDescribeForRange(data, bootstrap.cache.supportedSchemas())
	if err != nil {
		return nil, &OpDeclarationError{Err: err}
	}
	return spec, nil
}
func (bootstrap _bootstrapper) Plan(ctx context.Context, opCall *model.OpCall, rootCallID string, inputs map[string]interface{}) (*goop.Plan, error) {
	data, err := bootstrap.runCached(ctx, opCall, rootCallID, "plan", inputs)
	if err != nil {
		return nil, err
	}
	return goop.ParsePlanForRange(data, bootstrap.cache.supportedSchemas())
}
