package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func baselineSchemaGraph() Call {
	file := NewFile()
	return Serial(Parallel(Container("alpine:3").File("/marker", file).Cmd("sh", "-ce", `printf marker > /marker`), CallOp("./opaque").Inputs(map[string]any{"stdout": map[string]any{"outputParams": "literal $(unchanged)"}})), Container("alpine:3").File("/marker", file).Cmd("cat", "/marker"))
}

type schemaTypedOutputs struct {
	Count Number `opctl:"count,minimum=1,format=integer"`
	Token String `opctl:"token,secret,minLength=1"`
}

func schemaGraph(version int) Call {
	switch version {
	case 1:
		return baselineSchemaGraph()
	case 2:
		return streamGraph()
	case 3:
		return namedGraph()
	default:
		p, out := Outputs[schemaTypedOutputs](Container("alpine:3").Cmd("true"))
		return Serial(p, Container("alpine:3").EnvFrom("COUNT", out.Count).EnvFrom("TOKEN", out.Token).Cmd("true"))
	}
}

func TestNegotiatedSchemaGoldens(t *testing.T) {
	root, _ := filepath.Abs("testdata/schema-versions")
	helloDir(t)
	var baseline []byte
	for version := 1; version <= 4; version++ {
		r := schemaRange{1, version}
		contractVersion := 1
		request, _ := json.Marshal(request{ContractVersion: &contractVersion, Phase: "plan", OpspecSchema: r, Inputs: map[string]any{}, KnownOutputs: map[string]json.RawMessage{}})
		plan, err := emitForRange("plan", withoutOutputs(func(struct{}) Call { return schemaGraph(version) }), nil, r)
		if err != nil {
			t.Fatal(err)
		}
		desc, err := emitForRange[struct{}, struct{}]("describe", nil, nil, r)
		if err != nil {
			t.Fatal(err)
		}
		for name, data := range map[string][]byte{"request.json": append(request, '\n'), "plan.json": plan, "describe.json": desc} {
			path := filepath.Join(root, fmt.Sprintf("v%d", version), name)
			if os.Getenv("UPDATE_SCHEMA_GOLDEN") == "1" {
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, data, 0644); err != nil {
					t.Fatal(err)
				}
			}
			expected, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(expected, data) {
				t.Fatalf("fixture %s differs", path)
			}
		}
		data, err := emitForRange("plan", withoutOutputs(func(struct{}) Call { return baselineSchemaGraph() }), nil, r)
		if err != nil {
			t.Fatal(err)
		}
		var envelope map[string]json.RawMessage
		if err := json.Unmarshal(data, &envelope); err != nil {
			t.Fatal(err)
		}
		delete(envelope, "opspecSchema")
		body, _ := json.Marshal(envelope)
		if baseline != nil && !bytes.Equal(baseline, body) {
			t.Fatal("lowering changed the baseline body")
		}
		baseline = body
	}
	data, err := Plan(func(struct{}) Call { return baselineSchemaGraph() }, nil)
	if err != nil || !bytes.Contains(data, []byte(`"opspecSchema":4`)) {
		t.Fatalf("default must choose newest: %v", err)
	}
}

func TestNegotiationBuildsSelectedGraphOnce(t *testing.T) {
	for _, r := range []schemaRange{{1, 1}, {2, 3}, {1, 99}} {
		calls := 0
		_, err := emitForRange("plan", withoutOutputs(func(struct{}) Call {
			calls++
			_ = Container("unused").Cmd("outputParams stdout outputs")
			_, _ = Outputs[schemaTypedOutputs](Container("unused"))
			if false {
				return schemaGraph(4)
			}
			return Container("alpine:3").Cmd("echo", "stdout outputParams outputs")
		}), nil, r)
		if err != nil || calls != 1 {
			t.Fatalf("calls=%d err=%v", calls, err)
		}
	}
	calls := 0
	_, err := emitForRange("describe", func(struct{}) (Call, struct{}) { calls++; return nil, struct{}{} }, nil, schemaRange{1, 1})
	if err != nil || calls != 0 {
		t.Fatal("describe ran graph", err)
	}
	_, err = emitForRange("plan", func(struct{}) (Call, struct{}) { calls++; return nil, struct{}{} }, nil, schemaRange{5, 6})
	if err == nil || calls != 0 || !strings.Contains(err.Error(), "upgrade the Go op SDK") {
		t.Fatal("no shared version ran graph", err)
	}
}

