package ports

// influence http://bower.io/docs/pluggable-resolvers/
type TemplateSrc interface {
  GetTemplate(
  templateRef string,
  ) (pathToTemplateRootDir string, err error)
}
