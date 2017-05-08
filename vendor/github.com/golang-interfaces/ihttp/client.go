package ihttp

//go:generate counterfeiter -o fake.go --fake-name Fake ./ Client

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
