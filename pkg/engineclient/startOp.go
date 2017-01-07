package engineclient

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"io/ioutil"
	"net/http"
)

func (this _engineClient) StartOp(
	req model.StartOpReq,
) (
	opId string,
	err error,
) {

	reqBytes, err := this.jsonFormat.From(req)
	if nil != err {
		return
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http:%v/instances/starts", "localhost"),
		bytes.NewBuffer(reqBytes),
	)
	if nil != err {
		return
	}

	httpResp, err := this.httpClient.Do(httpReq)
	if nil != err {
		return
	}

	opIdBuffer, err := ioutil.ReadAll(httpResp.Body)
	if nil != err {
		return
	}

	opId = string(opIdBuffer)

	return

}
