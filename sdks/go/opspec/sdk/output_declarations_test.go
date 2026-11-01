package sdk

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type outputInputs struct {
	Text String `opctl:"text"`
}
type outputValues struct {
	Echoed   String `opctl:"echoed"`
	Fallback String `opctl:"fallback,default=default"`
}

func TestOutputContractFixtures(t *testing.T) {
	fixtures := map[string][]byte{}
	for _, name := range []string{"describe.outputs.json", "plan.outputs.json"} {
		data, err := os.ReadFile("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		fixtures[name] = data
	}
	helloDir(t)
	desc, err := emitForRange[outputInputs, outputValues]("describe", nil, nil, schemaRange{1, 1})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := emitForRange("plan", func(in outputInputs) (Call, outputValues) { return Container("alpine"), outputValues{Echoed: in.Text} }, map[string]any{"text": "public"}, schemaRange{1, 1})
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string][]byte{"describe.outputs.json": desc, "plan.outputs.json": plan} {
		want := fixtures[name]
		if !bytes.Equal(data, want) {
			t.Fatalf("%s=%s want %s", name, data, want)
		}
	}
}
func TestOutputDeclarations(t *testing.T) {
	type Out struct {
		Dir     Dir         `opctl:"dir,description='Built files'"`
		File    File        `opctl:"file"`
		Socket  Socket      `opctl:"socket"`
		String  String      `opctl:"string,secret,minLength=1"`
		Number  Number      `opctl:"number,minimum=0"`
		Boolean Boolean     `opctl:"boolean,default=false"`
		Array   Array[uint] `opctl:"array,minItems=1"`
		Object  Object[struct {
			Label string `json:"label"`
		}] `opctl:"object,required=label"`
	}
	data, err := DescribeWithOutputs[struct{}, Out]()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"outputs"`, `"isSecret":true`, `"minimum":0`, `"label"`, `"default":false`} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("missing %s: %s", want, data)
		}
	}
	_, err = PlanWithOutputs(func(struct{}) (Call, outputValues) { return Container("alpine"), outputValues{} }, nil)
	if err == nil || !strings.Contains(err.Error(), "echoed") {
		t.Fatalf("%v", err)
	}
	type D struct {
		Out Dir `opctl:"out,default=./"`
	}
	_, err = PlanWithOutputs(func(struct{}) (Call, D) { return Container("alpine"), D{Out: NewDir()} }, nil)
	if err == nil {
		t.Fatal("unproduced export accepted with default")
	}
	_, err = PlanWithOutputs(func(struct{}) (Call, D) { d := NewDir(); return Container("alpine").Dir("/out", d), D{Out: d} }, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = PlanWithOutputs(func(struct{}) (Call, D) { c, d := OpOutput[Dir](CallOp("./yaml"), "dir"); return c, D{Out: d} }, nil)
	if err != nil {
		t.Fatal(err)
	}
}
func TestInvalidOutputDeclarations(t *testing.T) {
	for _, describe := range []func(...Op) ([]byte, error){
		DescribeWithOutputs[struct{}, int],
		DescribeWithOutputs[struct{}, struct {
			X int `opctl:"x"`
		}],
		DescribeWithOutputs[struct{}, struct {
			x String `opctl:"x"`
		}],
		DescribeWithOutputs[struct{}, struct {
			X *String `opctl:"x"`
		}],
		DescribeWithOutputs[struct{}, struct {
			X String `opctl:"opctl_reserved"`
		}],
		DescribeWithOutputs[struct{}, struct {
			X String `opctl:"x"`
			Y String `opctl:"x"`
		}],
		DescribeWithOutputs[struct{}, struct {
			X Boolean `opctl:"x,secret"`
		}],
	} {
		if _, err := describe(); err == nil {
			t.Fatal("accepted invalid output declaration")
		}
	}
}

func TestMainWithOutputsHelper(t *testing.T) {
	if os.Getenv("OPCTL_SDK_OUTPUT_HELPER") != "1" {
		return
	}
	os.Args = []string{"hello", os.Getenv("OPCTL_SDK_OUTPUT_PHASE")}
	MainWithOutputs(func(in outputInputs) (Call, outputValues) {
		if os.Args[1] == "describe" {
			panic("describe invoked plan")
		}
		return Container("alpine"), outputValues{Echoed: in.Text}
	})
	os.Exit(0)
}
func TestMainWithOutputsContract(t *testing.T) {
	for _, phase := range []string{"describe", "plan"} {
		t.Run(phase, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "request.json")
			inputs := `{}`
			if phase == "plan" {
				inputs = `{"text":"hello"}`
			}
			request := `{"contractVersion":1,"phase":"` + phase + `","opspecSchema":{"min":1,"max":1},"inputs":` + inputs + `,"knownOutputs":{}}`
			if err := os.WriteFile(path, []byte(request), 0600); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(os.Args[0], "-test.run=^TestMainWithOutputsHelper$")
			cmd.Env = append(os.Environ(), "OPCTL_SDK_OUTPUT_HELPER=1", "OPCTL_SDK_OUTPUT_PHASE="+phase, "OPCTL_BOOTSTRAP_REQUEST="+path)
			data, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("%s: %v %s", phase, err, data)
			}
			if phase == "describe" && !bytes.Contains(data, []byte(`"outputs"`)) {
				t.Fatalf("%s", data)
			}
			if phase == "plan" && !bytes.Contains(data, []byte(`"echoed":"$(text)"`)) {
				t.Fatalf("%s", data)
			}
		})
	}
}

func TestInvalidOutputHandleCannotUseDefault(t *testing.T) {
	type Out struct {
		Dir Dir `opctl:"dir,default=./"`
	}
	_, invalid := OpOutput[Dir](CallOp("./yaml"), "dir")
	if _, err := PlanWithOutputs(func(struct{}) (Call, Out) { return Container("alpine"), Out{Dir: invalid} }, nil); err == nil {
		t.Fatal("invalid op output silently used its default")
	}
}
