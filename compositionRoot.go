package sdk_golang

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

type compositionRoot interface {
  SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase
}

func newCompositionRoot(
filesystem Filesystem,
) (compositionRoot compositionRoot, err error) {

  yamlCodec := newYamlCodec()

  setDescriptionOfOpUseCase := newSetDescriptionOfOpUseCase(
    filesystem,
    yamlCodec,
  )

  compositionRoot = &_compositionRoot{
    setDescriptionOfOpUseCase: setDescriptionOfOpUseCase,
  }

  return

}

type _compositionRoot struct {
  setDescriptionOfOpUseCase setDescriptionOfOpUseCase
}

func (this _compositionRoot) SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase {
  return this.setDescriptionOfOpUseCase
}
