package docker

import (
  "os"
  "strings"
  "fmt"
)

type getEngineBaseUrlUseCase interface {
  Execute(
  ) (baseUrl string, err error)
}

func newGetEngineBaseUrlUseCase(
) (getEngineBaseUrlUseCase getEngineBaseUrlUseCase) {

  getEngineBaseUrlUseCase = &_getEngineBaseUrlUseCase{}

  return

}

type _getEngineBaseUrlUseCase struct{}

func (this _getEngineBaseUrlUseCase) Execute(
) (baseUrl string, err error) {

  var hostname string
  dockerComposeHost, isDockerMachine := os.LookupEnv("DOCKER_HOST")
  if (isDockerMachine) {
    addrParts := strings.Split(dockerComposeHost, ":")
    hostname = strings.TrimPrefix(addrParts[1], "//")
  } else {
    hostname = "localhost"
  }

  baseUrl = fmt.Sprintf("http://%v:42224", hostname)

  return
}
