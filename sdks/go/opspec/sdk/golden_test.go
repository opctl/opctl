package sdk

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

type goldenInputs struct {
	Env    String  `opctl:"env,description='Target environment, dev or prod',default=dev,enum=dev|prod"`
	Count  Number  `opctl:"count,default=2,minimum=1,maximum=3,format=integer"`
	Token  String  `opctl:"token,secret,minLength=1"`
	Src    Dir     `opctl:"src,default=$(./testdata/src)"`
	Config File    `opctl:"config,default=$(./testdata/config.txt)"`
	Prod   Boolean `opctl:"wantProd"`
}

func goldenCall(in goldenInputs) Call {
	return Container("alpine:3").Cmd("echo", in.Env.Value()).EnvFrom("COUNT", in.Count).EnvFrom("TOKEN", in.Token).EnvFrom("WANT_PROD", in.Prod).Dir("/src", in.Src).File("/config", in.Config)
}
func TestInputsGolden(t *testing.T) {
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "inputs")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(previous)
	describe, err := emitForRange[goldenInputs, struct{}]("describe", nil, nil, schemaRange{1, 1})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := planAtSchema(goldenCall, map[string]any{"env": "prod", "count": 3, "wantProd": true}, 1)
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string][]byte{"describe.inputs.json": describe, "plan.inputs.json": plan} {
		path := filepath.Join(previous, "testdata", name)
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
			t.Fatalf("%s = %s; want %s", name, data, expected)
		}
	}
}
