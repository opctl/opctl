package vos

import (
	"os"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Vos

// virtual operating system interface
type Vos interface {
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

func New() Vos {
	return _vos{}
}

type _vos struct{}

func (this _vos) Exit(code int) {
	os.Exit(code)
}

func (this _vos) FindProcess(pid int) (*os.Process, error) {
	return os.FindProcess(pid)
}

func (this _vos) Getenv(key string) string {
	return os.Getenv(key)
}

func (this _vos) Getpid() int {
	return os.Getpid()
}

func (this _vos) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (this _vos) Getwd() (string, error) {
	return os.Getwd()
}
