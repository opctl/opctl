package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api"
)

// StartOp starts an op & returns its root op id (ROId)
func (c client) StartOp(
	ctx context.Context,
	req model.StartOpReq,
) (string, error) {

	reqBytes, err := json.Marshal(req)
	if nil != err {
		return "", err
	}

	reqURL := c.baseUrl
	reqURL.Path = path.Join(reqURL.Path, api.URLOps_Starts)

	httpReq, err := http.NewRequest(
		"POST",
		reqURL.String(),
		bytes.NewBuffer(reqBytes),
	)
	if nil != err {
		return "", err
	}

	httpReq = httpReq.WithContext(ctx)

	httpResp, err := c.httpClient.Do(httpReq)
	if nil != err {
		return "", err
	}
	// don't leak resources
	defer httpResp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if nil != err {
		return "", err
	}

	if http.StatusOK != httpResp.StatusCode {
		return "", errors.New(string(bodyBytes))
	}

	return string(bodyBytes), nil

}
