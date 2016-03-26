package git

type compositionRoot interface {
  GetTemplateUseCase() getTemplateUseCase
}

func newCompositionRoot(
) compositionRoot {

  return &_compositionRoot{
    getTemplateUseCase: newGetTemplateUseCase(),
  }

}

type _compositionRoot struct {
  getTemplateUseCase getTemplateUseCase
}

func (this _compositionRoot) GetTemplateUseCase(
) getTemplateUseCase {
  return this.getTemplateUseCase
}
