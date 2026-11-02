package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
)

func TestPrepareCopyPreservesCallers(t *testing.T) {
	root := t.TempDir()
	source, target := filepath.Join(root, "source"), filepath.Join(root, "copy")
	for name, data := range map[string]string{
		"go.mod":                 "module fixture\ngo 1.25\n",
		"caller/op.go":           "package main\n// caller must remain unchanged\n",
		"callee/op.yml":          "name: callee\nrun: {serial: []}\n",
		"opctlgen/stale/stub.go": "stale",
	} {
		path := filepath.Join(source, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	before, err := tree(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := prepareCopy(source, target, "/repo/sdk"); err != nil {
		t.Fatal(err)
	}
	caller, _ := os.ReadFile(filepath.Join(target, "caller/op.go"))
	if string(caller) != "package main\n// caller must remain unchanged\n" {
		t.Fatal("caller changed")
	}
	mod, _ := os.ReadFile(filepath.Join(target, "go.mod"))
	if !strings.Contains(string(mod), "require github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0") || !strings.Contains(string(mod), "=> /repo/sdk") {
		t.Fatalf("%s", mod)
	}
	if _, err := os.Stat(filepath.Join(target, "opctlgen")); !os.IsNotExist(err) {
		t.Fatal("stale stubs copied")
	}
	after, err := tree(source)
	if err != nil || before != after {
		t.Fatal("source changed", err)
	}
}

func TestPlanRejectsExecutionAndSecretEvents(t *testing.T) {
	name := "bootstrap build"
	bootstrap := model.Event{CallStarted: &model.CallStarted{Call: model.Call{Name: &name}}}
	if err := checkPlanEvent(bootstrap); err != nil {
		t.Fatal(err)
	}
	for _, event := range []model.Event{
		{CallStarted: &model.CallStarted{Call: model.Call{Op: &model.OpCall{}}}},
		{ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{Data: []byte("SECRET-MARKER")}},
		{ContainerStdErrWrittenTo: &model.ContainerStdErrWrittenTo{Data: []byte("PRIVATE-DEFAULT")}},
	} {
		if err := checkPlanEvent(event); err == nil {
			t.Fatal("plan accepted execution or a secret event")
		}
	}
}
