package client

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/opctl/opctl/sdks/go/node/api"
)

func (c apiClient) Liveness(
	ctx context.Context,
) error {

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL.String()+api.URLLiveness,
		nil,
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

	if http.StatusOK != httpResp.StatusCode {
		return errors.New(string(bodyBytes))
	}

	return nil
}
