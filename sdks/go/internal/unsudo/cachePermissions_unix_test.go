//go:build unix

package unsudo

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

func TestRestoreCacheTreeUnprivileged(t *testing.T) {
	if cachePath := os.Getenv("OPCTL_TEST_UNPRIVILEGED_CACHE"); cachePath != "" {
		foreignFile := filepath.Join(cachePath, "foreign", "file")
		if _, err := os.Stat(foreignFile); err != nil {
			if errors.Is(err, fs.ErrPermission) {
				t.Skipf("unprivileged helper cannot access fixture in this environment: %v", err)
			}
			t.Fatal(err)
		}
		if err := os.Chmod(foreignFile, 0666); !errors.Is(err, syscall.EPERM) {
			t.Fatalf("expected EPERM for another user's file, got %v", err)
		}
		if err := RestoreCacheTree(cachePath); err != nil {
			t.Fatalf("another user's entries should be skipped silently: %v", err)
		}
		for path, mode := range map[string]os.FileMode{
			foreignFile:                         0644,
			filepath.Join(cachePath, "foreign"): 0755,
			filepath.Join(cachePath, "owned"):   0666,
		} {
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != mode {
				t.Errorf("%s: got %v, want %v", path, info.Mode().Perm(), mode)
			}
		}
		return
	}
	if os.Geteuid() != 0 {
		t.Skip("requires root to create foreign-owned fixtures and drop privileges")
	}
	t.Setenv("TMPDIR", t.TempDir())
	// Avoid TMPDIR parents that only root can traverse, such as under sudo -E.
	fixturePath, err := os.MkdirTemp("/tmp", "opctl-cache-permissions-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fixturePath)
	if err := os.Chmod(fixturePath, 0755); err != nil {
		t.Fatal(err)
	}
	cachePath := filepath.Join(fixturePath, "cache")
	if err := os.MkdirAll(filepath.Join(cachePath, "foreign"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{cachePath, filepath.Join(cachePath, "foreign")} {
		if err := os.Chmod(path, 0755); err != nil {
			t.Fatal(err)
		}
	}
	foreignFile := filepath.Join(cachePath, "foreign", "file")
	if err := os.WriteFile(foreignFile, nil, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(foreignFile, 0644); err != nil {
		t.Fatal(err)
	}
	// a root-owned dir the unprivileged node can't even read must be skipped without a warning
	privatePath := filepath.Join(cachePath, "private")
	if err := os.Mkdir(privatePath, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(privatePath, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(privatePath, "file"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	ownedFile := filepath.Join(cachePath, "owned")
	if err := os.WriteFile(ownedFile, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chown(ownedFile, 65534, 65534); err != nil {
		t.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	helperPath := filepath.Join(fixturePath, "test-helper")
	if err := os.WriteFile(helperPath, contents, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(helperPath, 0755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(helperPath, "-test.run=^TestRestoreCacheTreeUnprivileged$", "-test.v")
	cmd.Env = append(os.Environ(), "OPCTL_TEST_UNPRIVILEGED_CACHE="+cachePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Credential: &syscall.Credential{Uid: 65534, Gid: 65534}}
	if output, err := cmd.CombinedOutput(); err != nil {
		if errors.Is(err, fs.ErrPermission) {
			t.Skipf("cannot execute unprivileged helper in this environment (check /tmp noexec and directory access): %v", err)
		}
		t.Fatalf("unprivileged helper: %v\n%s", err, output)
	} else {
		t.Logf("unprivileged helper:\n%s", output)
	}
}
