/*
Package handler implements an http.Handler for an opspec node
*/
package handler

import (
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/golang-interfaces/ihttp"
	"github.com/golang-interfaces/ijson"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
)

func New(
	core core.Core,
) http.Handler {

	router := mux.NewRouter()

	router.UseEncodedPath()

	handler := _handler{
		core:   core,
		http:   ihttp.New(),
		json:   ijson.New(),
		Router: router,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	router.HandleFunc(api.URLEvents_Stream, handler.events_streams).Methods(http.MethodGet)

	router.HandleFunc(api.URLLiveness, handler.liveness).Methods(http.MethodGet)

	router.HandleFunc(api.URLOps_Kills, handler.ops_kills).Methods(http.MethodPost)

	router.HandleFunc(api.URLOps_Starts, handler.ops_starts).Methods(http.MethodPost)

	router.HandleFunc(api.URLPkgs_Ref_Contents, handler.pkgs_ref_contents).Methods(http.MethodGet)

	router.HandleFunc(api.URLPkgs_Ref_Contents_Path, handler.pkgs_ref_contents_path).Methods(http.MethodGet)

	return router

}

type _handler struct {
	core core.Core
	http ihttp.IHTTP
	json ijson.IJSON
	*mux.Router
	upgrader iwebsocket.Upgrader
}
