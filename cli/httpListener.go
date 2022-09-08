package main

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	core "github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/handler"
	"github.com/opctl/opctl/webapp"
)

/*
*
listener is a generic interface for things which expose opctl via some protocol
*/
type listener interface {
	listen(
		ctx context.Context,
		address string,
	) error
}

/*
*
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

func (hd _httpListener) listen(
	ctx context.Context,
	address string,
) error {
	router := mux.NewRouter()
	router.UseEncodedPath()

	router.PathPrefix("/api/").Handler(
		hd.StripPrefixAndRecover(
			"/api/",
			handler.New(
				hd.opspecNodeCore,
			),
		),
	)

	buildFS, err := fs.Sub(webapp.Build, "build")
	if err != nil {
		return err
	}

	router.PathPrefix("/").Handler(http.FileServer(http.FS(buildFS)))

	httpServer := http.Server{
		Addr:    address,
		Handler: handlers.CORS()(router),
	}

	// catch signals to ensure shutdown properly happens
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		ctx, _ := context.WithTimeout(ctx, 5*time.Second)

		// little hammer
		httpServer.Shutdown(ctx)

		// big hammer
		httpServer.Close()
	}()

	fmt.Println(
		hd.cliColorer.Info(fmt.Sprintf("opctl listening on %s", address)),
	)

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// StripPrefix returns a handler that serves HTTP requests
// by removing the given prefix from the request URL's Path
// and invoking the handler h. StripPrefix handles a
// request for a path that doesn't begin with prefix by
// replying with an HTTP 404 not found error.
func (hd _httpListener) StripPrefixAndRecover(prefix string, h http.Handler) http.Handler {
	defer func() {
		// don't let panics from any operation kill the server.
		if panic := recover(); panic != nil {
			fmt.Printf("recovered from panic: %s\n", panic)
		}
	}()

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
