package ios

import (
	"os"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ IOS

// virtual operating system interface
type IOS interface {
	// Chmod changes the mode of the named file to mode.
	// If the file is a symbolic link, it changes the mode of the link's target.
	// If there is an error, it will be of type *PathError.
	Chmod(name string, mode os.FileMode) error

	// Create creates the named file with mode 0666 (before umask), truncating
	// it if it already exists. If successful, methods on the returned
	// File can be used for I/O; the associated file descriptor has mode
	// O_RDWR.
	// If there is an error, it will be of type *PathError.
	Create(name string) (*os.File, error)

	// Exit causes the current program to exit with the given status code.
	// Conventionally, code zero indicates success, non-zero an error.
	// The program terminates immediately; deferred functions are not run.
	Exit(code int)

	// FindProcess looks for a running process by its pid.
	//
	// The Process it returns can be used to obtain information
	// about the underlying operating system process.
	//
	// On Unix systems, FindProcess always succeeds and returns a Process
	// for the given pid, regardless of whether the process exists.
	FindProcess(pid int) (*os.Process, error)

	// Getenv retrieves the value of the environment variable named by the key.
	// It returns the value, which will be empty if the variable is not present.
	Getenv(key string) string

	// Getpid returns the process id of the caller.
	Getpid() int

	// Getwd returns a rooted path name corresponding to the
	// current directory. If the current directory can be
	// reached via multiple paths (due to symbolic links),
	// Getwd may return any one of them.
	Getwd() (string, error)

	// MkdirAll creates a directory named path,
	// along with any necessary parents, and returns nil,
	// or else returns an error.
	// The permission bits perm are used for all
	// directories that MkdirAll creates.
	// If path is already a directory, MkdirAll does nothing
	// and returns nil.
	MkdirAll(path string, perm os.FileMode) error

	// Open opens the named file for reading. If successful, methods on
	// the returned file can be used for reading; the associated file
	// descriptor has mode O_RDONLY.
	// If there is an error, it will be of type *PathError.
	Open(name string) (*os.File, error)

	// OpenFile is the generalized open call; most users will use Open
	// or Create instead. It opens the named file with specified flag
	// (O_RDONLY etc.) and perm, (0666 etc.) if applicable. If successful,
	// methods on the returned File can be used for I/O.
	// If there is an error, it will be of type *PathError.
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)

	// RemoveAll removes path and any children it contains.
	// It removes everything it can but returns the first error
	// it encounters. If the path does not exist, RemoveAll
	// returns nil (no error).
	RemoveAll(path string) error

	// Setenv sets the value of the environment variable named by the key.
	// It returns an error, if any.
	Setenv(key, value string) error

	// Stat returns a FileInfo describing the named file.
	// If there is an error, it will be of type *PathError.
	Stat(name string) (os.FileInfo, error)
}

func New() IOS {
	return _IOS{}
}

type _IOS struct{}

func (ios _IOS) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (ios _IOS) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (ios _IOS) Exit(code int) {
	os.Exit(code)
}

func (ios _IOS) FindProcess(pid int) (*os.Process, error) {
	return os.FindProcess(pid)
}

func (ios _IOS) Getenv(key string) string {
	return os.Getenv(key)
}

func (ios _IOS) Getpid() int {
	return os.Getpid()
}

func (ios _IOS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (ios _IOS) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (ios _IOS) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (ios _IOS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (ios _IOS) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (ios _IOS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (ios _IOS) Getwd() (string, error) {
	return os.Getwd()
}
