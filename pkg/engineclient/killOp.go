package engineclient

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"net/http"
)

func (this _engineClient) KillOp(
	req model.KillOpReq,
) (
	err error,
) {

	reqBytes, err := this.jsonFormat.From(req)
	if nil != err {
		return
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%v/instances/kills", "localhost"),
		bytes.NewBuffer(reqBytes),
	)
	if nil != err {
		return
	}

	_, err = this.httpClient.Do(httpReq)
	return

}
