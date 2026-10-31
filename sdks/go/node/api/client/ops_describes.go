package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/urltemplates"
)

func (c apiClient) DescribeOp(ctx context.Context, req model.DescribeOpReq) (*model.OpSpec, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	url := c.baseURL
	url.Path = path.Join(url.Path, urltemplates.Ops_Describes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("describe op: %s", body)
	}
	var spec model.OpSpec
	if err := json.NewDecoder(response.Body).Decode(&spec); err != nil {
		return nil, err
	}
	return &spec, nil
}
