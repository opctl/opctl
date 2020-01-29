package datadir

import (
	"fmt"
	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-utils/lockfile"
)

// DataDir is an interface exposing the functionality we require in conjunction with our "data dir".
type DataDir interface {
	// ClearNodeCreateErrorIfExists clears any recorded error which previously occured during node creation
	ClearNodeCreateErrorIfExists() error

	// InitAndLock initializes and locks an opctl data dir
	// - if dataDirPath is explicitly provided it will be used
	// - else if OPCTL_DATA_DIR env var is set it will be used
	// - else the OS specific "per user" app data path will be used.
	InitAndLock() error

	// RecordNodeCreateError records an error which occured during node creation
	RecordNodeCreateError(
		err error,
	) error

	// Path resolves the data dir path
	Path() string

	// EventDBPath resolves the eventdb path
	EventDBPath() string

	// TryGetNodeCreateError records an error which occured during node creation
	TryGetNodeCreateError() error
}

// resolvePath resolves the data dir path
func resolvePath(
	dataDirPath *string,
) (string, error) {
	if nil != dataDirPath {
		return filepath.Abs(*dataDirPath)
	} else if dataDirEnvVar, ok := os.LookupEnv("OPCTL_DATA_DIR"); ok {
		return filepath.Abs(dataDirEnvVar)
	}

	perUserAppDataPath, err := appdatapath.New().PerUser()
	return filepath.Join(
		perUserAppDataPath,
		"opctl",
	), err
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
	dataDirPath *string,
) (DataDir, error) {
	resolvedDataDirPath, err := resolvePath(dataDirPath)
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
	return filepath.Join(dd.resolvedPath, "dcg", "event.db")
}

func (dd _datadir) TryGetNodeCreateError() error {
	errBytes, err := ioutil.ReadFile(dd.nodeCreateErrorPath())
	if os.IsNotExist(err) {
		return nil
	} else if nil != err {
		return err
	}

	return fmt.Errorf(string(errBytes))
}

func (dd _datadir) RecordNodeCreateError(
	err error,
) error {
	return ioutil.WriteFile(dd.nodeCreateErrorPath(), []byte(err.Error()), 0775)
}

func (dd _datadir) ClearNodeCreateErrorIfExists() error {
	removeErr := os.Remove(dd.nodeCreateErrorPath())
	if !os.IsNotExist(removeErr) {
		return removeErr
	}

	return nil
}

func (dd _datadir) nodeCreateErrorPath() string {
	return filepath.Join(dd.resolvedPath, "node-create-error")
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
