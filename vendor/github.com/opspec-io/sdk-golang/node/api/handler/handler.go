/*
Package handler implements an http.Handler for an opspec node
*/
package handler

import (
	"github.com/gorilla/mux"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
)

func New(
	core core.Core,
) http.Handler {

	router := mux.NewRouter()

	router.Handle(
		api.Events_StreamsURLTpl,
		newGetEventStreamHandler(core),
	).Methods(http.MethodGet)

	router.Handle(
		api.LivenessURLTpl,
		newGetLivenessHandler(),
	).Methods(http.MethodGet)

	router.Handle(
		api.Ops_KillsURLTpl,
		newKillOpHandler(core),
	).Methods(http.MethodPost)

	router.Handle(
		api.Ops_StartsURLTpl,
		newStartOpHandler(core),
	).Methods(http.MethodPost)

	return router

}
