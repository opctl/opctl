package goop

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteStubs(t *testing.T) {
	root := t.TempDir()
	generated := []byte(GeneratedStubHeader + "repo#v1. DO NOT EDIT.\npackage greet\n")
	if err := WriteStubs(root, map[string][]byte{"greet": generated, "failed": nil}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "opctlgen/greet/stub.go")
	if got, err := os.ReadFile(path); err != nil || string(got) != string(generated) {
		t.Fatalf("%s %v", got, err)
	}
	ignore, err := os.ReadFile(filepath.Join(root, "opctlgen/.gitignore"))
	if err != nil || string(ignore) != "*\n" {
		t.Fatalf("%s %v", ignore, err)
	}
	if err := WriteStubs(root, map[string][]byte{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("stale generated package retained")
	}
}

func TestWriteStubsPreservesUserFiles(t *testing.T) {
	for _, kind := range []string{"file", "symlink", "stub", "extra", "package-symlink"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "opctlgen")
			generated := []byte(GeneratedStubHeader + "ref. DO NOT EDIT.\npackage greet\n")
			switch kind {
			case "file":
				if err := os.WriteFile(dir, []byte("keep"), 0644); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				if err := os.Symlink(t.TempDir(), dir); err != nil {
					t.Fatal(err)
				}
			default:
				if err := os.MkdirAll(filepath.Join(dir, "greet"), 0755); err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(dir, "greet/stub.go")
				if kind == "package-symlink" {
					if err := os.Symlink(t.TempDir(), filepath.Join(dir, "other")); err != nil {
						t.Fatal(err)
					}
				}
				if err := os.WriteFile(path, []byte("keep"), 0644); err != nil {
					t.Fatal(err)
				}
				if kind == "extra" {
					if err := os.WriteFile(path, generated, 0644); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(dir, "greet/user.txt"), []byte("keep"), 0644); err != nil {
						t.Fatal(err)
					}
				}
			}
			files := map[string][]byte{"greet": generated}
			if kind == "extra" {
				files = nil
			}
			if err := WriteStubs(root, files); err == nil {
				t.Fatal("accepted a user file or symlink")
			}
			if kind == "stub" {
				got, err := os.ReadFile(filepath.Join(dir, "greet/stub.go"))
				if err != nil || string(got) != "keep" {
					t.Fatal("user file overwritten")
				}
			}
		})
	}
}
