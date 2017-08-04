package model

import (
	"context"
	"io"
	"os"
)

type ReadSeekCloser interface {
	io.ReadCloser
	io.Seeker
}

// PkgHandle is a provider agnostic interface for interacting w/ pkg content
type PkgHandle interface {
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

	// Ref returns the pkgRef of the pkg
	Ref() string
}

type PkgManifest struct {
	Description string            `yaml:"description"`
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
