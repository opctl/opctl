// Package client implements a client for the opctl node api
package client

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"net/url"

	iwebsocket "github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/golang-interfaces/ihttp"
	"github.com/gorilla/websocket"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/sethgrid/pester"
)

// Opts is options for an api client node
type Opts struct {
	// RetryLogHook will be executed anytime a request is retried
	RetryLogHook func(err error)
}

// New returns a new api client node
// nil opts will be ignored
func New(
	baseURL url.URL,
	opts *Opts,
) node.Node {

	httpClient := pester.New()
	// 90 second timeout
	httpClient.MaxRetries = 90

	if opts != nil {
		// handle options
		httpClient.LogHook = func(errEntry pester.ErrEntry) {
			// wire up retry log hook
			opts.RetryLogHook(errEntry.Err)
		}
	}

	return &apiClient{
		baseURL:    baseURL,
		httpClient: httpClient,
		wsDialer:   websocket.DefaultDialer,
	}
}

type apiClient struct {
	baseURL    url.URL
	httpClient ihttp.Client
	wsDialer   iwebsocket.Dialer
}
