package model

import (
	"io"
)

type ReadSeekCloser interface {
	io.ReadCloser
	io.Seeker
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
}

// PullCreds contains optional authentication attributes
type PullCreds struct {
	Username,
	Password string
}
