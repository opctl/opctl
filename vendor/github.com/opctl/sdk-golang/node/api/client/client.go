// Package client implements a client for the opspec node api
package client

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Client

import (
	"context"
	"net/url"

	"github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/golang-interfaces/ihttp"
	"github.com/gorilla/websocket"
	"github.com/opctl/sdk-golang/model"
	"github.com/sethgrid/pester"
)

type Client interface {
	GetEventStream(
		req *model.GetEventStreamReq,
	) (
		stream chan model.Event,
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

	KillOp(
		ctx context.Context,
		req model.KillOpReq,
	) (
		err error,
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

	StartOp(
		ctx context.Context,
		req model.StartOpReq,
	) (
		opID string,
		err error,
	)
}

type Opts struct {
	// RetryLogHook will be executed anytime a request is retried
	RetryLogHook func(err error)
}

// New returns a new client
// nil opts will be ignored
func New(
	baseUrl url.URL,
	opts *Opts,
) Client {

	httpClient := pester.New()
	httpClient.Backoff = pester.ExponentialBackoff

	if nil != opts {
		// handle options
		httpClient.LogHook = func(errEntry pester.ErrEntry) {
			// wire up retry log hook
			opts.RetryLogHook(errEntry.Err)
		}
	}

	return &client{
		baseUrl:    baseUrl,
		httpClient: httpClient,
		wsDialer:   websocket.DefaultDialer,
	}
}

type client struct {
	baseUrl    url.URL
	httpClient ihttp.Client
	wsDialer   iwebsocket.Dialer
}
