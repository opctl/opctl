package engineclient

import (
  "github.com/opspec-io/sdk-golang/models"
  "net/http"
  "fmt"
  "bytes"
)

func (this _engineClient) KillOpRun(
req models.KillOpRunReq,
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
    fmt.Sprintf("http:%v/op-run-kills", engineProtocolRelativeBaseUrl),
    bytes.NewBuffer(reqBytes),
  )
  if (nil != err) {
    return
  }

  _, err = this.httpClient.Do(httpReq)
  return

}
