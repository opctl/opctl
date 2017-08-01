package client

import (
	"context"
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/http"
	"net/url"
	"strings"
)

func (c client) ListPkgContents(
	ctx context.Context,
	req model.ListPkgContentsReq,
) (
	[]*model.PkgContent,
	error,
) {

	// build url
	path := strings.Replace(api.URLPkgs_Ref_Contents, "{ref}", url.PathEscape(req.PkgRef), 1)

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

	defer httpResp.Body.Close()
	var contentList []*model.PkgContent

	return contentList, json.NewDecoder(httpResp.Body).Decode(&contentList)

}
