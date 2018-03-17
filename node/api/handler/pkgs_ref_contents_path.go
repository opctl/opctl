package handler

import (
	"github.com/gorilla/mux"
	"github.com/opspec-io/sdk-golang/model"
	"net/http"
	"net/url"
	"path"
	"time"
)

func (hdlr _handler) pkgs_ref_contents_path(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	vars := mux.Vars(httpReq)

	pkgRef, err := url.PathUnescape(vars["ref"])
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	contentPath, err := url.PathUnescape(vars["path"])
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

	pkgContentReader, err := opDirHandle.GetContent(
		httpReq.Context(),
		contentPath,
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusInternalServerError)
		return
	}

	hdlr.http.ServeContent(
		httpResp,
		httpReq,
		path.Base(contentPath),
		time.Time{},
		pkgContentReader,
	)
}
