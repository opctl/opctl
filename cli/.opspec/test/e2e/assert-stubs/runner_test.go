package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunnerStubProof(t *testing.T) {
	runner, err := filepath.Abs("../cmd.sh")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name, text := range map[string]string{"dockerd": "#!/bin/sh\nexit 0\n", "opctl": "#!/bin/sh\nexit 0\n", "assert-stubs": "#!/bin/sh\nprintf '%s\\n' \"$*\" > \"$PROOF_MARKER\"\n"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0755); err != nil {
			t.Fatal(err)
		}
	}
	command := exec.Command("sh", runner)
	command.Env = append(os.Environ(), "PATH="+dir+":"+os.Getenv("PATH"), "stubProof=true", "opRef=/fixture/caller", "expect=success", "authAddGithub=false", "assertStubsPath="+filepath.Join(dir, "assert-stubs"), "PROOF_MARKER="+filepath.Join(dir, "called"))
	if data, err := command.CombinedOutput(); err != nil {
		t.Fatalf("%v: %s", err, data)
	}
	data, err := os.ReadFile(filepath.Join(dir, "called"))
	if err != nil || strings.TrimSpace(string(data)) != "/fixture/caller /args.yml" {
		t.Fatalf("%s %v", data, err)
	}
}
