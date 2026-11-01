package node

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	opInterpreter "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op"
)

func TestBootstrapAuthFetchedProvenance(t *testing.T) {
	// A closed local port exercises the provider's cached fallback without a network fixture.
	resource, ref := "127.0.0.1/team/repo", "127.0.0.1/team/repo#v1/ops/check"
	data := t.TempDir()
	cache := filepath.Join(data, "ops")
	repo := resource + "#v1"
	path := filepath.Join(cache, ref)
	if err := os.MkdirAll(path, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "op.go"), []byte("package main"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cache, fmt.Sprintf(".%x", sha1.Sum([]byte(repo)))), []byte("cached"), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	creds := &model.Creds{Username: "source-user", Password: "synthetic-source-token"}
	bootstrap := new(fakes.FakeBootstrapper)
	caller := new(fakes.FakeCaller)
	started := make(chan *model.CallSpec, 1)
	caller.CallStub = func(_ context.Context, _ string, _ map[string]*model.Value, spec *model.CallSpec, _ string, _ *string, _ string) (map[string]*model.Value, error) {
		started <- spec
		return nil, nil
	}
	core := core{dataCachePath: cache, caller: caller, bootstrapper: bootstrap}
	req := model.StartOpReqOp{Ref: ref, PullCreds: creds}
	if _, err := core.DescribeOp(ctx, model.DescribeOpReq{Op: req}); err != nil {
		t.Fatal(err)
	}
	_, described, _ := bootstrap.DescribeArgsForCall(0)
	if described.BootstrapSource != resource || described.BootstrapCreds != creds || described.RepositoryRoot != filepath.Join(cache, repo) {
		t.Fatal("describe lost source repository/credential provenance")
	}
	if _, err := core.StartOp(ctx, model.StartOpReq{Op: req}); err != nil {
		t.Fatal(err)
	}
	var root *model.CallSpec
	select {
	case root = <-started:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if root.Op.BootstrapSource != resource || root.Op.BootstrapCreds != creds {
		t.Fatal("start lost explicit source credentials")
	}
	interpreted, err := opInterpreter.Interpret(ctx, nil, root.Op, "root", path, data)
	if err != nil {
		t.Fatal(err)
	}
	if interpreted.BootstrapSource != resource || interpreted.BootstrapCreds != creds {
		t.Fatal("root interpretation lost bootstrap provenance")
	}
	nested, err := opInterpreter.Interpret(ctx, nil, &model.OpCallSpec{Ref: ref, PullCreds: &model.CredsSpec{Username: "nested", Password: "nested-token"}}, "child", path, data)
	if err != nil {
		t.Fatal(err)
	}
	if nested.BootstrapSource != resource || nested.BootstrapCreds == nil || nested.BootstrapCreds.Username != "nested" {
		t.Fatal("nested call did not establish its own source credentials")
	}
	local, err := opInterpreter.Interpret(ctx, nil, &model.OpCallSpec{Ref: path}, "local", path, data)
	if err != nil {
		t.Fatal(err)
	}
	if local.BootstrapSource != "" || local.BootstrapCreds != nil {
		t.Fatal("local child inherited an unrelated repository's credentials")
	}
}
