package attach

import (
	"errors"
	"testing"
)

func TestCanForwardInput(t *testing.T) {
	for _, test := range []struct {
		name                          string
		stdinTTY, stdoutTTY           bool
		foregroundGroup, processGroup int
		err                           error
		want                          bool
	}{
		{name: "foreground", stdinTTY: true, stdoutTTY: true, foregroundGroup: 2, processGroup: 2, want: true},
		{name: "background", stdinTTY: true, stdoutTTY: true, foregroundGroup: 2, processGroup: 3},
		{name: "redirected stdin", stdoutTTY: true, foregroundGroup: 2, processGroup: 2},
		{name: "redirected stdout", stdinTTY: true, foregroundGroup: 2, processGroup: 2},
		{name: "foreground unavailable", stdinTTY: true, stdoutTTY: true, foregroundGroup: 2, processGroup: 2, err: errors.New("ioctl failed")},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := canForwardInput(test.stdinTTY, test.stdoutTTY, test.foregroundGroup, test.processGroup, test.err); got != test.want {
				t.Fatalf("canForwardInput = %v, want %v", got, test.want)
			}
		})
	}
}
