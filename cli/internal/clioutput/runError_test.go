package clioutput

import "testing"

func TestRunError(t *testing.T) {
	err := RunError{Message: "testing"}
	if err.Error() != "testing" {
		t.Error("run error Error() method is broken")
	}
}
