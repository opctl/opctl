package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
)

func TestCacheLifecycleAssertions(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		builds, describes, plans int
		warm, wantError          bool
	}{
		{"cold", 1, 1, 1, false, false}, {"warm", 0, 0, 1, true, false}, {"recompile", 1, 2, 1, true, true}, {"missing plan", 0, 2, 0, true, true},
	} {
		events := []model.Event{}
		cached := 1
		if tc.warm {
			cached = 2
		}
		cachedName := "bootstrap describe (cached)"
		for i := 0; i < cached; i++ {
			events = append(events, model.Event{CallStarted: &model.CallStarted{Call: model.Call{Name: &cachedName}}})
		}
		for name, count := range map[string]int{"bootstrap build": tc.builds, "bootstrap describe": tc.describes, "bootstrap plan": tc.plans} {
			for i := 0; i < count; i++ {
				name := name
				events = append(events, model.Event{CallStarted: &model.CallStarted{Call: model.Call{Name: &name}}})
			}
		}
		if err := checkCacheEvents(events, tc.warm); (err != nil) != tc.wantError {
			t.Fatalf("%s: %v", tc.name, err)
		}
	}
}

func TestCacheEntryRequiresMatchingMetadata(t *testing.T) {
	data := t.TempDir()
	root := filepath.Join(data, "go", "programs", "v1", "program-key")
	if err := os.MkdirAll(root, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "program"), []byte("executable"), 0700); err != nil {
		t.Fatal(err)
	}
	launcher := []byte("launcher executable")
	if err := os.WriteFile(filepath.Join(root, "launcher"), launcher, 0700); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"different-key", "program-key"} {
		metadata := []byte(fmt.Sprintf(`{"Key":%q,"LauncherChecksum":"%x","LauncherSize":%d}`, key, sha256.Sum256(launcher), len(launcher)))
		if err := os.WriteFile(filepath.Join(root, "metadata.json"), metadata, 0600); err != nil {
			t.Fatal(err)
		}
		_, _, err := entry(data)
		if (err == nil) != (key == "program-key") {
			t.Fatalf("key %s: %v", key, err)
		}
	}
}

func TestCacheEntryValidatesLauncher(t *testing.T) {
	for _, kind := range []string{"valid", "missing", "corrupt", "not-executable", "unexpected-file"} {
		t.Run(kind, func(t *testing.T) {
			data := t.TempDir()
			root := filepath.Join(data, "go", "programs", "v1", "key")
			if err := os.MkdirAll(root, 0700); err != nil {
				t.Fatal(err)
			}
			launcher := []byte("launcher executable")
			files := map[string][]byte{"program": []byte("op executable"), "launcher": launcher, "metadata.json": []byte(fmt.Sprintf(`{"Key":"key","LauncherChecksum":"%x","LauncherSize":%d}`, sha256.Sum256(launcher), len(launcher)))}
			for name, content := range files {
				if err := os.WriteFile(filepath.Join(root, name), content, 0700); err != nil {
					t.Fatal(err)
				}
			}
			switch kind {
			case "missing":
				os.Remove(filepath.Join(root, "launcher"))
			case "corrupt":
				os.WriteFile(filepath.Join(root, "launcher"), []byte("corrupt launcher"), 0700)
			case "not-executable":
				os.Chmod(filepath.Join(root, "launcher"), 0600)
			case "unexpected-file":
				os.WriteFile(filepath.Join(root, "sidecar"), nil, 0600)
			}
			_, _, err := entry(data)
			if (err == nil) != (kind == "valid") {
				t.Fatalf("%s: %v", kind, err)
			}
		})
	}
}
