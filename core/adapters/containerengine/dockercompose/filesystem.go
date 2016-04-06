package dockercompose

import (
  "io/ioutil"
  "path"
)

type filesystem interface {
  getPathToOpDockerComposeFile(
  opName string,
  ) (pathToOpDockerComposeFile string)

  saveOpDockerComposeFile(
  pathToOpDir string,
  data []byte,
  ) (err error)
}

type filesystemImpl struct{}

func (this filesystemImpl)  getPathToOpDockerComposeFile(
pathToOp string,
) (pathToOpDockerComposeFile string) {

  pathToOpDockerComposeFile = path.Join(pathToOp, "docker-compose.yml")

  return

}

func (this filesystemImpl)  saveOpDockerComposeFile(
pathToOpDir string,
data []byte,
) (err error) {

  err = ioutil.WriteFile(
    this.getPathToOpDockerComposeFile(pathToOpDir),
    data,
    0777,
  )

  return

}