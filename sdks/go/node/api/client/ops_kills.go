package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/api"
	"net/http"
	"path"
)

func (c client) KillOp(
	ctx context.Context,
	req model.KillOpReq,
) error {

	reqBytes, err := json.Marshal(req)
	if nil != err {
		return err
	}

	reqURL := c.baseUrl
	reqURL.Path = path.Join(reqURL.Path, api.URLOps_Kills)

	httpReq, err := http.NewRequest(
		"POST",
		reqURL.String(),
		bytes.NewBuffer(reqBytes),
	)
	if nil != err {
		return err
	}

	httpReq = httpReq.WithContext(ctx)

	httpResp, err := c.httpClient.Do(httpReq)
	if nil != err {
		return err
	}
	// don't leak resources
	defer httpResp.Body.Close()

	return nil

}
