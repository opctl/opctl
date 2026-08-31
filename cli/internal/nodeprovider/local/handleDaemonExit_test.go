//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/node/fakes"
)

const testAPIListenAddress = "127.0.0.1:42224"

// TestHandleDaemonExitReachable covers the CI flake this function exists for:
// two opctl invocations race to start a node, the loser dies with "address
// already in use", but the winner is serving the API so the command carries on.
func TestHandleDaemonExitReachable(t *testing.T) {
	for _, scenario := range []struct {
		name   string
		stderr string
	}{
		{
			name:   "lost the API port",
			stderr: "listen tcp 127.0.0.1:42224: bind: address already in use\n",
		},
		{
			name:   "lost the node lock",
			stderr: "node already running; to kill use \"sudo opctl node kill\"\n",
		},
		{
			name:   "lost the event db lock",
			stderr: "Cannot acquire directory lock on \"/root/opctl/dcg/events\"\n",
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			fakeNode := new(fakes.FakeNode)
			fakeNode.LivenessReturns(nil)

			actualNode, actualErr := handleDaemonExit(
				context.Background(),
				fakeNode,
				testAPIListenAddress,
				&exec.ExitError{},
				scenario.stderr,
			)

			if actualErr != nil {
				t.Fatalf("expected a reachable node API to be treated as success; got %v", actualErr)
			}
			if actualNode != fakeNode {
				t.Error("expected the api client node to be returned")
			}
		})
	}
}

func TestHandleDaemonExitUnreachable(t *testing.T) {
	fakeNode := new(fakes.FakeNode)
	fakeNode.LivenessReturns(errors.New("connection refused"))

	_, actualErr := handleDaemonExit(
		context.Background(),
		fakeNode,
		testAPIListenAddress,
		&exec.ExitError{},
		"  container runtime unusable\n",
	)

	if actualErr == nil {
		t.Fatal("expected an error when nothing serves the node API")
	}
	if actualErr.Error() != "container runtime unusable" {
		t.Errorf("expected the node's stderr to be surfaced; got %q", actualErr.Error())
	}
}

// TestHandleDaemonExitUnreachableWithoutStderr ensures we don't swallow the
// failure when the node exits silently.
func TestHandleDaemonExitUnreachableWithoutStderr(t *testing.T) {
	fakeNode := new(fakes.FakeNode)
	fakeNode.LivenessReturns(errors.New("connection refused"))

	_, actualErr := handleDaemonExit(
		context.Background(),
		fakeNode,
		testAPIListenAddress,
		errors.New("signal: killed"),
		"",
	)

	if actualErr == nil {
		t.Fatal("expected an error when nothing serves the node API")
	}
	if !strings.Contains(actualErr.Error(), testAPIListenAddress) ||
		!strings.Contains(actualErr.Error(), "signal: killed") {
		t.Errorf("expected the listen address & daemon error to be surfaced; got %q", actualErr.Error())
	}
}
