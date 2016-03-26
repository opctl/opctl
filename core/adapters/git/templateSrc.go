package git

import "github.com/dev-op-spec/engine/core/ports"

func NewTemplateSrc() ports.TemplateSrc {

  return templateSrc{
    compositionRoot:newCompositionRoot(),
  }

}

type templateSrc struct {
  compositionRoot compositionRoot
}

func (templateSrc templateSrc) GetTemplate(
templateRef string,
) (pathToTemplateRootDir string, err error) {
  return templateSrc.compositionRoot.
  GetTemplateUseCase().
  Execute(templateRef)
}
