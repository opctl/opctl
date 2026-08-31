package pidfile

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

func TestTryGetLock(t *testing.T) {
	dirPath := t.TempDir()

	lock, gotLock, err := TryGetLock(context.Background(), dirPath)
	if err != nil {
		t.Fatalf("expected to claim the lock of an unused dir; got %v", err)
	}
	if !gotLock {
		t.Fatal("expected to claim the lock of an unused dir")
	}
	defer lock.Close()

	pIDFileBytes, err := os.ReadFile(filepath.Join(dirPath, "pid.lock"))
	if err != nil {
		t.Fatalf("expected a pid file to have been written; got %v", err)
	}
	if string(pIDFileBytes) != strconv.Itoa(os.Getpid()) {
		t.Errorf("expected the pid file to hold our PID; got %q", string(pIDFileBytes))
	}
}

// TestTryGetLockHeld is the race this lock exists to prevent: two nodes
// starting at once must not both conclude they own the data dir, else they race
// to bind the API port and one dies with "address already in use".
//
// flock is held per open file description, so a second claim conflicts even
// from within this process.
func TestTryGetLockHeld(t *testing.T) {
	dirPath := t.TempDir()

	lock, gotLock, err := TryGetLock(context.Background(), dirPath)
	if err != nil || !gotLock {
		t.Fatalf("expected to claim the lock of an unused dir; got %v, %v", gotLock, err)
	}
	defer lock.Close()

	contendedLock, gotLock, err := TryGetLock(context.Background(), dirPath)
	if err != nil {
		t.Fatalf("expected a held lock to be reported, not to error; got %v", err)
	}
	if gotLock {
		contendedLock.Close()
		t.Fatal("expected not to claim a held lock")
	}
	if contendedLock != nil {
		t.Error("expected no lock handle when the lock is held")
	}
}

// TestTryGetLockReleased covers a node exiting: the lock must not outlive its
// holder, and the pid file left behind must not block the next node.
func TestTryGetLockReleased(t *testing.T) {
	dirPath := t.TempDir()

	lock, gotLock, err := TryGetLock(context.Background(), dirPath)
	if err != nil || !gotLock {
		t.Fatalf("expected to claim the lock of an unused dir; got %v, %v", gotLock, err)
	}
	lock.Close()

	reclaimedLock, gotLock, err := TryGetLock(context.Background(), dirPath)
	if err != nil {
		t.Fatalf("expected to reclaim a released lock; got %v", err)
	}
	if !gotLock {
		t.Fatal("expected to reclaim a released lock")
	}
	reclaimedLock.Close()
}

// TestTryGetLockStalePIDFile covers the fallback path taken on filesystems
// without lock support, where ownership is decided by the recorded PID.
func TestTryGetLockStalePIDFile(t *testing.T) {
	for _, scenario := range []struct {
		name     string
		contents string
	}{
		{name: "empty", contents: ""},
		{name: "garbled", contents: "not-a-pid"},
		{name: "dead process", contents: deadPID(t)},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			dirPath := t.TempDir()
			if err := os.WriteFile(
				filepath.Join(dirPath, "pid.lock"),
				[]byte(scenario.contents),
				0744,
			); err != nil {
				t.Fatal(err)
			}

			nodeProcess, err := TryGetProcess(context.Background(), dirPath)
			if err != nil {
				t.Fatalf("expected a stale pid file not to error; got %v", err)
			}
			if nodeProcess != nil {
				t.Error("expected no node process from a stale pid file")
			}

			lock, gotLock, err := TryGetLock(context.Background(), dirPath)
			if err != nil || !gotLock {
				t.Fatalf("expected to claim a stale lock; got %v, %v", gotLock, err)
			}
			lock.Close()
		})
	}
}

// deadPID returns the PID of a process which has exited & been reaped, i.e.
// what a node which died without cleaning up leaves behind.
func deadPID(
	t *testing.T,
) string {
	t.Helper()

	cmd := exec.Command("sh", "-c", "exit 0")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	return fmt.Sprintf("%d", cmd.Process.Pid)
}
