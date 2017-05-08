package client

import (
	"bytes"
	"encoding/json"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/http"
)

func (c client) KillOp(
	req model.KillOpReq,
) error {

	reqBytes, err := json.Marshal(req)
	if nil != err {
		return err
	}

	reqUrl := c.baseUrl
	reqUrl.Path = api.Ops_KillsURLTpl

	httpReq, err := http.NewRequest(
		"POST",
		reqUrl.String(),
		bytes.NewBuffer(reqBytes),
	)
	if nil != err {
		return err
	}

	_, err = c.httpClient.Do(httpReq)
	return err

}
