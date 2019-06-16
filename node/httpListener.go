package node

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/opctl/opctl/node/statik"
	"github.com/opctl/opctl/sdk/go/node/api/handler"
	"github.com/opctl/opctl/sdk/go/node/core"
	"github.com/opctl/opctl/util/clicolorer"
	"github.com/rakyll/statik/fs"
	"net/http"
	"strings"
)

/**
newHttpListener returns a listener which exposes opctl over http (& websockets via connection upgrades)
*/
func newHttpListener(
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
