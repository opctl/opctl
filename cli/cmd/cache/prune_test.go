package cache

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrune(t *testing.T) {
	for _, test := range []struct {
		name      string
		names     []string
		remaining []string
	}{
		{name: "all"},
		{name: "named", names: []string{"Go-mod_v1.2", "npm"}, remaining: []string{".hidden"}},
		{name: "missing", names: []string{"missing"}, remaining: []string{"Go-mod_v1.2", "npm", ".hidden"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			dataDir := t.TempDir()
			for _, name := range []string{"Go-mod_v1.2", "npm", ".hidden"} {
				writeContents(t, filepath.Join(dataDir, "caches", name))
			}
			writeContents(t, filepath.Join(dataDir, "nocache", "run", "npm"))
			writeContents(t, filepath.Join(dataDir, "ops"))

			if err := prune(dataDir, test.names); err != nil {
				t.Fatal(err)
			}
			entries, err := os.ReadDir(filepath.Join(dataDir, "caches"))
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != len(test.remaining) {
				t.Fatalf("got %d cache dirs, want %d", len(entries), len(test.remaining))
			}
			for _, name := range test.remaining {
				assertContents(t, filepath.Join(dataDir, "caches", name))
			}
			assertContents(t, filepath.Join(dataDir, "nocache", "run", "npm"))
			assertContents(t, filepath.Join(dataDir, "ops"))
		})
	}
}

func TestPruneRejectsInvalidNamesBeforeDeleting(t *testing.T) {
	for _, name := range []string{"..", "a/b", "../x", ".", "", "/absolute", `a\b`, "a b", "a\n"} {
		t.Run(name, func(t *testing.T) {
			dataDir := t.TempDir()
			writeContents(t, filepath.Join(dataDir, "caches", "valid"))
			writeContents(t, filepath.Join(dataDir, "x"))
			writeContents(t, filepath.Join(dataDir, "caches", "a", "b"))

			err := prune(dataDir, []string{"valid", name})
			if err == nil || !strings.Contains(err.Error(), "invalid cache name") {
				t.Fatalf("expected invalid cache name error, got %v", err)
			}
			assertContents(t, filepath.Join(dataDir, "caches", "valid"))
			assertContents(t, filepath.Join(dataDir, "x"))
			assertContents(t, filepath.Join(dataDir, "caches", "a", "b"))
		})
	}
}

func TestPruneMissingRoot(t *testing.T) {
	for _, names := range [][]string{nil, {"missing"}} {
		if err := prune(t.TempDir(), names); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPruneRemovesSymlinksWithoutFollowingThem(t *testing.T) {
	for _, test := range []struct {
		name  string
		names []string
	}{
		{name: "all"},
		{name: "named", names: []string{"linked", "nested"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			dataDir := t.TempDir()
			target := t.TempDir()
			writeContents(t, target)
			cacheRoot := filepath.Join(dataDir, "caches")
			if err := os.MkdirAll(filepath.Join(cacheRoot, "nested"), 0700); err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"linked", filepath.Join("nested", "link")} {
				if err := os.Symlink(target, filepath.Join(cacheRoot, name)); err != nil {
					t.Fatal(err)
				}
			}

			if err := prune(dataDir, test.names); err != nil {
				t.Fatal(err)
			}
			assertContents(t, target)
			entries, err := os.ReadDir(cacheRoot)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Fatalf("cache root still has entries: %v", entries)
			}
		})
	}
}

func TestPruneRejectsSymlinkCacheRoot(t *testing.T) {
	for _, names := range [][]string{nil, {"npm"}} {
		dataDir := t.TempDir()
		target := t.TempDir()
		writeContents(t, filepath.Join(target, "npm"))
		if err := os.Symlink(target, filepath.Join(dataDir, "caches")); err != nil {
			t.Fatal(err)
		}
		if err := prune(dataDir, names); err == nil {
			t.Fatal("expected an error for a symlink cache root")
		}
		assertContents(t, filepath.Join(target, "npm"))
	}
}

func writeContents(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "contents"), []byte("cached"), 0600); err != nil {
		t.Fatal(err)
	}
}

func assertContents(t *testing.T, dir string) {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(dir, "contents"))
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "cached" {
		t.Fatalf("contents in %q changed: %q", dir, contents)
	}
}
