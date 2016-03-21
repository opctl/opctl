package git

type getTemplateUcExecuter interface {

  Execute(
  templateRef string,
  ) (pathToTemplateRootDir string, err error)

}

func newGetTemplateUcExecuter() getTemplateUcExecuter {

  return &getTemplateUcExecuterImpl{}

}

type getTemplateUcExecuterImpl struct {}

func (uc getTemplateUcExecuterImpl) Execute(
templateRef string,
) (pathToTemplateRootDir string, err error) {

  return

}
