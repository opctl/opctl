// Package client implements a client for the opspec node api
package client

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"net/url"

	iwebsocket "github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/golang-interfaces/ihttp"
	"github.com/gorilla/websocket"
	"github.com/sethgrid/pester"
)

// Opts is options for APIClient
type Opts struct {
	// RetryLogHook will be executed anytime a request is retried
	RetryLogHook func(err error)
}

// New returns a new APIClient
// nil opts will be ignored
func New(
	baseURL url.URL,
	opts *Opts,
) *APIClient {

	httpClient := pester.New()
	httpClient.Backoff = pester.ExponentialBackoff

	if nil != opts {
		// handle options
		httpClient.LogHook = func(errEntry pester.ErrEntry) {
			// wire up retry log hook
			opts.RetryLogHook(errEntry.Err)
		}
	}

	return &APIClient{
		baseURL:    baseURL,
		httpClient: httpClient,
		wsDialer:   websocket.DefaultDialer,
	}
}

// APIClient is an OpNode that runs ops with a remote OpNode over a network
// connection
type APIClient struct {
	baseURL    url.URL
	httpClient ihttp.Client
	wsDialer   iwebsocket.Dialer
}
