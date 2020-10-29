package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api"
	"io/ioutil"
	"net/http"
	"path"
)

func (c client) AddAuth(
	ctx context.Context,
	req model.AddAuthReq,
) error {

	reqBytes, err := json.Marshal(req)
	if nil != err {
		return err
	}

	reqURL := c.baseUrl
	reqURL.Path = path.Join(reqURL.Path, api.URLAuths_Adds)

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

	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if nil != err {
		return err
	}

	if http.StatusCreated != httpResp.StatusCode {
		return errors.New(string(bodyBytes))
	}

	return nil

}
