package sdk

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type pathInputs struct {
	Home      Dir    `opctl:"HOME"`
	Src       Dir    `opctl:"src"`
	CachePath String `opctl:"cachePath"`
}

func pathCall(in pathInputs) Call {
	return Container("alpine:3").
		Dir("/src", OpDir("../..")).
		Dir("/scripts", OpDir("./scripts")).
		Dir("/nested", OpDir("../..").Dir("sub/dir")).
		File("/expected.txt", OpFile("./expected.txt")).
		Dir("/go/pkg", in.Home.Dir("go/pkg")).
		Dir("/cache", in.Home.Dir(in.CachePath.Value()).Dir("v1")).
		File("/config.yml", in.Src.File("config/app.yml")).
		Cmd("true")
}

func TestOpDirAndSubpathsGolden(t *testing.T) {
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	helloDir(t)
	data, err := planAtSchema(pathCall, map[string]any{"cachePath": ".cache/go-build"}, 1)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(previous, "testdata", "plan.op-dir.json")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	expected, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, expected) {
		t.Fatalf("plan = %s; want %s", data, expected)
	}
	for _, part := range []string{`"/src":"$(../..)"`, `"/nested":"$(../../sub/dir)"`, `"/expected.txt":"$(./expected.txt)"`, `"/go/pkg":"$(HOME/go/pkg)"`, `"/cache":"$(HOME/.cache/go-build/v1)"`, `"/config.yml":"$(src/config/app.yml)"`} {
		if !bytes.Contains(data, []byte(part)) {
			t.Fatalf("missing %s in %s", part, data)
		}
	}
}

func TestOpDirHandlesPassToOpCalls(t *testing.T) {
	data, err := Plan(func(struct{}) Call {
		return CallOp("../build").Inputs(map[string]any{"src": OpDir("./").Dir("app"), "config": OpFile("../config.yml")})
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{`"src":"$(./app)"`, `"config":"$(../config.yml)"`} {
		if !bytes.Contains(data, []byte(part)) {
			t.Fatalf("missing %s in %s", part, data)
		}
	}
}

func TestOpPathsRejectInvalidPaths(t *testing.T) {
	for _, path := range []string{"", ".", "..", "scripts", "/abs", ".scripts", "./$(x)", "./a(b", "../a)b"} {
		for name, handle := range map[string]Handle{"OpDir": OpDir(path), "OpFile": OpFile(path)} {
			_, err := Plan(func(struct{}) Call { return Container("a").EnvFrom("X", handle) }, nil)
			if err == nil || !strings.HasPrefix(err.Error(), name) {
				t.Fatalf("%s(%q) err = %v", name, path, err)
			}
		}
	}
}

func TestSubpathsRejectEscapesAndReferences(t *testing.T) {
	for _, path := range []string{"", "/abs", "..", "../x", "a/../b", "a/..", "$(x)", "a$b", "a(b", "a)b"} {
		for name, handle := range map[string]Handle{"Dir": OpDir("./x").Dir(path), "File": OpDir("./x").File(path)} {
			_, err := Plan(func(struct{}) Call { return Container("a").EnvFrom("X", handle) }, nil)
			if err == nil || !strings.Contains(err.Error(), "must be relative") {
				t.Fatalf("%s(%q) err = %v", name, path, err)
			}
		}
	}
	// Input text read during plan can't escape the directory or add references.
	for _, value := range []string{"../../etc", "$(token)", "x)/y"} {
		_, err := Plan(func(in pathInputs) Call { return Container("a").Dir("/x", in.Home.Dir(in.CachePath.Value())) }, map[string]any{"cachePath": value})
		if err == nil || !strings.Contains(err.Error(), "must be relative") {
			t.Fatalf("cachePath %q err = %v", value, err)
		}
	}
}

func TestSubpathsOfCallOutputsAreRejected(t *testing.T) {
	for name, fn := range map[string]func(struct{}) Call{
		"NewDir": func(struct{}) Call {
			out := NewDir()
			return Serial(Container("a").Dir("/out", out), Container("a").Dir("/sub", out.Dir("sub")))
		},
		"NewDir file": func(struct{}) Call {
			out := NewDir()
			return Serial(Container("a").Dir("/out", out), Container("a").File("/sub", out.File("a.txt")))
		},
		"OpOutput": func(struct{}) Call {
			call, out := OpOutput[Dir](CallOp("../x"), "dir")
			return Serial(call, Container("a").Dir("/sub", out.Dir("sub")))
		},
	} {
		_, err := Plan(fn, nil)
		if err == nil || err.Error() != "subpaths of call outputs are not supported yet" {
			t.Fatalf("%s err = %v", name, err)
		}
	}
	_, err := Plan(func(struct{}) Call { return Container("a").Dir("/x", RawRef[Dir]("$(x)").Dir("y")) }, nil)
	if err == nil || !strings.Contains(err.Error(), "subpaths of RawRef handles are not supported") {
		t.Fatalf("RawRef err = %v", err)
	}
	_, err = Plan(func(struct{}) Call { return Container("a").Dir("/x", Dir{}.Dir("y")) }, nil)
	if err == nil || err.Error() != "handle is not bound to an input" {
		t.Fatalf("unbound err = %v", err)
	}
}
