package dockercompose

import (
  "io/ioutil"
  "path"
)

type filesystem interface {
  getPathToOperationDockerComposeFile(
  operationName string,
  ) (pathToOperationDockerComposeFile string)

  saveOperationDockerComposeFile(
  pathToOperationDir string,
  data []byte,
  ) (err error)
}

type filesystemImpl struct{}

func (this filesystemImpl)  getPathToOperationDockerComposeFile(
pathToOperation string,
) (pathToOperationDockerComposeFile string) {

  pathToOperationDockerComposeFile = path.Join(pathToOperation, "docker-compose.yml")

  return

}

func (this filesystemImpl)  saveOperationDockerComposeFile(
pathToOperationDir string,
data []byte,
) (err error) {

  err = ioutil.WriteFile(
    this.getPathToOperationDockerComposeFile(pathToOperationDir),
    data,
    0777,
  )

  return

}