package api

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"runtime/debug"
	"strings"

	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	core "github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/handler"
	"github.com/opctl/opctl/webapp"
)

/*
*
Listen for API requests
*/
func Listen(
	ctx context.Context,
	address string,
	opctlNodeCore core.Core,
) error {
	router := mux.NewRouter()
	router.UseEncodedPath()

	router.PathPrefix("/api/").Handler(
		stripPrefixAndRecover(
			"/api/",
			handler.New(
				opctlNodeCore,
			),
		),
	)

	buildFS, err := fs.Sub(webapp.Build, "build")
	if err != nil {
		return err
	}

	router.PathPrefix("/").Handler(http.FileServer(http.FS(buildFS)))

	apiServer := http.Server{
		Addr:    address,
		Handler: handlers.CORS()(router),
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// little hammer
		apiServer.Shutdown(ctx)

		// big hammer
		apiServer.Close()
	}()

	if err := apiServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// StripPrefix returns a handler that serves HTTP requests
// by removing the given prefix from the request URL's Path
// and invoking the handler h. StripPrefix handles a
// request for a path that doesn't begin with prefix by
// replying with an HTTP 404 not found error.
func stripPrefixAndRecover(prefix string, h http.Handler) http.Handler {
	defer func() {
		// don't let panics from any operation kill the server.
		if panic := recover(); panic != nil {
			fmt.Printf("recovered from panic: %s\n%s", panic, string(debug.Stack()))
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
