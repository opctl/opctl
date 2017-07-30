package pkg

//go:generate counterfeiter -o ./fakeHandle.go --fake-name FakeHandle ./ Handle

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
)

// Handle is a provider agnostic interface for interacting w/ a package
type Handle interface {
	// ListContents lists contents of a package
	ListContents(
		ctx context.Context,
	) (
		[]*model.PkgContent,
		error,
	)

	// GetContent gets content from a package
	GetContent(
		ctx context.Context,
		contentPath string,
	) (
		model.ReadSeekCloser,
		error,
	)

	// Ref returns the pkgRef of the pkg
	Ref() string
}
