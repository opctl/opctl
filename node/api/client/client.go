// Package client implements a client for the opspec node api
package client

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Client

import (
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/golang-interfaces/vhttp"
	"github.com/gorilla/websocket"
	"github.com/opspec-io/sdk-golang/model"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	GetEventStream(
		req *model.GetEventStreamReq,
	) (
		stream chan model.Event,
		err error,
	)

	KillOp(
		req model.KillOpReq,
	) (
		err error,
	)

	StartOp(
		req model.StartOpReq,
	) (
		opId string,
		err error,
	)
}

func New(
	baseUrl url.URL,
) (Client, error) {
	return &client{
		baseUrl: baseUrl,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		wsDialer: websocket.DefaultDialer,
	}, nil
}

type client struct {
	baseUrl    url.URL
	httpClient vhttp.Client
	wsDialer   iwebsocket.Dialer
}
