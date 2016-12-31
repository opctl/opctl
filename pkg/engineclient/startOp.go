package engineclient

import (
  "bytes"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "net/http"
  "fmt"
  "io/ioutil"
)

func (this _engineClient) StartOp(
req model.StartOpReq,
) (
opId string,
err error,
) {

  engineProtocolRelativeBaseUrl, err := this.engineProvider.GetEngineProtocolRelativeBaseUrl()
  if (nil != err) {
    return
  }

  reqBytes, err := this.jsonFormat.From(req)
  if (nil != err) {
    return
  }

  httpReq, err := http.NewRequest(
    "POST",
    fmt.Sprintf("http:%v/instances/starts", engineProtocolRelativeBaseUrl),
    bytes.NewBuffer(reqBytes),
  )
  if (nil != err) {
    return
  }

  httpResp, err := this.httpClient.Do(httpReq)
  if (nil != err) {
    return
  }

  opIdBuffer, err := ioutil.ReadAll(httpResp.Body)
  if (nil != err) {
    return
  }

  opId = string(opIdBuffer)

  return

}
