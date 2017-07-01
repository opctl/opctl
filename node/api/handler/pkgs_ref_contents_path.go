package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"time"
)

func (hdlr _handler) pkgs_ref_contents_path(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	vars := mux.Vars(httpReq)
	contentPath := vars["path"]

	pkgContent, err := hdlr.core.GetPkgContent(
		vars["ref"],
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
		pkgContent,
	)
}
