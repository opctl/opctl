package pkg

//go:generate counterfeiter -o ./fakeHandle.go --fake-name fakeHandle ./ handle

import (
	"github.com/opspec-io/sdk-golang/model"
)

// handle represents an open pkg
type handle interface {
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
}
