package git

type compositionRoot interface {
  GetTemplateUcExecuter() getTemplateUcExecuter
}

func newCompositionRoot(
) compositionRoot {

  return &_compositionRoot{
    getTemplateUcExecuter: newGetTemplateUcExecuter(),
  }

}

type _compositionRoot struct {
  getTemplateUcExecuter getTemplateUcExecuter
}

func (_compositionRoot _compositionRoot) GetTemplateUcExecuter(
) getTemplateUcExecuter {
  return _compositionRoot.getTemplateUcExecuter
}
