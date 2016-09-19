package opspec

import (
  "net/http"
)

type httpClient interface {

  // Do is implemented by net/http/client
  Do(
  req *http.Request,
  ) (
  *http.Response,
  error,
  )

}
