//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"os"
	"testing"
)

// TestDaemonEnv verifies the daemonized node inherits proxy vars (so op
// containers can reach the network behind a forward proxy) but does not
// inherit arbitrary env. Save/restore env so the host can't perturb results.
func TestDaemonEnv(t *testing.T) {
	for _, name := range []string{"HTTP_PROXY", "https_proxy", "NO_PROXY", "OPCTL_UNRELATED"} {
		if orig, ok := os.LookupEnv(name); ok {
			defer os.Setenv(name, orig)
		} else {
			defer os.Unsetenv(name)
		}
	}

	os.Setenv("HTTP_PROXY", "http://proxy.example:3128")
	os.Setenv("https_proxy", "http://proxy.example:3128")
	os.Unsetenv("NO_PROXY")
	os.Setenv("OPCTL_UNRELATED", "should-not-propagate")

	env := daemonEnv()

	has := func(s string) bool {
		for _, e := range env {
			if e == s {
				return true
			}
		}
		return false
	}

	if !has("HTTP_PROXY=http://proxy.example:3128") {
		t.Errorf("expected HTTP_PROXY to be passed through; got %v", env)
	}
	if !has("https_proxy=http://proxy.example:3128") {
		t.Errorf("expected https_proxy to be passed through; got %v", env)
	}
	for _, e := range env {
		if e == "NO_PROXY=" {
			t.Errorf("did not expect an unset NO_PROXY to be added; got %v", env)
		}
		if len(e) >= 15 && e[:15] == "OPCTL_UNRELATED" {
			t.Errorf("did not expect arbitrary env to be inherited; got %v", env)
		}
	}
}
