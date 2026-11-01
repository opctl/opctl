package opstubs

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/fakes"
)

func TestRefreshPartialSuccess(t *testing.T) {
	root := t.TempDir()
	for name, text := range map[string]string{"go.mod": "module app\ngo 1.25\n", "op.go": "package main\nimport (\n_ \"app/opctlgen/good\" //opctl:op ./good\n_ \"app/opctlgen/bad\" //opctl:op repo#v1\n_ \"app/opctlgen/undeclared\"\n)\n"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	api := new(fakes.FakeNode)
	api.DescribeOpStub = func(_ context.Context, req model.DescribeOpReq) (*model.OpSpec, error) {
		if req.Op.Ref == "repo#v1" {
			return nil, errors.New("unreachable")
		}
		return &model.OpSpec{Name: "good", Outputs: map[string]*model.ParamSpec{"text": {String: &model.StringParamSpec{}}}}, nil
	}
	err := RefreshModule(context.Background(), root, api)
	if err == nil || !strings.Contains(err.Error(), "bad") {
		t.Fatalf("%v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "opctlgen/good/stub.go")); err != nil {
		t.Fatal(err)
	}
	var warning bytes.Buffer
	WarnRefresh(context.Background(), root, api, "", &warning)
	if strings.Count(warning.String(), "warning:") != 1 || !strings.Contains(warning.String(), "opctl op gen") {
		t.Fatalf("%s", warning.String())
	}
	count := api.DescribeOpCallCount()
	WarnRefresh(context.Background(), "repo#v1", api, "", &warning)
	WarnRefresh(context.Background(), root, api, root, &warning)
	if api.DescribeOpCallCount() != count {
		t.Fatal("refreshed a remote or cached op")
	}
}

func TestRefreshStubFreeModule(t *testing.T) {
	root := t.TempDir()
	for name, text := range map[string]string{"go.mod": "module app\ngo 1.25\n", "op.go": "package main\n"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	api := new(fakes.FakeNode)
	if err := RefreshModule(context.Background(), root, api); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(root, "opctlgen")); !os.IsNotExist(err) {
		t.Fatalf("created opctlgen without stubs: %v", err)
	}
	if api.DescribeOpCallCount() != 0 {
		t.Fatal("described a stub-free module")
	}
}
