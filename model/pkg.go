package model

//go:generate counterfeiter -o ../pkg/fakeHandle.go --fake-name FakeHandle ./ DataHandle

import (
	"context"
	"io"
	"os"
)

type ReadSeekCloser interface {
	io.ReadCloser
	io.Seeker
}

// DataHandle is a provider agnostic interface for interacting w/ data
type DataHandle interface {
	// ListContents lists contents of a package
	ListContents(
		ctx context.Context,
	) (
		[]*PkgContent,
		error,
	)

	// GetContent gets content from a package
	GetContent(
		ctx context.Context,
		contentPath string,
	) (
		ReadSeekCloser,
		error,
	)

	// Path the local path of the pkg
	// returns nil if pkg doesn't exist locally
	Path() *string

	// Ref returns a ref to the data
	Ref() string
}

type PkgManifest struct {
	Description string            `yaml:"description"`
	Icon        string            `yaml:"icon,omitempty"`
	Inputs      map[string]*Param `yaml:"inputs,omitempty"`
	Name        string            `yaml:"name"`
	Outputs     map[string]*Param `yaml:"outputs,omitempty"`
	Run         *SCG              `yaml:"run,omitempty"`
	Version     string            `yaml:"version,omitempty"`
}

type PkgContent struct {
	Path string
	Size int64
	Mode os.FileMode
}

// PullCreds contains optional authentication attributes
type PullCreds struct {
	Username,
	Password string
}
