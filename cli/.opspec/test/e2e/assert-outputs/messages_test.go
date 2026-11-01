package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRunnerExpectedMessage(t *testing.T) {
	runner, err := filepath.Abs("../cmd.sh")
	if err != nil {
		t.Fatal(err)
	}
	for _, message := range []string{"unknown field Gone", "different error"} {
		dir := t.TempDir()
		for name, text := range map[string]string{"dockerd": "#!/bin/sh\nexit 0\n", "opctl": "#!/bin/sh\nif [ \"$1\" = run ]; then printf '%s\\n' \"$FAKE_MESSAGE\"; exit 1; fi\nexit 0\n"} {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0755); err != nil {
				t.Fatal(err)
			}
		}
		cmd := exec.Command("sh", runner)
		cmd.Env = append(os.Environ(), "PATH="+dir+":"+os.Getenv("PATH"), "expect=failure", "expectMessage=unknown field Gone", "FAKE_MESSAGE="+message, "authAddGithub=false")
		output, err := cmd.CombinedOutput()
		if (err == nil) != (message == "unknown field Gone") {
			t.Fatalf("%s: %v\n%s", message, err, output)
		}
	}
}
