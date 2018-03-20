package handler

import (
	"github.com/gorilla/mux"
	"github.com/opspec-io/sdk-golang/model"
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

	var pullCreds *model.PullCreds
	pullUsername, pullPassword, hasBasicAuth := httpReq.BasicAuth()
	if hasBasicAuth {
		pullCreds = &model.PullCreds{
			Username: pullUsername,
			Password: pullPassword,
		}
	}

	opDirHandle, err := hdlr.core.ResolvePkg(
		httpReq.Context(),
		pkgRef,
		pullCreds,
	)
	if nil != err {
		var status int
		switch err.(type) {
		case model.ErrDataProviderAuthentication:
			status = http.StatusUnauthorized
		case model.ErrDataProviderAuthorization:
			status = http.StatusForbidden
		case model.ErrDataRefResolution:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		http.Error(httpResp, err.Error(), status)
		return
	}

	pkgContentsList, err := opDirHandle.ListContents(
		httpReq.Context(),
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	httpResp.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := hdlr.json.NewEncoder(httpResp).Encode(pkgContentsList); nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}
}
