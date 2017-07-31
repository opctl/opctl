package client

import (
	"context"
	"github.com/jfbus/httprs"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/http"
	"net/url"
	"strings"
)

func (c client) GetPkgContent(
	ctx context.Context,
	req model.GetPkgContentReq,
) (
	model.ReadSeekCloser,
	error,
) {

	// build url
	path := strings.Replace(api.URLPkgs_Ref_Contents_Path, "{ref}", url.PathEscape(req.PkgRef), 1)
	path = strings.Replace(path, "{path}", url.PathEscape(req.ContentPath), 1)

	httpReq, err := http.NewRequest(
		"GET",
		c.baseUrl.String()+path,
		nil,
	)
	if nil != err {
		return nil, err
	}

	httpReq = httpReq.WithContext(ctx)

	httpResp, err := c.httpClient.Do(httpReq)
	if nil != err {
		return nil, err
	}

	// @TODO: rework to be true read seek closer; httprs seems to do a call w/ no end range
	return httprs.NewHttpReadSeeker(httpResp), nil

}
