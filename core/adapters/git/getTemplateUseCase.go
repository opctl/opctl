package git

type getTemplateUseCase interface {

  Execute(
  templateRef string,
  ) (pathToTemplateRootDir string, err error)

}

func newGetTemplateUseCase() getTemplateUseCase {

  return &getTemplateUseCaseImpl{}

}

type getTemplateUseCaseImpl struct {}

func (uc getTemplateUseCaseImpl) Execute(
templateRef string,
) (pathToTemplateRootDir string, err error) {

  return

}
