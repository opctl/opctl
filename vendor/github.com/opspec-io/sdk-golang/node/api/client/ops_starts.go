package client

import (
	"bytes"
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"io/ioutil"
	"net/http"
)

// StartOp starts an op & returns its root op id (ROID)
func (c client) StartOp(
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

	httpResp, err := c.httpClient.Do(httpReq)
	if nil != err {
		return "", err
	}

	opIdBuffer, err := ioutil.ReadAll(httpResp.Body)
	return string(opIdBuffer), nil

}
