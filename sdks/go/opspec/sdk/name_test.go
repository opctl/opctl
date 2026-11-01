package sdk

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestContainerNameGolden(t *testing.T) {
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	helloDir(t)
	data, err := planAtSchema(func(struct{}) Call {
		return Container("golang:1.25").Name("run-a-go-service-api").Cmd("go", "run", "main.go")
	}, nil, 1)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(previous, "testdata", "plan.container-name.json")
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
}

func TestContainerNameIsAnEscapedLiteral(t *testing.T) {
	data, err := Plan(func(struct{}) Call { return Container("a").Name("api-$(token)") }, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"name":"api-\\$(token)"`)) {
		t.Fatalf("plan = %s", data)
	}
}

func TestContainerNameCopiesAreIndependent(t *testing.T) {
	base := Container("a").Name("first")
	_ = base.Name("second")
	data, err := Plan(func(struct{}) Call { return base }, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"name":"first"`)) {
		t.Fatalf("derived call mutated base: %s", data)
	}
}
