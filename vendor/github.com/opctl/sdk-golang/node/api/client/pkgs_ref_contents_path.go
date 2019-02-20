package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfbus/httprs"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/api"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (c client) GetData(
	ctx context.Context,
	req model.GetDataReq,
) (
	model.ReadSeekCloser,
	error,
) {

	// build url
	path := strings.Replace(api.URLPkgs_Ref_Contents_Path, "{ref}", url.PathEscape(req.PkgRef), 1)
	path = strings.Replace(path, "{path}", url.PathEscape(req.ContentPath), 1)

	httpReq, err := http.NewRequest(
		"GET",
		c.baseUrl.String()+path,
		nil,
	)
	if nil != err {
		return nil, err
	}

	httpReq = httpReq.WithContext(ctx)
	if nil != req.PullCreds {
		httpReq.SetBasicAuth(
			req.PullCreds.Username,
			req.PullCreds.Password,
		)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if nil != err {
		return nil, err
	}

	if httpResp.StatusCode >= 400 {
		defer httpResp.Body.Close()

		switch httpResp.StatusCode {
		case http.StatusUnauthorized:
			return nil, model.ErrDataProviderAuthentication{}
		case http.StatusForbidden:
			return nil, model.ErrDataProviderAuthorization{}
		case http.StatusNotFound:
			return nil, model.ErrDataRefResolution{}
		default:
			body, err := ioutil.ReadAll(httpResp.Body)
			if nil != err {
				return nil, fmt.Errorf(
					"Error encountered parsing response w/ status code '%v'; error was %v",
					httpResp.StatusCode,
					err.Error(),
				)
			}
			return nil, errors.New(string(body))
		}
	}

	// @TODO: rework to be true read seek closer; httprs seems to do a call w/ no end range
	return httprs.NewHttpReadSeeker(httpResp), nil

}
