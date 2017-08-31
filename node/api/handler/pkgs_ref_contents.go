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

	pkgHandle, err := hdlr.core.ResolvePkg(
		pkgRef,
		pullCreds,
	)
	if nil != err {
		var status int
		switch err.(type) {
		case model.ErrPkgPullAuthentication:
			status = http.StatusUnauthorized
		case model.ErrPkgPullAuthorization:
			status = http.StatusForbidden
		case model.ErrPkgNotFound:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		http.Error(httpResp, err.Error(), status)
		return
	}

	pkgContentsList, err := pkgHandle.ListContents(
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
