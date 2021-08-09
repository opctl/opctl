package client

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c apiClient) getWithAuth(ctx context.Context, path string, pullCreds *model.Creds) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL.String()+path,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if pullCreds != nil {
		httpReq.SetBasicAuth(pullCreds.Username, pullCreds.Password)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	readBody := func() (string, error) {
		defer httpResp.Body.Close()
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return "", fmt.Errorf(
				"error encountered parsing response (status code %v): %w",
				httpResp.StatusCode,
				err,
			)
		}
		return string(body), nil
	}

	if httpResp.StatusCode >= 400 {
		switch httpResp.StatusCode {
		case http.StatusUnauthorized:
			body, err := readBody()
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("%w: %s", model.ErrDataProviderAuthentication{}, body)
		case http.StatusForbidden:
			body, err := readBody()
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("%w: %s", model.ErrDataProviderAuthorization{}, body)
		case http.StatusNotFound:
			return nil, model.ErrDataRefResolution{}
		default:
			body, err := readBody()
			if err != nil {
				return nil, err
			}
			return nil, errors.New(body)
		}
	}

	return httpResp, nil
}
