package pidfile

import (
	"context"
	"errors"
	"io"
	"os"
	"strconv"
	"syscall"
)

// TryGetLock attempts to claim dirPath's node lock for the current process.
//
// The lock is an advisory whole-file lock (flock) on <dirPath>/pid.lock, held
// for as long as the returned io.Closer stays open. Claiming it is atomic, so
// two nodes racing to start can't both win, and the kernel drops it when the
// holder dies, so it can't go stale. The file's contents are the holder's PID,
// which TryGetProcess uses to find the node to signal.
//
// A false return means a live node already holds the lock; the returned
// io.Closer is nil in that case. Otherwise the caller owns the lock and must
// keep the io.Closer open for the lifetime of the node.
func TryGetLock(
	ctx context.Context,
	dirPath string,
) (io.Closer, bool, error) {
	file, err := os.OpenFile(
		constructPIDFilePath(dirPath),
		os.O_CREATE|os.O_RDWR,
		0744,
	)
	if err != nil {
		return nil, false, err
	}

	switch flockErr := syscall.Flock(
		int(file.Fd()),
		syscall.LOCK_EX|syscall.LOCK_NB,
	); {
	case flockErr == nil:
		// the lock is ours until we close the file (or die)
	case errors.Is(flockErr, syscall.EWOULDBLOCK):
		file.Close()
		return nil, false, nil
	default:
		// the filesystem doesn't support locking (some network mounts); fall
		// back to asking whether the recorded PID is still alive. Less reliable
		// -- it can't tell us about a node which is starting up -- but it's what
		// we've always done.
		nodeProcess, err := TryGetProcess(ctx, dirPath)
		if err != nil || nodeProcess != nil {
			file.Close()
			return nil, false, err
		}
	}

	if err := writePID(file); err != nil {
		file.Close()
		return nil, false, err
	}

	return file, true, nil
}

// writePID records the current process' PID as the sole contents of file.
func writePID(
	file *os.File,
) error {
	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.WriteAt(
		[]byte(strconv.Itoa(os.Getpid())),
		0,
	); err != nil {
		return err
	}

	return file.Sync()
}
