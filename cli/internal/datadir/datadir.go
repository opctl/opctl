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

	// EventDBPath resolves the eventdb path
	EventDBPath() string
}

// ensureExists ensures resolvedDataDirPath exists
func ensureExists(
	resolvedDataDirPath string,
) error {
	if err := os.MkdirAll(resolvedDataDirPath, 0775|os.ModeSetgid); nil != err {
		return err
	}
	return nil
}

// New returns an initialized data dir instance
func New(
	dataDirPath string,
) (DataDir, error) {
	resolvedDataDirPath, err := filepath.Abs(dataDirPath)
	if nil != err {
		return nil, fmt.Errorf("error initializing opctl data dir; error was %v", err)
	}

	if err := ensureExists(resolvedDataDirPath); nil != err {
		return nil, fmt.Errorf("error initializing opctl data dir; error was %v", err)
	}

	// ensure we can write
	if err := ioutil.WriteFile(
		filepath.Join(resolvedDataDirPath, "write-test"),
		[]byte(""),
		0775,
	); nil != err {
		return nil, fmt.Errorf("error initializing opctl data dir; error was %v", err)
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

func (dd _datadir) EventDBPath() string {
	return filepath.Join(dd.resolvedPath, "dcg", "events")
}

func (dd _datadir) InitAndLock() error {
	if err := ensureExists(dd.resolvedPath); nil != err {
		return err
	}

	lockFilePath := filepath.Join(
		dd.resolvedPath,
		"pid.lock",
	)

	// claim resolvedDataDirPath as ours
	lockFile := lockfile.New()
	if err := lockFile.Lock(lockFilePath); nil != err {
		pIDOfExistingNode := lockFile.PIdOfOwner(lockFilePath)
		return fmt.Errorf("node already running w/ PId: %v", pIDOfExistingNode)
	}

	return nil
}
