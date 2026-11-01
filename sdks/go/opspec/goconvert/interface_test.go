package goconvert

import (
	"context"
	"encoding/json"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderInterface(t *testing.T) {
	source := []byte("name: interface-test\ndescription: releases opctl to [Github Container Registry](https://ghcr.io/opctl/opctl/)\ninputs:\n  MYSQL_USER:\n    description: |\n      line one\n      line `two`, 'three'\n    string: {default: 'a,b', constraints: {enum: ['a,b', 'a|b']}}\n  tags: {array: {default: [], constraints: {items: {type: string}}}}\n  settings: {object: {constraints: {additionalProperties: {type: number}}}}\n  n: {number: {default: 2, constraints: {minimum: 1, format: integer}}}\n  flag: {boolean: {default: false}}\noutputs: {result: {description: output, string: {}}}\nrun: {container: {image: {ref: alpine}, stdout: '$(result)'}}\n")
	a, err := Analyze(context.Background(), t.TempDir(), source, nil)
	if err != nil || len(a.Refusals) > 0 {
		t.Fatalf("%v %+v", err, a)
	}
	code, err := RenderInterface(a)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(code), "MysqlUser") || !strings.Contains(string(code), "opctl.Op{") {
		t.Fatal(string(code))
	}
	if strings.Index(string(code), "False ") > strings.Index(string(code), "Flag ") {
		t.Fatal("parameter order changed")
	}
	desc := describeCode(t, code)
	if desc.Name != a.Spec.Name || desc.Description != a.Spec.Description {
		t.Fatalf("metadata changed: %+v", desc)
	}
	for _, group := range []struct{ want, got map[string]*model.ParamSpec }{{a.Spec.Inputs, desc.Inputs}, {a.Spec.Outputs, desc.Outputs}} {
		x, _ := json.Marshal(group.want)
		y, _ := json.Marshal(group.got)
		if string(x) != string(y) {
			t.Fatalf("interface changed:\n%s\n%s", x, y)
		}
	}
}
func describeCode(t *testing.T, code []byte) *model.OpSpec {
	t.Helper()
	dir := t.TempDir()
	sdk, err := filepath.Abs("../sdk")
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string][]byte{"op.go": code, "go.mod": []byte("module ops.local/test\ngo 1.25\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => " + sdk + "\n")} {
		if err = os.WriteFile(filepath.Join(dir, name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	request, _ := json.Marshal(goop.NewRequest("describe", nil))
	requestPath := filepath.Join(dir, "request.json")
	if err = os.WriteFile(requestPath, request, 0600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "run", ".", "describe")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=readonly", "GOWORK=off", "GOTOOLCHAIN=local", "OPCTL_BOOTSTRAP_REQUEST="+requestPath)
	out, err := cmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			t.Fatalf("%v %s\n%s", err, e.Stderr, code)
		}
		t.Fatal(err)
	}
	var response struct{ Opspec *model.OpSpec }
	if err = json.Unmarshal(out, &response); err != nil {
		t.Fatal(err)
	}
	return response.Opspec
}
