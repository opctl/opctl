package docker

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
  opBundlePath string,
  ) (isExistent bool)

  saveOpDockerComposeFile(
  opBundlePath string,
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
opBundlePath string,
) (isExistent bool) {

  pathToFile := this.getPathToOpDockerComposeFile(opBundlePath)

  if _, err := os.Stat(pathToFile); err == nil {
    return true
  }

  return false

}

func (this _filesys)  saveOpDockerComposeFile(
opBundlePath string,
data []byte,
) (err error) {

  err = ioutil.WriteFile(
    this.getPathToOpDockerComposeFile(opBundlePath),
    data,
    0777,
  )

  return

}
