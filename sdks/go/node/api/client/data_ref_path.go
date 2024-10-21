package client

import (
	"context"
	"net/url"
	"strings"

	"github.com/jfbus/httprs"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/urltemplates"
)

func (c apiClient) GetData(
	ctx context.Context,
	req model.GetDataReq,
) (
	model.ReadSeekCloser,
	error,
) {
	path := strings.Replace(urltemplates.Data_Ref, "{ref}", url.PathEscape(req.DataRef), 1)

	httpResp, err := c.getWithAuth(ctx, path, req.PullCreds)
	if err != nil {
		return nil, err
	}

	// @TODO: rework to be true read seek closer; httprs seems to do a call w/ no end range
	return httprs.NewHttpReadSeeker(httpResp), nil
}
