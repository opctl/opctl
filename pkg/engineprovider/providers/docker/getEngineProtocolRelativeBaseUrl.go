package docker

import (
  "os"
  "strings"
  "fmt"
)

func (this _engineProvider) GetEngineProtocolRelativeBaseUrl(
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
