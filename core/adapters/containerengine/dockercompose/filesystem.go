package dockercompose

import (
  "io/ioutil"
  "path"
)

type filesystem interface {
  getPathToDevOpDockerComposeFile(
  devOpName string,
  ) (pathToDevOpDockerComposeFile string)

  saveDevOpDockerComposeFile(
  pathToDevOpDir string,
  data []byte,
  ) (err error)
}

type filesystemImpl struct{}

func (this filesystemImpl)  getPathToDevOpDockerComposeFile(
pathToDevOp string,
) (pathToDevOpDockerComposeFile string) {

  pathToDevOpDockerComposeFile = path.Join(pathToDevOp, "docker-compose.yml")

  return

}

func (this filesystemImpl)  saveDevOpDockerComposeFile(
pathToDevOpDir string,
data []byte,
) (err error) {

  err = ioutil.WriteFile(
    this.getPathToDevOpDockerComposeFile(pathToDevOpDir),
    data,
    0777,
  )

  return

}