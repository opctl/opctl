package goconvert

import (
	"errors"
	"github.com/opctl/opctl/sdks/go/opspec"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConversionFiles(t *testing.T) {
	for _, mode := range []string{"new", "shared", "own", "bad-own", "app"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "hello")
			if err := os.Mkdir(dir, 0755); err != nil {
				t.Fatal(err)
			}
			mod := "module app\ngo 1.25\n"
			if mode == "shared" || mode == "own" {
				mod += "require github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\n"
			}
			if mode != "new" {
				path := root
				if mode == "own" || mode == "bad-own" {
					path = dir
				}
				if err := os.WriteFile(filepath.Join(path, "go.mod"), []byte(mod), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("keep\n"), 0600); err != nil {
				t.Fatal(err)
			}
			files, err := PrepareFiles(dir, opspec.GoSDK{})
			if mode == "bad-own" {
				if err == nil || !strings.Contains(err.Error(), "without the Go op SDK") {
					t.Fatal(err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			own := mode == "new" || mode == "app"
			if own != (files.Files["go.mod"] != nil) {
				t.Fatal(files)
			}
			if own && string(files.Files[".gitignore"]) != "keep\n/opctlgen/\n" {
				t.Fatal(string(files.Files[".gitignore"]))
			}
		})
	}
}
func TestInterruptedConversionKeepsYAML(t *testing.T) {
	dir := t.TempDir()
	original := []byte("name: original\n")
	if err := os.WriteFile(filepath.Join(dir, "op.yml"), original, 0600); err != nil {
		t.Fatal(err)
	}
	stop := errors.New("interrupted")
	err := writeFiles(dir, map[string][]byte{"op.go": []byte("package main\n")}, func() error { return stop })
	if !errors.Is(err, stop) {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "op.yml"))
	if err != nil || string(data) != string(original) {
		t.Fatal("lost YAML")
	}
	if _, err = os.Stat(filepath.Join(dir, "op.go")); err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(filepath.Join(dir, "op.go.tmp")); !os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func TestFailedSumWriteRemovesOnlyNewModule(t *testing.T) {
	for _, fresh := range []bool{false, true} {
		dir := t.TempDir()
		original := []byte("name: retry\n")
		if err := os.WriteFile(filepath.Join(dir, "op.yml"), original, 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(filepath.Join(dir, "go.sum"), 0755); err != nil {
			t.Fatal(err)
		}
		files := map[string][]byte{"go.sum": []byte("sum"), "op.go": []byte("package main\n")}
		mod := []byte("module ops.local/retry\ngo 1.25\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\n")
		if fresh {
			files["go.mod"] = mod
		} else if err := os.WriteFile(filepath.Join(dir, "go.mod"), mod, 0600); err != nil {
			t.Fatal(err)
		}
		if err := WriteFiles(dir, files); err == nil {
			t.Fatal("expected go.sum write failure")
		}
		data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if fresh && !os.IsNotExist(err) {
			t.Fatalf("left a new module behind: %v", err)
		}
		if !fresh && (err != nil || string(data) != string(mod)) {
			t.Fatalf("changed reused module: %v", err)
		}
		data, err = os.ReadFile(filepath.Join(dir, "op.yml"))
		if err != nil || string(data) != string(original) {
			t.Fatal("lost YAML")
		}
		retry, err := PrepareFiles(dir, opspec.GoSDK{})
		if err != nil {
			t.Fatal(err)
		}
		if retry.NewModule != fresh {
			t.Fatalf("retry reused an incomplete module: %+v", retry)
		}
	}
}
