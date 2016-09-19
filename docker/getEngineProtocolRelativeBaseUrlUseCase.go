package docker

import (
  "os"
  "strings"
  "fmt"
)

type getEngineProtocolRelativeBaseUrlUseCase interface {
  Execute(
  ) (protocolRelativeBaseUrl string, err error)
}

func newGetEngineProtocolRelativeBaseUrlUseCase(
) (getEngineProtocolRelativeBaseUrlUseCase getEngineProtocolRelativeBaseUrlUseCase) {

  getEngineProtocolRelativeBaseUrlUseCase = &_getEngineProtocolRelativeBaseUrlUseCase{}

  return

}

type _getEngineProtocolRelativeBaseUrlUseCase struct{}

func (this _getEngineProtocolRelativeBaseUrlUseCase) Execute(
) (protocolRelativeBaseUrl string, err error) {

  var hostname string
  dockerComposeHost, isDockerMachine := os.LookupEnv("DOCKER_HOST")
  if (isDockerMachine) {
    addrParts := strings.Split(dockerComposeHost, ":")
    hostname = strings.TrimPrefix(addrParts[1], "//")
  } else {
    hostname = "localhost"
  }

  protocolRelativeBaseUrl = fmt.Sprintf("//%v:42224", hostname)

  return
}
