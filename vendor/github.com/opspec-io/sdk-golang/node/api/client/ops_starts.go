package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"io/ioutil"
	"net/http"
)

// StartOp starts an op & returns its root op id (ROId)
func (c client) StartOp(
	ctx context.Context,
	req model.StartOpReq,
) (string, error) {

	reqBytes, err := json.Marshal(req)
	if nil != err {
		return "", nil
	}

	reqUrl := c.baseUrl
	reqUrl.Path = api.Ops_StartsURLTpl

	httpReq, err := http.NewRequest(
		"POST",
		reqUrl.String(),
		bytes.NewBuffer(reqBytes),
	)
	if nil != err {
		return "", nil
	}

	httpReq.WithContext(ctx)

	httpResp, err := c.httpClient.Do(httpReq)
	if nil != err {
		return "", err
	}
	// don't leak resources
	defer httpResp.Body.Close()

	opIdBuffer, err := ioutil.ReadAll(httpResp.Body)
	return string(opIdBuffer), nil

}
