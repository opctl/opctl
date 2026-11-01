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

type metadataInputs struct {
	Who String `opctl:"who,default=world"`
}

func metadataPlan(in metadataInputs) Call {
	return Container("alpine:3").Cmd("sh", "-ce", `echo "hello $WHO"`).EnvFrom("WHO", in.Who)
}

// opDir switches to a fresh op directory named name that holds files.
func opDir(t *testing.T, name string, files map[string]string) string {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), name)
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for file, content := range files {
		if err := os.WriteFile(filepath.Join(dir, file), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatal(err)
		}
	})
	return previous
}

var greeterFiles = map[string]string{
	"op.go":      "// Greets the caller.\n//\n// Runs echo in alpine:\n//\n//\tsh -ce 'echo hello'\npackage main\n",
	"gen.go":     "//go:build ignore\n\n// Generates nothing.\npackage main\n",
	"op_test.go": "// Tests the greeter.\npackage main\n",
}

func describedMetadata(t *testing.T, data []byte) (name string, description *string) {
	t.Helper()
	var response struct {
		Opspec map[string]json.RawMessage `json:"opspec"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(response.Opspec["name"], &name); err != nil {
		t.Fatal(err)
	}
	if raw, ok := response.Opspec["description"]; ok {
		description = new(string)
		if err := json.Unmarshal(raw, description); err != nil {
			t.Fatal(err)
		}
	}
	return name, description
}

func TestDescribeReadsPackageDocComment(t *testing.T) {
	previous := opDir(t, "greeter", greeterFiles)
	data, err := emitForRange[metadataInputs, struct{}]("describe", nil, nil, schemaRange{1, 1})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(previous, "testdata", "describe.metadata.json")
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
		t.Fatalf("describe = %s; want %s", data, expected)
	}
	_, description := describedMetadata(t, data)
	if description == nil || !strings.HasPrefix(*description, "Greets the caller.\n\nRuns echo in alpine:\n\n\tsh -ce") || strings.Contains(*description, "Generates") || strings.Contains(*description, "Tests") || strings.HasSuffix(*description, "\n") {
		t.Fatalf("description = %v", description)
	}
	plan, err := Plan(metadataPlan, map[string]any{"who": "world"})
	if err != nil {
		t.Fatal(err)
	}
	if _, description := describedMetadata(t, plan); description != nil {
		t.Fatalf("plan emitted description %q", *description)
	}
}

func TestDescribeWithoutMainPackageDocHasNoDescription(t *testing.T) {
	for name, files := range map[string]map[string]string{
		"library": {"lib.go": "// Library docs.\npackage lib\n"},
		"no-go":   {"README.md": "docs"},
		"no-doc":  {"op.go": "package main\n"},
		"ignored": {"op.go": "package main\n", "gen.go": "//go:build ignore\n\n// Generator.\npackage main\n"},
	} {
		t.Run(name, func(t *testing.T) {
			opDir(t, name, files)
			data, err := Describe[struct{}]()
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(data, []byte(`"description"`)) {
				t.Fatalf("describe = %s", data)
			}
		})
	}
}

func TestDescribeFailsWhenThePackageDocCannotBeRead(t *testing.T) {
	opDir(t, "mixed", map[string]string{"a.go": "package main\n", "b.go": "package other\n"})
	_, err := Describe[struct{}]()
	if err == nil || !strings.Contains(err.Error(), "read package doc comment") {
		t.Fatalf("err = %v", err)
	}
}

func TestOpOverridesMetadata(t *testing.T) {
	opDir(t, "greeter", greeterFiles)
	for _, tc := range []struct {
		op                Op
		name, description string
	}{
		{Op{}, "greeter", "Greets the caller."},
		{Op{Name: "hello"}, "hello", "Greets the caller."},
		{Op{Description: "Deploys the API."}, "greeter", "Deploys the API."},
		{Op{Name: "hello", Description: "Deploys the API."}, "hello", "Deploys the API."},
	} {
		for label, describe := range map[string]func(...Op) ([]byte, error){"Describe": Describe[metadataInputs], "DescribeWithOutputs": DescribeWithOutputs[metadataInputs, struct{}]} {
			data, err := describe(tc.op)
			if err != nil {
				t.Fatal(err)
			}
			name, description := describedMetadata(t, data)
			if name != tc.name || description == nil || !strings.HasPrefix(*description, tc.description) {
				t.Fatalf("%s(%+v) = %s", label, tc.op, data)
			}
		}
		plan, err := emitForRange("plan", withoutOutputs(metadataPlan), map[string]any{"who": "world"}, schemaRange{1, 1}, tc.op)
		if err != nil {
			t.Fatal(err)
		}
		if name, description := describedMetadata(t, plan); name != tc.name || description != nil {
			t.Fatalf("plan(%+v) = %s", tc.op, plan)
		}
	}
}

func TestOpAcceptsAtMostOneOverride(t *testing.T) {
	opDir(t, "greeter", greeterFiles)
	for label, emit := range map[string]func() ([]byte, error){
		"Describe": func() ([]byte, error) { return Describe[metadataInputs](Op{}, Op{}) },
		"DescribeWithOutputs": func() ([]byte, error) {
			return DescribeWithOutputs[metadataInputs, struct{}](Op{Name: "a"}, Op{Name: "b"})
		},
		"plan": func() ([]byte, error) {
			return emitForRange("plan", withoutOutputs(metadataPlan), nil, schemaRange{1, 1}, Op{}, Op{})
		},
	} {
		if _, err := emit(); err == nil || err.Error() != "pass at most one opctl.Op" {
			t.Fatalf("%s err = %v", label, err)
		}
	}
}

func TestMainAppliesOpOverridesHelper(t *testing.T) {
	if os.Getenv("OPCTL_SDK_TEST_HELPER") != "1" {
		return
	}
	os.Args = []string{"greeter", os.Getenv("OPCTL_SDK_TEST_PHASE")}
	Main(metadataPlan, Op{Name: "hello", Description: "Deploys the API."})
	os.Exit(0)
}

func TestMainAppliesOpOverrides(t *testing.T) {
	executable, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	opDir(t, "greeter", greeterFiles)
	for _, phase := range []string{"describe", "plan"} {
		request := filepath.Join(t.TempDir(), "request.json")
		if err := os.WriteFile(request, []byte(`{"contractVersion":1,"phase":"`+phase+`","opspecSchema":{"min":1,"max":1},"inputs":`+map[string]string{"describe": `{}`, "plan": `{"who":"world"}`}[phase]+`,"knownOutputs":{}}`), 0600); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(executable, "-test.run=^TestMainAppliesOpOverridesHelper$")
		cmd.Env = append(os.Environ(), "OPCTL_SDK_TEST_HELPER=1", "OPCTL_SDK_TEST_PHASE="+phase, "OPCTL_BOOTSTRAP_REQUEST="+request)
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("%s: %v: %s", phase, err, &stderr)
		}
		name, description := describedMetadata(t, stdout.Bytes())
		if name != "hello" || (phase == "describe") != (description != nil && *description == "Deploys the API.") || (phase == "plan" && description != nil) {
			t.Fatalf("%s = %s", phase, &stdout)
		}
	}
}
