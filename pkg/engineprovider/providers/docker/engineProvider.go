package docker

import (
  "github.com/opspec-io/sdk-golang/pkg/engineprovider"
)

func New(
) (engineProvider engineprovider.EngineProvider) {

  engineProvider = &_engineProvider{
    containerChecker:newContainerChecker(),
    containerRemover:newContainerRemover(),
    containerStarter:newContainerStarter(newPathNormalizer()),
  }

  return

}

type _engineProvider struct {
  containerChecker containerChecker
  containerRemover containerRemover
  containerStarter containerStarter
}
