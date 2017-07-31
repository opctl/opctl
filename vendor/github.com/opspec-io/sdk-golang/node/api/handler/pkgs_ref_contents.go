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

	pkgContentsList, err := pkgHandle.ListContents(
		httpReq.Context(),
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := hdlr.json.NewEncoder(httpResp).Encode(pkgContentsList); nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}
	httpResp.Header().Set("Content-Type", "application/json")
}
