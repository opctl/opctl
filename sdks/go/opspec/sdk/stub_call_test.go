package sdk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStubCall(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module app\ngo 1.21\n"), 0600); err != nil {
		t.Fatal(err)
	}
	previous, _ := os.Getwd()
	defer os.Chdir(previous)
	for _, dir := range []string{".", ".opspec/release", "internal/helper"} {
		cwd := filepath.Join(root, dir)
		if err := os.MkdirAll(cwd, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.Chdir(cwd); err != nil {
			t.Fatal(err)
		}
		for _, target := range []string{".opspec/greet", "../outside"} {
			data, err := Plan(func(struct{}) Call {
				return StubCall("greet", target, true, map[string]Handle{"name": StringOf("Ada"), "optional": String{}}, "name")
			}, nil)
			if err != nil {
				t.Fatal(err)
			}
			rel, _ := filepath.Rel(cwd, filepath.Join(root, target))
			if !strings.HasPrefix(rel, "../") {
				rel = "./" + rel
			}
			if !strings.Contains(string(data), `"ref":"`+filepath.ToSlash(rel)+`"`) || strings.Contains(string(data), "optional") {
				t.Fatalf("%s", data)
			}
		}
	}
	data, err := Plan(func(struct{}) Call {
		return StubCall("git", "example.com/op#v1/path", false, map[string]Handle{"name": RawRef[String]("$(name)")}, "name")
	}, nil)
	if err != nil || !strings.Contains(string(data), "example.com/op#v1/path") {
		t.Fatalf("%s %v", data, err)
	}
	_, err = Plan(func(struct{}) Call {
		return StubCall("greet", "git#v1", false, map[string]Handle{"name": String{}}, "name")
	}, nil)
	if err == nil || !strings.Contains(err.Error(), `op stub greet: input "name" is required`) {
		t.Fatalf("%v", err)
	}
	_, err = Plan(func(struct{}) Call {
		return StubCall("greet", "git#v1", false, map[string]Handle{"name": RawRef[String]("bad")}, "name")
	}, nil)
	if err == nil || !strings.Contains(err.Error(), "RawRef needs") {
		t.Fatalf("%v", err)
	}
}

func TestStubCallMissingModule(t *testing.T) {
	dir := t.TempDir()
	previous, _ := os.Getwd()
	defer os.Chdir(previous)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	_, err := Plan(func(struct{}) Call { return StubCall("greet", "greet", true, nil) }, nil)
	if err == nil || !strings.Contains(err.Error(), "go.mod") {
		t.Fatalf("%v", err)
	}
}
