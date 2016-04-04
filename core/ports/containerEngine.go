package ports

type ContainerEngine interface {
  InitOperation(
  pathToOperationDir string,
  ) (err error)

  RunOperation(
  pathToOperationDir string,
  ) (exitCode int, err error)
}
