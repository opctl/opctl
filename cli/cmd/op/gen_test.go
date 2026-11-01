package op

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/fakes"
)

func TestGenCommand(t *testing.T) {
	root := t.TempDir()
	opDir := filepath.Join(root, ".opspec/release")
	if err := os.MkdirAll(opDir, 0755); err != nil {
		t.Fatal(err)
	}
	for path, text := range map[string]string{filepath.Join(root, "go.mod"): "module app\ngo 1.25\n", filepath.Join(opDir, "op.go"): "package main\nimport _ \"app/opctlgen/greet\" //opctl:op ../greet\n"} {
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	previous, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(previous)
	api := new(fakes.FakeNode)
	api.DescribeOpReturns(&model.OpSpec{Name: "greet"}, nil)
	saved := node
	node = api
	defer func() { node = saved }()
	cmd := newGenCmd()
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "opctlgen/greet/stub.go")); err != nil {
		t.Fatal(err)
	}
	cmd = newGenCmd()
	cmd.SetArgs([]string{"repo#v1"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "local op refs") {
		t.Fatalf("%v", err)
	}
}
