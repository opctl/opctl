package opspec

import (
  "github.com/opspec-io/sdk-golang/models"
  "net/http"
  "github.com/opspec-io/sdk-golang/adapters"
  "fmt"
  "bytes"
)

//go:generate counterfeiter -o ./fakeKillOpRunUseCase.go --fake-name fakeKillOpRunUseCase ./ killOpRunUseCase


type killOpRunUseCase interface {
  Execute(
  req models.KillOpRunReq,
  ) (
  err error,
  )
}

func newKillOpRunUseCase(
engineHost adapters.EngineHost,
httpClient httpClient,
jsonFormat format,
) killOpRunUseCase {

  return &_killOpRunUseCase{
    engineHost:engineHost,
    httpClient:httpClient,
    jsonFormat:jsonFormat,
  }

}

type _killOpRunUseCase struct {
  engineHost adapters.EngineHost
  httpClient httpClient
  jsonFormat format
}

func (this _killOpRunUseCase) Execute(
req models.KillOpRunReq,
) (
err error,
) {

  engineProtocolRelativeBaseUrl, err := this.engineHost.GetEngineProtocolRelativeBaseUrl()
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
