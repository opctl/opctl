package opspec

//go:generate counterfeiter -o ./fakeStartOpRunUseCase.go --fake-name fakeStartOpRunUseCase ./ startOpRunUseCase

import (
  "bytes"
  "github.com/opspec-io/sdk-golang/models"
  "net/http"
  "fmt"
  "github.com/opspec-io/sdk-golang/adapters"
  "io/ioutil"
)

type startOpRunUseCase interface {
  Execute(
  req models.StartOpRunReq,
  ) (
  opRunId string,
  err error,
  )
}

func newStartOpRunUseCase(
engineHost adapters.EngineHost,
httpClient httpClient,
jsonFormat format,
) startOpRunUseCase {

  return &_startOpRunUseCase{
    engineHost:engineHost,
    httpClient:httpClient,
    jsonFormat:jsonFormat,
  }

}

type _startOpRunUseCase struct {
  engineHost adapters.EngineHost
  httpClient httpClient
  jsonFormat format
}

func (this _startOpRunUseCase) Execute(
req models.StartOpRunReq,
) (
opRunId string,
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
