// Package core defines the core interface for an opspec node
package core

import "github.com/opspec-io/sdk-golang/model"
import "io"

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

type ReadSeekCloser interface {
	io.ReadCloser
	io.Seeker
}

type Core interface {
	GetEventStream(
		req *model.GetEventStreamReq,
		eventChannel chan *model.Event,
	) error

	KillOp(
		req model.KillOpReq,
	)

	StartOp(
		req model.StartOpReq,
	) (
		callId string,
		err error,
	)

	ListPkgContents(
		pkgRef string,
	) (
		[]*model.PkgContent,
		error,
	)

	GetPkgContent(
		pkgRef string,
		path string,
	) (ReadSeekCloser, error)
}
