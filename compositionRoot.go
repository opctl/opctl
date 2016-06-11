package sdk

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

type compositionRoot interface {
  CreateOpUseCase() createOpUseCase
  SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase
}

func newCompositionRoot(
filesystem Filesystem,
) (compositionRoot compositionRoot) {

  yamlCodec := newYamlCodec()

  createOpUseCase := newCreateOpUseCase(
    filesystem,
    yamlCodec,
  )

  setDescriptionOfOpUseCase := newSetDescriptionOfOpUseCase(
    filesystem,
    yamlCodec,
  )

  compositionRoot = &_compositionRoot{
    createOpUseCase:createOpUseCase,
    setDescriptionOfOpUseCase: setDescriptionOfOpUseCase,
  }

  return

}

type _compositionRoot struct {
  createOpUseCase              createOpUseCase
  setDescriptionOfOpUseCase setDescriptionOfOpUseCase
}

func (this _compositionRoot) CreateOpUseCase() createOpUseCase {
  return this.createOpUseCase
}

func (this _compositionRoot) SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase {
  return this.setDescriptionOfOpUseCase
}
