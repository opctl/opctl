package datadir

import (
	"fmt"
	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/golang-utils/lockfile"
)

// DataDir is an interface exposing the functionality we require in conjunction with our "data dir".
type DataDir interface {
	// InitAndLock initializes and locks an opctl data dir
	InitAndLock() error

	// Path resolves the data dir path
	Path() string
}

// ensureExists ensures resolvedDataDirPath exists
func ensureExists(
	resolvedDataDirPath string,
) error {
	if err := os.MkdirAll(resolvedDataDirPath, 0775|os.ModeSetgid); err != nil {
		return err
	}
	return nil
}

// New returns an initialized data dir instance
func New(
	dataDirPath string,
) (DataDir, error) {
	resolvedDataDirPath, err := filepath.Abs(dataDirPath)
	if err != nil {
		return nil, fmt.Errorf("error initializing opctl data dir: %w", err)
	}

	if err := ensureExists(resolvedDataDirPath); err != nil {
		return nil, fmt.Errorf("error initializing opctl data dir: %w", err)
	}

	// ensure we can write
	if err := ioutil.WriteFile(
		filepath.Join(resolvedDataDirPath, "write-test"),
		[]byte(""),
		0775,
	); err != nil {
		return nil, fmt.Errorf("error initializing opctl data dir: %w", err)
	}

	return _datadir{
		resolvedPath: resolvedDataDirPath,
	}, err
}

type _datadir struct {
	resolvedPath string
}

func (dd _datadir) Path() string {
	return filepath.Join(dd.resolvedPath)
}

func (dd _datadir) InitAndLock() error {
	if err := ensureExists(dd.resolvedPath); err != nil {
		return err
	}

	lockFilePath := filepath.Join(
		dd.resolvedPath,
		"pid.lock",
	)

	// claim resolvedDataDirPath as ours
	lockFile := lockfile.New()
	if err := lockFile.Lock(lockFilePath); err != nil {
		pIDOfExistingNode := lockFile.PIdOfOwner(lockFilePath)
		return fmt.Errorf("node already running with PID: %d", pIDOfExistingNode)
	}

	return nil
}
