package engineclient

import (
  "bytes"
  "github.com/opspec-io/sdk-golang/models"
  "net/http"
  "fmt"
  "io/ioutil"
)

func (this _engineClient) StartOpRun(
req models.StartOpRunReq,
) (
opRunId string,
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
    fmt.Sprintf("http:%v/op-runs", engineProtocolRelativeBaseUrl),
    bytes.NewBuffer(reqBytes),
  )
  if (nil != err) {
    return
  }

  httpResp, err := this.httpClient.Do(httpReq)
  if (nil != err) {
    return
  }

  opRunIdBuffer, err := ioutil.ReadAll(httpResp.Body)
  if (nil != err) {
    return
  }

  opRunId = string(opRunIdBuffer)

  return

}
