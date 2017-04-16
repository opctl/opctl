package vos

import (
	"github.com/virtual-go/fs"
	"os"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ VOS

// virtual operating system interface
type VOS interface {
	fs.FS
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

	// Setenv sets the value of the environment variable named by the key.
	// It returns an error, if any.
	Setenv(key, value string) error
}

func New(fs fs.FS) VOS {
	return _VOS{
		FS: fs,
	}
}

type _VOS struct {
	fs.FS
}

func (this _VOS) Exit(code int) {
	os.Exit(code)
}

func (this _VOS) FindProcess(pid int) (*os.Process, error) {
	return os.FindProcess(pid)
}

func (this _VOS) Getenv(key string) string {
	return os.Getenv(key)
}

func (this _VOS) Getpid() int {
	return os.Getpid()
}

func (this _VOS) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (this _VOS) Getwd() (string, error) {
	return os.Getwd()
}
