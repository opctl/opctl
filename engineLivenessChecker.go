package opspec

//go:generate counterfeiter -o ./fakeEngineLivenessChecker.go --fake-name fakeEngineLivenessChecker ./ engineLivenessChecker

import (
  "time"
  "errors"
  "github.com/opspec-io/sdk-golang/adapters"
  "net/http"
  "fmt"
  "bytes"
)

type engineLivenessChecker interface {
  Execute(
  ) (err error)
}

func newEngineLivenessChecker(
engineHost adapters.EngineHost,
httpClient httpClient,
) engineLivenessChecker {

  return &_engineLivenessChecker{
    engineHost:engineHost,
    httpClient:httpClient,
  }

}

type _engineLivenessChecker struct {
  engineHost adapters.EngineHost
  httpClient httpClient
}

func (this _engineLivenessChecker) Execute(
) (err error) {

  engineProtocolRelativeBaseUrl, err := this.engineHost.GetEngineProtocolRelativeBaseUrl()
  if (nil != err) {
    return
  }

  httpReq, err := http.NewRequest(
    "GET",
    fmt.Sprintf("http:%v/liveness", engineProtocolRelativeBaseUrl),
    bytes.NewBuffer([]byte{}),
  )
  if (nil != err) {
    return
  }

  timeout := time.After(15 * time.Second)
  tick := time.Tick(500 * time.Millisecond)

  // retry until we're timed out or engine is alive
  for {
    select {
    // Got a timeout! fail with a timeout error
    case <-timeout:
      err = errors.New("Timeout exceeded while liveness checking engine. \n")
      return
    // Got a tick, check liveness
    case <-tick:
      _, err = this.httpClient.Do(httpReq)
      if err == nil {
        return
      }
    }
  }

  return

}
