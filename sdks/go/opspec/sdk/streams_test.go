package sdk

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func streamGraph() Call {
	producer, version := Container("alpine:3").Cmd("sh", "-ce", "printf 'v1.2.3\\n'; printf 'diagnostic\\n' >&2").Stdout()
	producer, diagnostic := producer.Stderr()
	return Serial(producer, Container("alpine:3").EnvFrom("VERSION", version).EnvFrom("DIAGNOSTIC", diagnostic).Cmd("sh", "-ce", `test "$VERSION" = v1.2.3 && test "$DIAGNOSTIC" = diagnostic`))
}
func TestStreamGolden(t *testing.T) {
	fixturePath, _ := filepath.Abs("testdata/plan.streams.json")
	helloDir(t)
	data, err := Plan(func(struct{}) Call { return streamGraph() }, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(expected) {
		t.Fatalf("stream fixture mismatch: %s", data)
	}
	again, err := Plan(func(struct{}) Call { return streamGraph() }, nil)
	if err != nil || string(data) != string(again) {
		t.Fatal("unstable plan", err)
	}
	var response response
	if err = json.Unmarshal(data, &response); err != nil || response.OpspecSchema != 2 {
		t.Fatal("wrong stream schema", err)
	}
}
func TestStreamVisibility(t *testing.T) {
	for _, tc := range []struct {
		name  string
		build func(ContainerCall, String) Call
		bad   bool
	}{
		{"serial", func(p ContainerCall, v String) Call { return Serial(p, Container("x").EnvFrom("V", v)) }, false},
		{"several readers", func(p ContainerCall, v String) Call {
			return Serial(p, Parallel(Container("x").EnvFrom("V", v), Container("x").EnvFrom("V", v)))
		}, false},
		{"post parallel", func(p ContainerCall, v String) Call {
			return Serial(Parallel(p, Container("x")), Container("x").EnvFrom("V", v))
		}, false},
		{"omitted", func(p ContainerCall, v String) Call { return Container("x").EnvFrom("V", v) }, true},
		{"before", func(p ContainerCall, v String) Call { return Serial(Container("x").EnvFrom("V", v), p) }, true},
		{"self", func(p ContainerCall, v String) Call { return p.EnvFrom("V", v) }, true},
		{"parallel sibling", func(p ContainerCall, v String) Call { return Parallel(p, Container("x").EnvFrom("V", v)) }, true},
		{"twice", func(p ContainerCall, v String) Call { return Serial(p, p) }, true},
		{"copied producer", func(p ContainerCall, v String) Call { return Serial(p, p.Cmd("other")) }, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Plan(func(struct{}) Call { p, v := Container("x").Stdout(); return tc.build(p, v) }, nil)
			if (err != nil) != tc.bad {
				t.Fatalf("err=%v expected failure=%v", err, tc.bad)
			}
		})
	}
}
func TestStreamCopyReuseAndRead(t *testing.T) {
	_, err := Plan(func(struct{}) Call {
		base := Container("x")
		a, v := base.Stdout()
		again, v2 := a.Stdout()
		b, w := base.Stdout()
		if v.handle() != v2.handle() || w.handle() == v.handle() {
			t.Fatal("bad accessor reuse")
		}
		return Serial(again, b, Container("x").EnvFrom("A", v).EnvFrom("B", w))
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Plan(func(struct{}) Call { p, v := Container("x").Stderr(); v.Value(); return p }, nil)
	if err == nil || !strings.Contains(err.Error(), "call outputs can't be read during plan") {
		t.Fatal(err)
	}
	p, v := Container("x").Stdout()
	_, err = Plan(func(struct{}) Call { return p.EnvFrom("V", v) }, nil)
	if err == nil {
		t.Fatal("accepted accessor outside plan")
	}
}
func TestCallResultCreationSite(t *testing.T) {
	_, err := Plan(func(struct{}) Call { _, v := Container("x").Stdout(); return Container("x").EnvFrom("V", v) }, nil)
	if err == nil || !strings.Contains(err.Error(), "streams_test.go:") {
		t.Fatalf("diagnostic lost author creation site: %v", err)
	}
}

func TestStreamExportVisibility(t *testing.T) {
	type Outputs struct {
		Version String `opctl:"version,secret"`
	}
	data, err := PlanWithOutputs(func(struct{}) (Call, Outputs) { p, v := Container("x").Stdout(); return p, Outputs{v} }, nil)
	if err != nil || !strings.Contains(string(data), `"version":"$(opctl_string_1)"`) {
		t.Fatalf("%s %v", data, err)
	}
	_, err = PlanWithOutputs(func(struct{}) (Call, Outputs) { _, v := Container("x").Stdout(); return Container("x"), Outputs{v} }, nil)
	if err == nil {
		t.Fatal("exported invisible stream")
	}
}

func TestStreamMainHelper(t *testing.T) {
	if os.Getenv("OPCTL_STREAM_HELPER") != "1" {
		return
	}
	os.Args = []string{"hello", os.Getenv("OPCTL_STREAM_PHASE")}
	Main(func(struct{}) Call {
		if os.Getenv("OPCTL_TYPED_HELPER") == "1" {
			return typedFixtureGraph()
		}
		if os.Getenv("OPCTL_NAMED_HELPER") == "1" {
			return namedGraph()
		}
		return streamGraph()
	})
	os.Exit(0)
}
func TestStreamSchemaRequests(t *testing.T) {
	for _, tc := range []struct {
		phase    string
		min, max int
		bad      bool
	}{{"plan", 1, 1, true}, {"plan", 1, 2, false}, {"plan", 2, 2, false}, {"plan", 3, 3, true}, {"describe", 1, 1, false}, {"describe", 2, 2, true}} {
		t.Run(fmt.Sprintf("%s-%d-%d", tc.phase, tc.min, tc.max), func(t *testing.T) {
			request := fmt.Sprintf(`{"contractVersion":1,"phase":%q,"opspecSchema":{"min":%d,"max":%d},"inputs":{},"knownOutputs":{}}`, tc.phase, tc.min, tc.max)
			path := filepath.Join(t.TempDir(), "request.json")
			if err := os.WriteFile(path, []byte(request), 0600); err != nil {
				t.Fatal(err)
			}
			command := exec.Command(os.Args[0], "-test.run=^TestStreamMainHelper$")
			command.Env = append(os.Environ(), "OPCTL_STREAM_HELPER=1", "OPCTL_STREAM_PHASE="+tc.phase, "OPCTL_BOOTSTRAP_REQUEST="+path)
			var stderr strings.Builder
			command.Stderr = &stderr
			out, err := command.Output()
			if (err != nil) != tc.bad {
				t.Fatalf("unexpected result %s %s %v", out, stderr.String(), err)
			}
			if tc.bad && (len(out) != 0 || !strings.Contains(stderr.String(), "upgrade opctl")) {
				t.Fatalf("unsupported schema emitted a response: %s %s", out, stderr.String())
			}
		})
	}
}
