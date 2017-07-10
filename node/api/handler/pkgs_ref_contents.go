package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

func (hdlr _handler) pkgs_ref_contents(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	vars := mux.Vars(httpReq)

	pkgRef, err := url.PathUnescape(vars["ref"])
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHandle, err := hdlr.core.ResolvePkg(
		pkgRef,
		nil,
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgContentsList, err := pkgHandle.ListContents()
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	// @TODO: replace w/ json.NewEncoder(w).Encode(p); more performant
	pkgContentsListBytes, err := hdlr.json.Marshal(pkgContentsList)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	httpResp.WriteHeader(http.StatusOK)
	httpResp.Header().Set("Content-Type", "application/json")
	httpResp.Write(pkgContentsListBytes)
}
