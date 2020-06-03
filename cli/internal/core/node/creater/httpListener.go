package creater

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"os"
	"os/signal"
	"syscall"
	"time"

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
		address string,
	) error
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
	address string,
) error {
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

	statikFS, err := fs.New()
	if nil != err {
		return err
	}
	router.PathPrefix("/").Handler(http.FileServer(statikFS))

	httpServer := http.Server{
		Addr:    address,
		Handler: handlers.CORS()(router),
	}

	// catch signals to ensure shutdown properly happens
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		// little hammer
		httpServer.Shutdown(ctx)

		// big hammer
		httpServer.Close()
	}()

	fmt.Println(
		hd.cliColorer.Info(fmt.Sprintf("Binding opctl API to %s", address)),
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
