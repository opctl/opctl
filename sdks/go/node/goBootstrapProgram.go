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
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

func (bootstrap _bootstrapper) runCached(ctx context.Context, op *model.OpCall, root, phase string, inputs map[string]interface{}) ([]byte, error) {
	session, _ := ctx.Value(goBootstrapSessionKey{}).(*goBootstrapSession)
	if session == nil || session.path != op.OpPath {
		ctx, session = withGoBootstrapSession(ctx, op.OpPath)
		defer session.close()
	}
	if session.program == nil {
		program, release, err := bootstrap.cache.prepare(ctx, bootstrap, op, root)
		if err != nil {
			return nil, err
		}
		session.program, session.release = program, release
	}
	program := session.program
	spec := program.source.ContainerSpec(phase)
	spec.Image.Ref, spec.Image.Prepared = program.image.Ref, true
	executable := filepath.Join(program.path, "program")
	spec.Files["/opctl/program"] = model.Value{File: &executable}
	spec.ReadOnlyMounts["/opctl/program"] = true
	spec.Cmd = []interface{}{"sh", "/opctl/launcher", "/opctl/program", phase}
	invoke := func() ([]byte, error) {
		return bootstrap.invokeProgram(ctx, op, root, "bootstrap "+phase, spec, program.source.ModuleRoot, false, goop.NewRequest(phase, inputs))
	}
	if phase == "describe" {
		return bootstrap.describeCached(ctx, op, root, program, invoke)
	}
	return invoke()
}

func (bootstrap _bootstrapper) buildProgram(ctx context.Context, op *model.OpCall, root string, source *goop.Source, image *containerruntime.ImageIdentity, staging string) error {
	spec := source.ContainerSpec("plan")
	spec.Image.Ref, spec.Image.Prepared = image.Ref, true
	spec.Dirs["/opctl/build"] = model.Value{Dir: &staging}
	delete(spec.Files, "/opctl/request.json")
	env := spec.EnvVars.(map[string]interface{})
	delete(env, "OPCTL_BOOTSTRAP_REQUEST")
	env["GOARCH"], env["GOAMD64"] = image.Arch, "v1"
	auths := []model.Auth{}
	if bootstrap.cache.stateStore != nil {
		var err error
		auths, err = bootstrap.cache.stateStore.AuthSnapshot()
		if err != nil {
			return err
		}
	}
	auth, err := prepareBootstrapAuth(auths, op.BootstrapSource, op.BootstrapCreds)
	if err != nil {
		return err
	}
	if auth.private != "" {
		env["GOPRIVATE"], env["GONOPROXY"], env["GONOSUMDB"] = auth.private, auth.private, auth.private
	}
	if len(auth.archive) > 0 {
		ctx = containerruntime.WithBootstrapInput(ctx, bytes.NewReader(auth.archive))
	}
	spec.Cmd = []interface{}{"sh", "-ce", "go version > /opctl/build/go-version; go build " + goop.BuildFlags + " -o /opctl/build/program ."}
	_, err = bootstrap.invokeProgram(ctx, op, root, "bootstrap build", spec, source.ModuleRoot, true, nil)
	return err
}

func (bootstrap _bootstrapper) invokeProgram(ctx context.Context, op *model.OpCall, root, name string, spec *model.ContainerCallSpec, moduleRoot string, shared bool, request interface{}) ([]byte, error) {
	id, err := uniquestring.Construct()
	if err != nil {
		return nil, err
	}
	scratch := filepath.Join(bootstrap.dataDirPath, "dcg", id)
	if err := unsudo.CreateDir(scratch); err != nil {
		return nil, err
	}
	defer os.RemoveAll(scratch)
	scope := map[string]*model.Value{"module": {Dir: &moduleRoot}}
	if request != nil {
		data, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(scratch, "request.json")
		if err := unsudo.CreateFile(path, data); err != nil {
			return nil, err
		}
		if err := os.Chmod(path, 0600); err != nil {
			return nil, err
		}
		scope["request"] = &model.Value{File: &path}
	}
	if shared {
		err = prepareGoCaches(spec, scope, bootstrap.dataDirPath, scratch)
	} else {
		err = prepareGoLauncher(spec, scope, scratch)
	}
	if err != nil {
		return nil, err
	}
	var stdout bytes.Buffer
	ctx = withStdoutCapture(ctx, &stdout)
	var parent *string
	if op.OpID != "" {
		parent = &op.OpID
	}
	if root == "" {
		root = id
	}
	_, err = bootstrap.caller.Call(ctx, id, scope, &model.CallSpec{Container: spec, Name: &name}, op.OpPath, parent, root)
	releaseCompletion(ctx, id)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if err != nil {
		if diagnostic := strings.TrimSpace(string(stdoutCaptureFrom(ctx).diagnostics)); diagnostic != "" {
			return nil, fmt.Errorf("Go op %s: %w\n%s", name, err, diagnostic)
		}
		return nil, fmt.Errorf("Go op %s: %w", name, err)
	}
	return stdout.Bytes(), nil
}