func TestNegotiationMainHelper(t *testing.T) {
	if os.Getenv("OPCTL_NEGOTIATION_HELPER") != "1" {
		return
	}
	os.Args = []string{"hello", os.Getenv("OPCTL_NEGOTIATION_PHASE")}
	Main(func(struct{}) Call {
		var v int
		fmt.Sscan(os.Getenv("OPCTL_NEGOTIATION_FEATURE"), &v)
		return schemaGraph(v)
	})
	os.Exit(0)
}
func TestNegotiationMain(t *testing.T) {
	helloDir(t)
	for _, tc := range []struct {
		schema, phase string
		feature       int
		want          string
		target        int
	}{
		{`{"min":2,"max":3}`, "plan", 1, "", 3},
		{`{"min":1,"max":99}`, "plan", 1, "", 4},
		{`{"min":1,"max":1}`, "describe", 4, "", 1},
		{`{"min":1,"max":1}`, "plan", 2, "this op uses container stdout/stderr bindings (opspec >= 2); upgrade opctl to a version supporting opspec 2", 0},
		{`{"min":1,"max":2}`, "plan", 3, "this op uses container outputs file (opspec >= 3); upgrade opctl to a version supporting opspec 3", 0},
		{`{"min":1,"max":3}`, "plan", 4, "this op uses typed container outputs (opspec >= 4); upgrade opctl to a version supporting opspec 4", 0},
		{`{"min":5,"max":6}`, "plan", 1, "upgrade the Go op SDK", 0},
	} {
		t.Run(fmt.Sprintf("%s/%s/%d", tc.phase, tc.schema, tc.feature), func(t *testing.T) {
			data := fmt.Sprintf(`{"contractVersion":1,"phase":%q,"opspecSchema":%s}`, tc.phase, tc.schema)
			output, diagnostic, err := runNegotiationMain(t, data, tc.phase, tc.feature)
			if tc.want != "" {
				if err == nil || len(output) != 0 || !strings.Contains(diagnostic, tc.want) {
					t.Fatalf("err=%v diagnostic=%s stdout length=%d", err, diagnostic, len(output))
				}
				return
			}
			var response response
			if err != nil || json.Unmarshal(output, &response) != nil || response.OpspecSchema != tc.target {
				t.Fatalf("target %d: %v %s", tc.target, err, diagnostic)
			}
		})
	}
	for _, schema := range []string{`null`, `{}`, `{"min":1}`, `{"max":4}`, `{"min":null,"max":4}`, `{"min":1,"max":null}`, `{"min":0,"max":4}`, `{"min":-1,"max":4}`, `{"min":4,"max":1}`, `{"min":1.5,"max":4}`, `{"min":1,"max":9999999999999999999999999}`, `"bad"`} {
		output, _, err := runNegotiationMain(t, `{"contractVersion":1,"phase":"plan","opspecSchema":`+schema+`}`, "plan", 1)
		if err == nil || len(output) != 0 {
			t.Fatalf("malformed range accepted: %s", schema)
		}
	}
	output, _, err := runNegotiationMain(t, `{"contractVersion":1,"phase":"plan"}`, "plan", 1)
	if err == nil || len(output) != 0 {
		t.Fatal("missing range accepted")
	}
}
func runNegotiationMain(t *testing.T, data, phase string, feature int) ([]byte, string, error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestNegotiationMainHelper$")
	cmd.Env = append(os.Environ(), "OPCTL_NEGOTIATION_HELPER=1", "OPCTL_NEGOTIATION_PHASE="+phase, fmt.Sprintf("OPCTL_NEGOTIATION_FEATURE=%d", feature), "OPCTL_BOOTSTRAP_REQUEST="+path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	return out, stderr.String(), err
}

func TestNullBootstrapVersionsDoNotReuseEarlierFields(t *testing.T) {
	for _, data := range []string{
		`{"contractVersion":1,"contractVersion":null,"phase":"plan","opspecSchema":{"min":1,"max":4}}`,
		`{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"min":null,"max":4}}`,
		`{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":4,"max":null}}`,
		`{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":4},"opspecSchema":null}`,
	} {
		output, _, err := runNegotiationMain(t, data, "plan", 1)
		if err == nil || len(output) != 0 {
			t.Fatal("null schema range reused an earlier value")
		}
	}
}
