package client

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/opctl/opctl/sdks/go/node/api"
)

func (c client) Liveness(
	ctx context.Context,
) error {

	httpReq, err := http.NewRequest(
		"GET",
		c.baseUrl.String()+api.URLLiveness,
		nil,
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

	if http.StatusOK != httpResp.StatusCode {
		return errors.New(string(bodyBytes))
	}

	return nil
}
