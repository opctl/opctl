package creater

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	_ "github.com/opctl/opctl/cli/internal/core/node/creater/statik"
	"github.com/opctl/opctl/sdks/go/node/api/handler"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/rakyll/statik/fs"
)

/**
listener is a generic interface for things which expose opctl via some protocol
*/
type listener interface {
	Listen(
		ctx context.Context,
	) <-chan error
}

/**
newHTTPListener returns a listener which exposes opctl over http (& websockets via connection upgrades)
*/
func newHTTPListener(
	opspecNodeCore core.Core,
) listener {
	return _httpListener{
		opspecNodeCore: opspecNodeCore,
		cliColorer:     clicolorer.New(),
	}
}

type _httpListener struct {
	opspecNodeCore core.Core
	cliColorer     clicolorer.CliColorer
}

func (hd _httpListener) Listen(
	ctx context.Context,
) <-chan error {
	router := mux.NewRouter()
	router.UseEncodedPath()

	router.PathPrefix("/api/").Handler(
		hd.StripPrefix(
			"/api/",
			handler.New(
				hd.opspecNodeCore,
			),
		),
	)

	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		statikFS, err := fs.New()
		if nil != err {
			errChan <- err
			return
		}
		router.PathPrefix("/").Handler(http.FileServer(statikFS))

		httpServer := http.Server{Addr: ":42224", Handler: handlers.CORS()(router)}

		ctx, cancel := context.WithCancel(ctx)
		go func() {
			err := httpServer.ListenAndServe()
			if nil != err {
				errChan <- err
				cancel()
			}
		}()

		<-ctx.Done()
		httpServer.Shutdown(ctx)
	}()

	fmt.Println(hd.cliColorer.Info("http listener bound to 0.0.0.0:42224"))
	return errChan
}

// StripPrefix returns a handler that serves HTTP requests
// by removing the given prefix from the request URL's Path
// and invoking the handler h. StripPrefix handles a
// request for a path that doesn't begin with prefix by
// replying with an HTTP 404 not found error.
func (hd _httpListener) StripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, prefix)
			h.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}
