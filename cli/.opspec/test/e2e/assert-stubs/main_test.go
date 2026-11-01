package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

func TestBuildCounts(t *testing.T) {
	caller := "/proof/caller"
	name := "bootstrap build"
	event := func(path string) model.Event {
		return model.Event{CallStarted: &model.CallStarted{Call: model.Call{Name: &name}, Ref: path}}
	}
	cold := []model.Event{event(caller), event("/proof/gogreet")}
	if err := checkBuilds(cold, caller, 1, 1); err != nil {
		t.Fatal(err)
	}
	if err := checkBuilds(nil, caller, 0, 0); err != nil {
		t.Fatal(err)
	}
	if err := checkBuilds(cold, caller, 1, -1); err != nil {
		t.Fatal(err)
	}
	if err := checkBuilds(cold, caller, 0, 0); err == nil {
		t.Fatal("warm build accepted")
	}
	if err := checkBuilds(append(cold, event(caller)), caller, 1, -1); err == nil {
		t.Fatal("multiple caller builds accepted")
	}
}

func TestFixtureCopiesHaveIndependentColdCaches(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.MkdirAll(filepath.Join(source, "caller"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "go.mod"), []byte("module fixture\ngo 1.25\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "caller/op.go"), []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	keys := []string{}
	for _, name := range []string{"first", "second"} {
		target := filepath.Join(root, name)
		if err := copyFixture(source, target); err != nil {
			t.Fatal(err)
		}
		module, err := goop.Discover(filepath.Join(target, "caller"), "")
		if err != nil {
			t.Fatal(err)
		}
		key, err := module.Fingerprint(context.Background(), nil, goop.Toolchain{}, goop.BuildRecipe{})
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, key)
	}
	if keys[0] == keys[1] {
		t.Fatal("separate proof runs share a bootstrap cache key, so the second cold run can be warm")
	}
}

func TestFixtureCopyAndInterfaceEdits(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	target := filepath.Join(root, "copy")
	for name, text := range map[string]string{"go.mod": "module fixture\ngo 1.25\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => ../sdk\n", "yaml/op.yml": "name: yaml\ninputs:\n  label: {string: {}}\n  src: {dir: {}}\noutputs:\n  text: {string: {}}\n", "opctlgen/greet/stub.go": "stale"} {
		path := filepath.Join(source, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := copyFixture(source, target); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(target, "opctlgen")); !os.IsNotExist(err) {
		t.Fatal("copied editor stubs")
	}
	path := filepath.Join(target, "yaml/op.yml")
	if err := addOptional(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(data), "optionalAdded") {
		t.Fatalf("%s %v", data, err)
	}
	if err := removeLabel(path); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(path)
	if err != nil || strings.Contains(string(data), "  label:") {
		t.Fatalf("%s %v", data, err)
	}
	mod, err := os.ReadFile(filepath.Join(target, "go.mod"))
	if err != nil || !strings.Contains(string(mod), filepath.Join(root, "sdk")) {
		t.Fatalf("%s %v", mod, err)
	}
}
