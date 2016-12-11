package engineclient

import (
  "github.com/opspec-io/sdk-golang/pkg/model"
  "net/http"
  "fmt"
  "bytes"
)

func (this _engineClient) KillOp(
req model.KillOpReq,
) (
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
    fmt.Sprintf("http:%v/instances/kills", engineProtocolRelativeBaseUrl),
    bytes.NewBuffer(reqBytes),
  )
  if (nil != err) {
    return
  }

  _, err = this.httpClient.Do(httpReq)
  return

}
