package client

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api"
)

func (c apiClient) ListDescendants(
	ctx context.Context,
	req model.ListDescendantsReq,
) (
	[]*model.DirEntry,
	error,
) {
	path := strings.Replace(api.URLData_Ref, "{ref}", url.PathEscape(req.DataRef), 1)

	httpResp, err := c.getWithAuth(ctx, path, req.PullCreds)
	if err != nil {
		return nil, err
	}

	var contentList []*model.DirEntry
	return contentList, json.NewDecoder(httpResp.Body).Decode(&contentList)
}
