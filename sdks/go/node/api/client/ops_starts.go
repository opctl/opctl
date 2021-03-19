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
func (c apiClient) StartOp(
	ctx context.Context,
	req model.StartOpReq,
) (string, error) {

	// if remote node; need to embed local file/dir args
	if c.baseURL.Hostname() != "localhost" && c.baseURL.Hostname() != "127.0.0.1" {
		err := embedLocalFilesAndDirs(req.Args)
		if err != nil {
			return "", err
		}
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	reqURL := c.baseURL
	reqURL.Path = path.Join(reqURL.Path, api.URLOps_Starts)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		reqURL.String(),
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return "", err
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	// don't leak resources
	defer httpResp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}

	if http.StatusCreated != httpResp.StatusCode {
		return "", errors.New(string(bodyBytes))
	}

	return string(bodyBytes), nil

}
