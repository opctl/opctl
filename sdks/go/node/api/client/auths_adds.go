package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
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
	if err != nil {
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
	if err != nil {
		return err
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	// don't leak resources
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	if http.StatusCreated != httpResp.StatusCode {
		return errors.New(string(bodyBytes))
	}

	return nil

}
