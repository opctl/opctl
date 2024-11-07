package datadir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rogpeppe/go-internal/lockedfile"
	"golang.org/x/sys/unix"
)

// DataDir is an interface exposing the functionality we require in conjunction with our "data dir".
type DataDir interface {
	// InitAndLock initializes and locks an opctl data dir
	InitAndLock() (unlock func(), err error)

	// Path resolves the data dir path
	Path() string
}

// ensureExists ensures resolvedDataDirPath exists
func ensureExists(
	resolvedDataDirPath string,
) error {
	// don't de-privilege group
	unix.Umask(0002)

	return os.MkdirAll(resolvedDataDirPath, 0770|os.ModeSetgid)
}

// New returns an initialized data dir instance
func New(
	dataDirPath string,
) (DataDir, error) {
	resolvedPath, err := filepath.Abs(dataDirPath)
	if err != nil {
		return nil, fmt.Errorf("error initializing opctl data dir: %w", err)
	}

	if err := ensureExists(resolvedPath); err != nil {
		return nil, fmt.Errorf("error initializing opctl data dir: %w", err)
	}

	// ensure we can write
	if err := os.WriteFile(
		filepath.Join(resolvedPath, "write-test"),
		[]byte(""),
		0775,
	); err != nil {
		return nil, fmt.Errorf("error initializing opctl data dir: %w", err)
	}

	return _datadir{
		resolvedPath: resolvedPath,
	}, err
}

type _datadir struct {
	resolvedPath string
}

func (dd _datadir) Path() string {
	return filepath.Join(dd.resolvedPath)
}

func (dd _datadir) InitAndLock() (unlock func(), err error) {
	if err := ensureExists(dd.resolvedPath); err != nil {
		return nil, err
	}

	lockedFile, err := lockedfile.Create(
		filepath.Join(
			dd.resolvedPath,
			"pid.lock",
		),
	)
	if err != nil {
		return func() {}, nil
	}

	_, err = lockedFile.Write(
		[]byte(
			fmt.Sprintf("%d", os.Getpid()),
		),
	)

	return func() { lockedFile.Close() }, err
}
