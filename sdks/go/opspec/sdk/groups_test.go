package sdk

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func composeCall(struct{}) Call {
	return Serial(Parallel(Container("alpine:3").Cmd("echo", "a"), Container("alpine:3").Cmd("echo", "b")), Container("alpine:3").Cmd("echo", "done"))
}

func TestComposeGolden(t *testing.T) {
	expected, err := os.ReadFile("testdata/plan.compose.json")
	if err != nil {
		t.Fatal(err)
	}
	previous, _ := os.Getwd()
	dir := filepath.Join(t.TempDir(), "compose")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(previous)
	actual, err := planAtSchema(composeCall, nil, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Fatalf("compose = %s; want %s", actual, expected)
	}
}

func TestGroups(t *testing.T) {
	for _, tc := range []struct {
		name          string
		call          Call
		want, wantErr string
	}{
		{"drop empty children", Serial(Parallel(), Serial(Parallel()), Container("alpine")), `"serial":[{"container":{"image":{"ref":"alpine"}}}]`, ""},
		{"empty serial root", Serial(), "", "op graph has no calls"},
		{"empty parallel root", Parallel(Serial()), "", "op graph has no calls"},
		{"first child error", Serial(Container("alpine").Env("$first", ""), Container("alpine").Env("$second", "")), "", "$first"},
		{"nested error", Parallel(Serial(Container("alpine").EnvFrom("A", String{}))), "", "handle is not bound"},
		{"nil child", Serial(nil), "", "must return a call"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := Plan(func(struct{}) Call { return tc.call }, nil)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("error = %v; want %s", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Contains(data, []byte(tc.want)) || bytes.Contains(data, []byte(`"container":null`)) {
				t.Fatalf("graph = %s", data)
			}
		})
	}
}

func TestGroupsCopyTheirChildren(t *testing.T) {
	calls := []Call{Container("alpine").Cmd("original")}
	group := Parallel(calls...)
	calls[0] = Container("alpine").Cmd("changed")
	data, err := Plan(func(struct{}) Call { return group }, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"original"`)) || bytes.Contains(data, []byte(`"changed"`)) {
		t.Fatalf("caller mutated group: %s", data)
	}
}
