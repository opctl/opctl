package containerruntime

import (
	"os"
	"reflect"
	"testing"
)

func TestProxyEnvVars(t *testing.T) {
	t.Run("proxy vars in the node env are added to the container env", func(t *testing.T) {
		clearProxyEnv(t)
		t.Setenv("HTTP_PROXY", "http://proxy:3128")
		t.Setenv("https_proxy", "http://proxy:3129")

		actual := ProxyEnvVars(map[string]string{})

		expected := map[string]string{
			"HTTP_PROXY":  "http://proxy:3128",
			"https_proxy": "http://proxy:3129",
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("an op-provided value for the same key is not overridden", func(t *testing.T) {
		clearProxyEnv(t)
		t.Setenv("HTTP_PROXY", "http://node-proxy:3128")

		actual := ProxyEnvVars(map[string]string{
			"HTTP_PROXY": "http://op-proxy:8080",
		})

		if len(actual) != 0 {
			t.Fatalf("expected op-provided HTTP_PROXY to be left untouched, got %v", actual)
		}
	})

	t.Run("when no proxy vars are set the result is empty", func(t *testing.T) {
		clearProxyEnv(t)

		actual := ProxyEnvVars(map[string]string{
			"SOME_OTHER_VAR": "value",
		})

		if len(actual) != 0 {
			t.Fatalf("expected no propagated vars, got %v", actual)
		}
	})
}

// clearProxyEnv unsets every proxy env var so the host's own environment cannot
// influence the result, restoring prior values when the test completes.
func clearProxyEnv(t *testing.T) {
	t.Helper()
	for _, name := range proxyEnvVarNames {
		prev, had := os.LookupEnv(name)
		if had {
			t.Cleanup(func() { os.Setenv(name, prev) })
		} else {
			t.Cleanup(func() { os.Unsetenv(name) })
		}
		os.Unsetenv(name)
	}
}
