package goop

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func stubTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, text := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestScanStubs(t *testing.T) {
	root := stubTree(t, map[string]string{
		"go.mod":                  "module example.com/app\ngo 1.25\n",
		".opspec/a/op.go":         "package main\nimport _ \"example.com/app/helper\"\n",
		"helper/h.go":             "//go:build ignored\n\npackage helper\nimport _ \"example.com/app/opctlgen/greet\"\n",
		".opspec/b/op.go":         "package main\nimport _ \"example.com/app/opctlgen/greet\" //opctl:op ../greet\n",
		".opspec/b/op_test.go":    "package main\nimport _ \"example.com/app/opctlgen/test\" //opctl:op repo#v1/path\n",
		"nested/go.mod":           "module example.com/nested\ngo 1.25\n",
		"nested/bad.go":           "this is ignored",
		"opctlgen/greet/stale.go": "this is ignored",
	})
	source, err := Discover(filepath.Join(root, ".opspec/a"), "")
	if err != nil {
		t.Fatal(err)
	}
	uses, err := ScanStubs(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(uses) != 1 || uses[0].Name != "greet" || uses[0].ModuleRef != ".opspec/greet" || !uses[0].Local {
		t.Fatalf("%+v", uses)
	}
	all, err := ScanModuleStubs(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatalf("%+v", all)
	}
}

func TestScanStubErrors(t *testing.T) {
	for _, tc := range []struct{ source, want string }{
		{"import _ \"app/opctlgen/nope\"", `no ref is declared for op stub "nope"`},
		{"import _ \"app/opctlgen/greet\" //opctl:op", "needs one ref"},
		{"import _ \"app/opctlgen/greet\" //opctl:op ../a ../b", "needs one ref"},
		{"import _ \"app/opctlgen/greet\" // opctl:op ../greet", "without a space"},
		{"//opctl:op ../greet\nimport _ \"app/opctlgen/greet\"", "after an op stub import"},
		{"import _ \"fmt\" //opctl:op ../greet", "after an op stub import"},
		{"import _ \"app/opctlgen/my-sql\" //opctl:op ../greet", "lowercase Go package name"},
		{"import _ \"app/opctlgen/internal\" //opctl:op ../greet", "lowercase Go package name"},
		{"import _ \"app/opctlgen/for\" //opctl:op ../greet", "lowercase Go package name"},
		{"import _ \"app/opctlgen/a/b\" //opctl:op ../greet", "one level"},
		{"import _ \"app/opctlgen/greet\" //opctl:op greet", "use ./ or ../"},
		{"import _ \"app/opctlgen/greet\" //opctl:op /absolute", "use ./ or ../"},
		{"func f(){ //opctl:op ../greet\n}", "after an op stub import"},
	} {
		t.Run(tc.want+tc.source, func(t *testing.T) {
			root := stubTree(t, map[string]string{"go.mod": "module app\ngo 1.25\n", "op.go": "package main\n" + tc.source + "\n"})
			source, err := Discover(root, "")
			if err != nil {
				t.Fatal(err)
			}
			_, err = ScanStubs(source)
			if err == nil || !strings.Contains(err.Error(), tc.want) || !strings.Contains(err.Error(), "op.go:") {
				t.Fatalf("%v", err)
			}
		})
	}
}

func TestScanStubConflictsAndBounds(t *testing.T) {
	root := stubTree(t, map[string]string{"go.mod": "module app\ngo 1.25\n", "op.go": "package main\nimport _ \"app/opctlgen/greet\" //opctl:op ../outside\n", "helper/h.go": "package helper\nimport _ \"app/opctlgen/greet\" //opctl:op repo#v1\n"})
	source, err := Discover(root, "")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ScanStubs(source)
	if err == nil || !strings.Contains(err.Error(), "different refs") || !strings.Contains(err.Error(), "helper/h.go:2") {
		t.Fatalf("%v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "helper/h.go"), []byte("package helper\n"), 0644); err != nil {
		t.Fatal(err)
	}
	source, err = Discover(root, root)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ScanStubs(source)
	if err == nil || !strings.Contains(err.Error(), "outside fetched repository") {
		t.Fatalf("%v", err)
	}
}
