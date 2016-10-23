package http

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
