package consumenodeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
	"net/http"
)

func (this consumeNodeApi) StartOp(
	req model.StartOpReq,
) (
	opId string,
	err error,
) {

	reqBytes, err := json.Marshal(req)
	if nil != err {
		return
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%v/ops/starts", "localhost:42224"),
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
