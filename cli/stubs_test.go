package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoGeneratedTestSuiteStubsAreTracked(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "-c", "safe.directory="+root, "ls-files", "test-suite")
	cmd.Dir = ".."
	data, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "/opctlgen/") {
		t.Fatalf("generated op stubs are tracked:\n%s", data)
	}
}
