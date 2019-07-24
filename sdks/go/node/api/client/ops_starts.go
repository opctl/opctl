package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opctl/opctl/sdks/go/node/api"
	"github.com/opctl/opctl/sdks/go/types"
	"io/ioutil"
	"net/http"
	"path"
)

// StartOp starts an op & returns its root op id (ROId)
func (c client) StartOp(
	ctx context.Context,
	req types.StartOpReq,
) (string, error) {

	reqBytes, err := json.Marshal(req)
	if nil != err {
		return "", nil
	}

	reqURL := c.baseUrl
	reqURL.Path = path.Join(reqURL.Path, api.URLOps_Starts)

	httpReq, err := http.NewRequest(
		"POST",
		reqURL.String(),
		bytes.NewBuffer(reqBytes),
	)
	if nil != err {
		return "", nil
	}

	httpReq = httpReq.WithContext(ctx)

	httpResp, err := c.httpClient.Do(httpReq)
	if nil != err {
		return "", err
	}
	// don't leak resources
	defer httpResp.Body.Close()

	opIDBuffer, err := ioutil.ReadAll(httpResp.Body)
	return string(opIDBuffer), nil

}
