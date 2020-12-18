package ihttp

import (
  "io"
  "time"
  "net/http"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ IHTTP

type IHTTP interface {

  // ServeContent replies to the request using the content in the
  // provided ReadSeeker. The main benefit of ServeContent over io.Copy
  // is that it handles Range requests properly, sets the MIME type, and
  // handles If-Match, If-Unmodified-Since, If-None-Match, If-Modified-Since,
  // and If-Range requests.
  //
  // If the response's Content-Type header is not set, ServeContent
  // first tries to deduce the type from name's file extension and,
  // if that fails, falls back to reading the first block of the content
  // and passing it to DetectContentType.
  // The name is otherwise unused; in particular it can be empty and is
  // never sent in the response.
  //
  // If modtime is not the zero time or Unix epoch, ServeContent
  // includes it in a Last-Modified header in the response. If the
  // request includes an If-Modified-Since header, ServeContent uses
  // modtime to decide whether the content needs to be sent at all.
  //
  // The content's Seek method must work: ServeContent uses
  // a seek to the end of the content to determine its size.
  //
  // If the caller has set w's ETag header formatted per RFC 7232, section 2.3,
  // ServeContent uses it to handle requests using If-Match, If-None-Match, or If-Range.
  //
  // Note that *os.File implements the io.ReadSeeker interface.
  ServeContent(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker)
}

func New() IHTTP {
  return _IHTTP{}
}

type _IHTTP struct{}

func (ihttp _IHTTP) ServeContent(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker) {
  http.ServeContent(w, req, name, modtime, content)
}
