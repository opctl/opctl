package sdk

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

type compositionRoot interface {
  AddOpUseCase() addOpUseCase
  SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase
}

func newCompositionRoot(
filesystem Filesystem,
) (compositionRoot compositionRoot) {

  yamlCodec := newYamlCodec()

  addOpUseCase := newAddOpUseCase(
    filesystem,
    yamlCodec,
  )

  setDescriptionOfOpUseCase := newSetDescriptionOfOpUseCase(
    filesystem,
    yamlCodec,
  )

  compositionRoot = &_compositionRoot{
    addOpUseCase:addOpUseCase,
    setDescriptionOfOpUseCase: setDescriptionOfOpUseCase,
  }

  return

}

type _compositionRoot struct {
  addOpUseCase              addOpUseCase
  setDescriptionOfOpUseCase setDescriptionOfOpUseCase
}

func (this _compositionRoot) AddOpUseCase() addOpUseCase {
  return this.addOpUseCase
}

func (this _compositionRoot) SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase {
  return this.setDescriptionOfOpUseCase
}
