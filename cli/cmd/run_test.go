package cmd

import "testing"

func TestShowLiveGraph(t *testing.T) {
	if !showLiveGraph(false, false) {
		t.Error("an unfocused terminal should show live progress")
	}
	if showLiveGraph(false, true) {
		t.Error("a focused container should suppress live progress")
	}
	if showLiveGraph(true, false) {
		t.Error("--no-progress should suppress live progress even without focus")
	}
}
