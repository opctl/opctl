package node

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

//counterfeiter:generate -o internal/fakes/bootstrapper.go . bootstrapper
type bootstrapper interface {
	Plan(ctx context.Context, opCall *model.OpCall, rootCallID string) (*model.OpSpec, error)
}

func newBootstrapper(caller caller, dataDirPath string) bootstrapper {
	return _bootstrapper{caller: caller, dataDirPath: dataDirPath}
}

type _bootstrapper struct {
	caller      caller
	dataDirPath string
}

func (bootstrap _bootstrapper) Plan(ctx context.Context, opCall *model.OpCall, rootCallID string) (*model.OpSpec, error) {
	spec, mountRoot, err := goop.ContainerSpec(opCall.OpPath, "plan")
	if err != nil {
		return nil, err
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
	request, err := json.Marshal(goop.NewRequest("plan", nil))
	if err != nil {
		return nil, err
	}
	requestPath := filepath.Join(scratch, "request.json")
	if err := os.WriteFile(requestPath, request, 0600); err != nil {
		return nil, err
	}
	scope := map[string]*model.Value{"module": {Dir: &mountRoot}, "request": {File: &requestPath}}
	var stdout bytes.Buffer
	name := "bootstrap plan"
	_, err = bootstrap.caller.Call(withStdoutCapture(ctx, &stdout), id, scope, &model.CallSpec{Container: spec, Name: &name}, opCall.OpPath, &opCall.OpID, rootCallID)
	if err != nil {
		return nil, fmt.Errorf("Go op bootstrap plan: %w", err)
	}
	return goop.ParsePlan(stdout.Bytes())
}
