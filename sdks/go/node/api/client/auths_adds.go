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

func (c apiClient) AddAuth(
	ctx context.Context,
	req model.AddAuthReq,
) error {

	reqBytes, err := json.Marshal(req)
	if nil != err {
		return err
	}

	reqURL := c.baseURL
	reqURL.Path = path.Join(reqURL.Path, api.URLAuths_Adds)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		reqURL.String(),
		bytes.NewBuffer(reqBytes),
	)
	if nil != err {
		return err
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if nil != err {
		return err
	}
	// don't leak resources
	defer httpResp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if nil != err {
		return err
	}

	if http.StatusCreated != httpResp.StatusCode {
		return errors.New(string(bodyBytes))
	}

	return nil

}
