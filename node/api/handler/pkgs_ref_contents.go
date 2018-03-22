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

	dataRef, err := url.PathUnescape(vars["ref"])
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

	opHandle, err := hdlr.core.ResolveData(
		httpReq.Context(),
		dataRef,
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

	dataNodesList, err := opHandle.ListDescendants(
		httpReq.Context(),
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	httpResp.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := hdlr.json.NewEncoder(httpResp).Encode(dataNodesList); nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}
}
