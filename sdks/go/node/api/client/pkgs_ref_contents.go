package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/opctl/opctl/sdks/go/node/api"
	"github.com/opctl/opctl/sdks/go/types"
)

func (c client) ListDescendants(
	ctx context.Context,
	req types.ListDescendantsReq,
) (
	[]*types.DirEntry,
	error,
) {

	// build url
	path := strings.Replace(api.URLPkgs_Ref_Contents, "{ref}", url.PathEscape(req.PkgRef), 1)

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

	defer httpResp.Body.Close()

	if httpResp.StatusCode >= 400 {
		switch httpResp.StatusCode {
		case http.StatusUnauthorized:
			return nil, types.ErrDataProviderAuthentication{}
		case http.StatusForbidden:
			return nil, types.ErrDataProviderAuthorization{}
		case http.StatusNotFound:
			return nil, types.ErrDataRefResolution{}
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

	var contentList []*types.DirEntry
	return contentList, json.NewDecoder(httpResp.Body).Decode(&contentList)

}
