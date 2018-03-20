// Package client implements a client for the opspec node api
package client

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Client

import (
	"context"
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/golang-interfaces/ihttp"
	"github.com/gorilla/websocket"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/sethgrid/pester"
	"net/url"
)

type Client interface {
	GetEventStream(
		req *model.GetEventStreamReq,
	) (
		stream chan model.Event,
		err error,
	)

	// GetContent gets content from a package
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	GetPkgContent(
		ctx context.Context,
		req model.GetPkgContentReq,
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

	// ListContents lists contents of a package
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	ListPkgContents(
		ctx context.Context,
		req model.ListPkgContentsReq,
	) (
		[]*model.PkgContent,
		error,
	)

	StartOp(
		ctx context.Context,
		req model.StartOpReq,
	) (
		opId string,
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
