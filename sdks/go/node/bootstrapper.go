package node

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

//counterfeiter:generate -o internal/fakes/bootstrapper.go . bootstrapper
type bootstrapper interface {
	Describe(ctx context.Context, opCall *model.OpCall, rootCallID string) (*model.OpSpec, error)
	Plan(ctx context.Context, opCall *model.OpCall, rootCallID string, inputs map[string]interface{}) (*goop.Plan, error)
}

func newBootstrapper(caller caller, dataDirPath string, caches ...*goBootstrapCache) bootstrapper {
	bootstrap := _bootstrapper{caller: caller, dataDirPath: dataDirPath}
	if len(caches) > 0 {
		bootstrap.cache = caches[0]
	}
	return bootstrap
}

type _bootstrapper struct {
	cache       *goBootstrapCache
	caller      caller
	dataDirPath string
}

func (bootstrap _bootstrapper) Describe(ctx context.Context, opCall *model.OpCall, rootCallID string) (*model.OpSpec, error) {
	data, err := bootstrap.run(ctx, opCall, rootCallID, "describe", nil)
	if err != nil {
		return nil, err
	}
	spec, err := goop.ParseDescribe(data)
	if err != nil {
		return nil, &OpDeclarationError{Err: err}
	}
	return spec, nil
}
func (bootstrap _bootstrapper) Plan(ctx context.Context, opCall *model.OpCall, rootCallID string, inputs map[string]interface{}) (*goop.Plan, error) {
	data, err := bootstrap.run(ctx, opCall, rootCallID, "plan", inputs)
	if err != nil {
		return nil, err
	}
	return goop.ParsePlan(data)
}
func (bootstrap _bootstrapper) run(ctx context.Context, opCall *model.OpCall, rootCallID, phase string, inputs map[string]interface{}) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	spec, mountRoot, err := goop.ContainerSpecInRepository(opCall.OpPath, phase, opCall.RepositoryRoot)
	if err != nil {
		return nil, &OpDeclarationError{Err: err}
	}
	if bootstrap.cache != nil {
		ref := spec.Image.Ref
		image := &model.ContainerCallImage{Ref: &ref}
		if auth := bootstrap.cache.stateStore.TryGetAuth(ref); auth != nil {
			image.PullCreds = &auth.Creds
		}
		identity, err := bootstrap.cache.runtime.ResolveImage(ctx, image, rootCallID, discardBootstrapEvents{})
		if err != nil {
			return nil, err
		}
		spec.Image.Ref, spec.Image.Prepared = identity.Ref, true
	}
	id, err := uniquestring.Construct()
	if err != nil {
		return nil, err
	}
	scratch := filepath.Join(bootstrap.dataDirPath, "dcg", id)
	if err := unsudo.CreateDir(scratch); err != nil {
		return nil, err
	}
	defer os.RemoveAll(scratch)
	request, err := json.Marshal(goop.NewRequest(phase, inputs))
	if err != nil {
		return nil, err
	}
	requestPath := filepath.Join(scratch, "request.json")
	if err := os.WriteFile(requestPath, request, 0600); err != nil {
		return nil, err
	}
	scope := map[string]*model.Value{"module": {Dir: &mountRoot}, "request": {File: &requestPath}}
	var stdout bytes.Buffer
	name := "bootstrap " + phase
	var parent *string
	if opCall.OpID != "" {
		parent = &opCall.OpID
	}
	if rootCallID == "" {
		rootCallID = id
	}
	ctx = withStdoutCapture(ctx, &stdout)
	_, err = bootstrap.caller.Call(ctx, id, scope, &model.CallSpec{Container: spec, Name: &name}, opCall.OpPath, parent, rootCallID)
	releaseCompletion(ctx, id)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if err != nil {
		if diagnostic := strings.TrimSpace(string(stdoutCaptureFrom(ctx).diagnostics)); diagnostic != "" {
			return nil, fmt.Errorf("Go op bootstrap %s: %w\n%s", phase, err, diagnostic)
		}
		return nil, fmt.Errorf("Go op bootstrap %s: %w", phase, err)
	}
	return stdout.Bytes(), nil
}
