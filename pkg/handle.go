package pkg

//go:generate counterfeiter -o ./fakeHandle.go --fake-name FakeHandle ./ Handle

import (
	"github.com/opspec-io/sdk-golang/model"
)

// Handle provides a source agnostic interface for interacting w/ packages
type Handle interface {
	// ListContents lists contents of a package
	ListContents() (
		[]*model.PkgContent,
		error,
	)

	// GetContent gets content from a package
	GetContent(
		contentPath string,
	) (
		model.ReadSeekCloser,
		error,
	)

	// Ref returns the absolute pkg ref of the open pkg
	Ref() string
}
