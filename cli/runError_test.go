package main

import "testing"

func TestRunError(t *testing.T) {
	err := RunError{message: "testing"}
	if err.Error() != "testing" {
		t.Error("run error Error() method is broken")
	}
}
