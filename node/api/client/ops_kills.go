package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/http"
)

func (c client) KillOp(
	ctx context.Context,
	req model.KillOpReq,
) error {

	reqBytes, err := json.Marshal(req)
	if nil != err {
		return err
	}

	reqUrl := c.baseUrl
	reqUrl.Path = api.URLOps_Kills

	httpReq, err := http.NewRequest(
		"POST",
		reqUrl.String(),
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
