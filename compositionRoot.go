package opspec

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

type compositionRoot interface {
  CreateOpUseCase() createOpUseCase
  SetCollectionDescriptionUseCase() setCollectionDescriptionUseCase
  SetOpDescriptionUseCase() setOpDescriptionUseCase
  TryResolveDefaultCollectionUseCase() tryResolveDefaultCollectionUseCase
}

func newCompositionRoot(
filesystem Filesystem,
) (compositionRoot compositionRoot) {

  yamlCodec := newYamlCodec()

  createOpUseCase := newCreateOpUseCase(
    filesystem,
    yamlCodec,
  )

  setCollectionDescriptionUseCase := newSetCollectionDescriptionUseCase(
    filesystem,
    yamlCodec,
  )

  setOpDescriptionUseCase := newSetOpDescriptionUseCase(
    filesystem,
    yamlCodec,
  )

  tryResolveDefaultCollectionUseCase := newTryResolveDefaultCollectionUseCase(
    filesystem,
  )

  compositionRoot = &_compositionRoot{
    createOpUseCase:createOpUseCase,
    setCollectionDescriptionUseCase:setCollectionDescriptionUseCase,
    setOpDescriptionUseCase: setOpDescriptionUseCase,
    tryResolveDefaultCollectionUseCase: tryResolveDefaultCollectionUseCase,
  }

  return

}

type _compositionRoot struct {
  createOpUseCase                    createOpUseCase
  setCollectionDescriptionUseCase    setCollectionDescriptionUseCase
  setOpDescriptionUseCase            setOpDescriptionUseCase
  tryResolveDefaultCollectionUseCase tryResolveDefaultCollectionUseCase
}

func (this _compositionRoot) CreateOpUseCase(
) createOpUseCase {
  return this.createOpUseCase
}

func (this _compositionRoot) SetCollectionDescriptionUseCase(
) setCollectionDescriptionUseCase {
  return this.setCollectionDescriptionUseCase
}

func (this _compositionRoot) SetOpDescriptionUseCase(
) setOpDescriptionUseCase {
  return this.setOpDescriptionUseCase
}

func (this _compositionRoot) TryResolveDefaultCollectionUseCase(
) tryResolveDefaultCollectionUseCase {
  return this.tryResolveDefaultCollectionUseCase
}
