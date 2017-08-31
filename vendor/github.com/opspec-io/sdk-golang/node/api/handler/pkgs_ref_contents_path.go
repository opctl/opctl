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

	pkgHandle, err := hdlr.core.ResolvePkg(
		pkgRef,
		nil,
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

	pkgContentReader, err := pkgHandle.GetContent(
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
