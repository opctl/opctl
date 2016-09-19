package opspec

import (
  "github.com/opspec-io/sdk-golang/adapters"
  "net/http"
  "time"
)

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

type compositionRoot interface {
  CreateCollectionUseCase() createCollectionUseCase
  CreateOpUseCase() createOpUseCase
  GetCollectionUseCase() getCollectionUseCase
  GetEventStreamUseCase() getEventStreamUseCase
  GetOpUseCase() getOpUseCase
  KillOpRunUseCase() killOpRunUseCase
  SetCollectionDescriptionUseCase() setCollectionDescriptionUseCase
  SetOpDescriptionUseCase() setOpDescriptionUseCase
  StartOpRunUseCase() startOpRunUseCase
  TryResolveDefaultCollectionUseCase() tryResolveDefaultCollectionUseCase
}

func newCompositionRoot(
engineHost adapters.EngineHost,
filesystem filesystem,
) (compositionRoot compositionRoot) {

  yamlFormat := newYamlFormat()
  jsonFormat := newJsonFormat()

  httpClient := http.DefaultClient
  httpClient.Timeout = time.Second * 10

  opViewFactory := newOpViewFactory(
    filesystem,
    yamlFormat,
  )

  collectionViewFactory := newCollectionViewFactory(
    filesystem,
    opViewFactory,
    yamlFormat,
  )

  compositionRoot = &_compositionRoot{
    createCollectionUseCase:newCreateCollectionUseCase(filesystem, yamlFormat),
    createOpUseCase:newCreateOpUseCase(filesystem, yamlFormat),
    getCollectionUseCase:newGetCollectionUseCase(collectionViewFactory),
    getEventStreamUseCase:newGetEventStreamUseCase(engineHost, jsonFormat),
    getOpUseCase:newGetOpUseCase(opViewFactory),
    killOpRunUseCase:newKillOpRunUseCase(engineHost, httpClient, jsonFormat),
    setCollectionDescriptionUseCase:newSetCollectionDescriptionUseCase(filesystem, yamlFormat),
    setOpDescriptionUseCase: newSetOpDescriptionUseCase(filesystem, yamlFormat),
    startOpRunUseCase: newStartOpRunUseCase(engineHost, httpClient, jsonFormat),
    tryResolveDefaultCollectionUseCase: newTryResolveDefaultCollectionUseCase(filesystem),
  }

  return

}

type _compositionRoot struct {
  createCollectionUseCase            createCollectionUseCase
  createOpUseCase                    createOpUseCase
  getCollectionUseCase               getCollectionUseCase
  getEventStreamUseCase              getEventStreamUseCase
  getOpUseCase                       getOpUseCase
  killOpRunUseCase                   killOpRunUseCase
  setCollectionDescriptionUseCase    setCollectionDescriptionUseCase
  setOpDescriptionUseCase            setOpDescriptionUseCase
  startOpRunUseCase                  startOpRunUseCase
  tryResolveDefaultCollectionUseCase tryResolveDefaultCollectionUseCase
}

func (this _compositionRoot) CreateCollectionUseCase(
) createCollectionUseCase {
  return this.createCollectionUseCase
}

func (this _compositionRoot) CreateOpUseCase(
) createOpUseCase {
  return this.createOpUseCase
}

func (this _compositionRoot) GetCollectionUseCase(
) getCollectionUseCase {
  return this.getCollectionUseCase
}

func (this _compositionRoot) GetEventStreamUseCase(
) getEventStreamUseCase {
  return this.getEventStreamUseCase
}

func (this _compositionRoot) GetOpUseCase(
) getOpUseCase {
  return this.getOpUseCase
}

func (this _compositionRoot) KillOpRunUseCase() killOpRunUseCase {
  return this.killOpRunUseCase
}

func (this _compositionRoot) SetCollectionDescriptionUseCase(
) setCollectionDescriptionUseCase {
  return this.setCollectionDescriptionUseCase
}

func (this _compositionRoot) SetOpDescriptionUseCase(
) setOpDescriptionUseCase {
  return this.setOpDescriptionUseCase
}

func (this _compositionRoot) StartOpRunUseCase(
) startOpRunUseCase {
  return this.startOpRunUseCase
}

func (this _compositionRoot) TryResolveDefaultCollectionUseCase(
) tryResolveDefaultCollectionUseCase {
  return this.tryResolveDefaultCollectionUseCase
}
