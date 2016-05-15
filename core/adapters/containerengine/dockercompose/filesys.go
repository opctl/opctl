package dockercompose

import (
  "io/ioutil"
  "path"
  "os"
)

type filesys interface {
  getPathToOpDockerComposeFile(
  opName string,
  ) (pathToOpDockerComposeFile string)

  isDockerComposeFileExistent(
  pathToOpDir string,
  ) (isExistent bool)

  saveOpDockerComposeFile(
  pathToOpDir string,
  data []byte,
  ) (err error)
}

type _filesys struct{}

func (this _filesys)  getPathToOpDockerComposeFile(
pathToOp string,
) (pathToOpDockerComposeFile string) {

  pathToOpDockerComposeFile = path.Join(pathToOp, "docker-compose.yml")

  return

}

func (this _filesys) isDockerComposeFileExistent(
pathToOpDir string,
) (isExistent bool) {

  pathToFile := this.getPathToOpDockerComposeFile(pathToOpDir)

  if _, err := os.Stat(pathToFile); err == nil {
    return true
  }

  return false

}

func (this _filesys)  saveOpDockerComposeFile(
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
