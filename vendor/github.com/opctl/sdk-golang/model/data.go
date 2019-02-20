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
	// ListDescendants lists descendant of the data node pointed to by the current handle
	ListDescendants(
		ctx context.Context,
	) (
		[]*DirEntry,
		error,
	)

	// GetContent gets data from the current handle
	GetContent(
		ctx context.Context,
		contentPath string,
	) (
		ReadSeekCloser,
		error,
	)

	// Path returns the local path of the data; may or may not be same as Ref
	// returns nil if data doesn't exist locally
	Path() *string

	// Ref returns a ref to the data; may or may not be same as Path
	Ref() string
}

// DirEntry represents an entry in a directory (a sub directory or file)
type DirEntry struct {
	Path string
	Size int64
	Mode os.FileMode
}

// Value represents a typed value
type Value struct {
	Array   []interface{}          `json:"array,omitempty"`
	Boolean *bool                  `json:"boolean,omitempty"`
	Dir     *string                `json:"dir,omitempty"`
	File    *string                `json:"file,omitempty"`
	Number  *float64               `json:"number,omitempty"`
	Object  map[string]interface{} `json:"object,omitempty"`
	Socket  *string                `json:"socket,omitempty"`
	String  *string                `json:"string,omitempty"`
}

// PullCreds contains authentication attributes for auth'ing w/ a data provider
type PullCreds struct {
	Username,
	Password string
}
