package fs

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ FS

import (
	"os"
)

// virtual filesystem interface
type FS interface {
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

	// Stat returns a FileInfo describing the named file.
	// If there is an error, it will be of type *PathError.
	Stat(name string) (os.FileInfo, error)
}
