package docker

import (
  "github.com/opspec-io/sdk-golang/adapters"
)

func New(
) (engineHost adapters.EngineHost) {

  engineHost = &_engineHost{
    compositionRoot:newCompositionRoot(),
  }

  return

}

type _engineHost struct {
  compositionRoot compositionRoot
}

func (this _engineHost) EnsureEngineRunning(
image string,
) (err error) {

  return this.
  compositionRoot.
    EnsureEngineRunningUseCase().
    Execute(image)

}

func (this _engineHost) GetEngineBaseUrl(
) (
baseUrl string,
err error,
) {

  return this.
  compositionRoot.
    GetEngineBaseUrlUseCase().
    Execute()

}
