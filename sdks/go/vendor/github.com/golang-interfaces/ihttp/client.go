package ihttp

//go:generate counterfeiter -o fakeClient.go --fake-name FakeClient ./ Client

import (
	"net/http"
)

type Client interface {
	// Do is implemented by net/http/Client
	Do(
		req *http.Request,
	) (
		*http.Response,
		error,
	)
}
