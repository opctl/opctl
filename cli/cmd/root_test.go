package cmd

import (
	"io"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCachePruneWiring(t *testing.T) {
	for _, test := range []struct {
		name    string
		args    []string
		names   []string
		useFlag bool
		invalid bool
	}{
		{name: "env data dir", args: []string{"cache", "prune"}},
		{name: "flag overrides env", args: []string{"cache", "prune", "npm", "gomod"}, names: []string{"npm", "gomod"}, useFlag: true},
		{name: "invalid name", args: []string{"cache", "prune", "npm", "../x"}, invalid: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			originalConfig := nodeConfig
			t.Cleanup(func() { nodeConfig = originalConfig })
			dataDir := filepath.Join(t.TempDir(), "env-data")
			t.Setenv("OPCTL_DATA_DIR", dataDir)
			rootCmd, err := NewRootCmd()
			if err != nil {
				t.Fatal(err)
			}
			pruneCmd, _, err := rootCmd.Find([]string{"cache", "prune"})
			if err != nil || pruneCmd.CommandPath() != "opctl cache prune" {
				t.Fatalf("cache prune command not found: %v", err)
			}
			called := false
			// Exercise argument and flag handling without sudo or a running node.
			pruneCmd.RunE = func(cmd *cobra.Command, args []string) error {
				called = true
				if nodeConfig.DataDir != dataDir {
					t.Errorf("data dir = %q, want %q", nodeConfig.DataDir, dataDir)
				}
				if !slices.Equal(args, test.names) {
					t.Errorf("names = %v, want %v", args, test.names)
				}
				return nil
			}
			args := test.args
			if test.useFlag {
				dataDir = filepath.Join(t.TempDir(), "flag-data")
				args = append([]string{"--data-dir", dataDir}, args...)
			}
			rootCmd.SetArgs(args)
			rootCmd.SetOut(io.Discard)
			rootCmd.SetErr(io.Discard)
			err = rootCmd.Execute()
			if test.invalid {
				if err == nil || !strings.Contains(err.Error(), "invalid cache name") || called {
					t.Fatalf("expected rejection before RunE, called = %v, error = %v", called, err)
				}
			} else if err != nil || !called {
				t.Fatalf("called = %v, error = %v", called, err)
			}
		})
	}
}
