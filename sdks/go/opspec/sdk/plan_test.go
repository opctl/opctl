package sdk

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func helloDir(t *testing.T) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "hello")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatal(err)
		}
	})
}

func TestPlan(t *testing.T) {
	expected, err := os.ReadFile("testdata/plan.hello.json")
	if err != nil {
		t.Fatal(err)
	}
	helloDir(t)
	actual, err := Plan(func(struct{}) Call { return Container("alpine:3").Cmd("echo", "hello").Env("A", "b").WorkDir("/") }, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Fatalf("plan = %s; want %s", actual, expected)
	}
}

func TestPlanRejectsMissingCall(t *testing.T) {
	if _, err := Plan(func(struct{}) Call { return Call{} }, nil); err == nil {
		t.Fatal("missing container call accepted")
	}
}

func TestContainerCallsCanBeReused(t *testing.T) {
	base := Container("alpine:3").Env("A", "b")
	_ = base.Env("A", "changed")
	data, err := Plan(func(struct{}) Call { return base }, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"A":"b"`)) {
		t.Fatalf("derived call mutated base: %s", data)
	}
}

func TestMainHelper(t *testing.T) {
	if os.Getenv("OPCTL_SDK_TEST_HELPER") != "1" {
		return
	}
	os.Args = []string{"hello", os.Getenv("OPCTL_SDK_TEST_PHASE")}
	Main(func(struct{}) Call { return Container("alpine:3").Cmd("echo", "hello").Env("A", "b").WorkDir("/") })
	os.Exit(0)
}

func TestMainContract(t *testing.T) {
	for _, tc := range []struct{ name, phase, request, wantErr string }{
		{"plan", "plan", `{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":1},"inputs":{},"knownOutputs":{}}`, ""},
		{"describe", "describe", `{"contractVersion":1,"phase":"describe","opspecSchema":{"min":1,"max":1},"inputs":{},"knownOutputs":{}}`, ""},
		{"unknown phase", "other", `{"contractVersion":1,"phase":"other","opspecSchema":{"min":1,"max":1}}`, "unsupported bootstrap phase"},
		{"phase mismatch", "plan", `{"contractVersion":1,"phase":"describe","opspecSchema":{"min":1,"max":1}}`, "phase"},
		{"contract version", "plan", `{"contractVersion":2,"phase":"plan","opspecSchema":{"min":1,"max":1}}`, "contract version"},
		{"schema version", "plan", `{"contractVersion":1,"phase":"plan","opspecSchema":{"min":2,"max":2}}`, "upgrade opctl"},
		{"invalid JSON", "plan", `{`, "request"},
		{"inputs", "plan", `{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":1},"inputs":{"A":"b"}}`, "unknown input"},
		{"known outputs", "plan", `{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":1},"knownOutputs":{"A":"b"}}`, "not supported yet"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			helloDir(t)
			request := filepath.Join(t.TempDir(), "request.json")
			if err := os.WriteFile(request, []byte(tc.request), 0600); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(os.Args[0], "-test.run=^TestMainHelper$")
			cmd.Env = append(os.Environ(), "OPCTL_SDK_TEST_HELPER=1", "OPCTL_SDK_TEST_PHASE="+tc.phase, "OPCTL_BOOTSTRAP_REQUEST="+request)
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err := cmd.Run()
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(stderr.String(), tc.wantErr) || stdout.Len() != 0 {
					t.Fatalf("err=%v stdout=%s stderr=%s", err, &stdout, &stderr)
				}
				return
			}
			if err != nil {
				t.Fatalf("%v: %s", err, &stderr)
			}
			var result struct {
				ContractVersion, OpspecSchema int
				Opspec                        struct {
					Name string
					Run  json.RawMessage
				}
			}
			if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if result.ContractVersion != 1 || result.OpspecSchema != 1 || result.Opspec.Name != "hello" {
				t.Fatalf("response = %s", &stdout)
			}
			if tc.phase == "describe" && len(result.Opspec.Run) != 0 {
				t.Fatalf("describe emitted op graph: %s", &stdout)
			}
			if tc.phase == "plan" && len(result.Opspec.Run) == 0 {
				t.Fatal("plan missing op graph")
			}
		})
	}
}
