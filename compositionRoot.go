package opspec

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

type compositionRoot interface {
  CreateCollectionUseCase() createCollectionUseCase
  CreateOpUseCase() createOpUseCase
  GetCollectionUseCase() getCollectionUseCase
  GetOpUseCase() getOpUseCase
  SetCollectionDescriptionUseCase() setCollectionDescriptionUseCase
  SetOpDescriptionUseCase() setOpDescriptionUseCase
  TryResolveDefaultCollectionUseCase() tryResolveDefaultCollectionUseCase
}

func newCompositionRoot(
filesystem Filesystem,
) (compositionRoot compositionRoot) {

  yamlCodec := newYamlCodec()

  opViewFactory := newOpViewFactory(
    filesystem,
    yamlCodec,
  )

  collectionViewFactory := newCollectionViewFactory(
    filesystem,
    opViewFactory,
    yamlCodec,
  )

  createCollectionUseCase := newCreateCollectionUseCase(
    filesystem,
    yamlCodec,
  )

  createOpUseCase := newCreateOpUseCase(
    filesystem,
    yamlCodec,
  )

  getCollectionUseCase := newGetCollectionUseCase(collectionViewFactory)

  getOpUseCase := newGetOpUseCase(opViewFactory)

  setCollectionDescriptionUseCase := newSetCollectionDescriptionUseCase(
    filesystem,
    yamlCodec,
  )

  setOpDescriptionUseCase := newSetOpDescriptionUseCase(
    filesystem,
    yamlCodec,
  )

  tryResolveDefaultCollectionUseCase := newTryResolveDefaultCollectionUseCase(filesystem)

  compositionRoot = &_compositionRoot{
    createCollectionUseCase:createCollectionUseCase,
    createOpUseCase:createOpUseCase,
    getCollectionUseCase:getCollectionUseCase,
    getOpUseCase:getOpUseCase,
    setCollectionDescriptionUseCase:setCollectionDescriptionUseCase,
    setOpDescriptionUseCase: setOpDescriptionUseCase,
    tryResolveDefaultCollectionUseCase: tryResolveDefaultCollectionUseCase,
  }

  return

}

type _compositionRoot struct {
  createCollectionUseCase            createCollectionUseCase
  createOpUseCase                    createOpUseCase
  getCollectionUseCase               getCollectionUseCase
  getOpUseCase                       getOpUseCase
  setCollectionDescriptionUseCase    setCollectionDescriptionUseCase
  setOpDescriptionUseCase            setOpDescriptionUseCase
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

func (this _compositionRoot) GetOpUseCase(
) getOpUseCase {
  return this.getOpUseCase
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
