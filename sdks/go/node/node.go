package node

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/core.go . OpNode

// OpNode is the main structure to run and interact with ops
type OpNode interface {
	// AddAuth records authentication within the core
	AddAuth(
		ctx context.Context,
		req model.AddAuthReq,
	) error

	GetEventStream(
		ctx context.Context,
		req *model.GetEventStreamReq,
	) (
		<-chan model.Event,
		error,
	)

	// KillOp kills a running op
	KillOp(
		ctx context.Context,
		req model.KillOpReq,
	) (
		err error,
	)

	// StartOp starts an op and returns the root call ID
	StartOp(
		ctx context.Context,
		req model.StartOpReq,
	) (
		rootCallID string,
		err error,
	)

	// GetData gets data
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	GetData(
		ctx context.Context,
		req model.GetDataReq,
	) (
		model.ReadSeekCloser,
		error,
	)

	// ListDescendants lists file system entries
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	ListDescendants(
		ctx context.Context,
		req model.ListDescendantsReq,
	) (
		[]*model.DirEntry,
		error,
	)
}
