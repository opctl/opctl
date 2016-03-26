package dockercompose

import (
  "io/ioutil"
  "path"
  "errors"
)

type filesystem interface {
  getRelPathToDevOpDir(
  devOpName string,
  ) (relPathToDevOpDir string, err error)

  getRelPathToDevOpDockerComposeFile(
  devOpName string,
  ) (relPathToDevOpDockerComposeFile string, err error)

  saveDevOpDockerComposeFile(
  devOpName string,
  data []byte,
  ) (err error)
}

const (
  relPathToDevOpSpecDir = "./.dev-op-spec"

  relPathToDevOpsDir = relPathToDevOpSpecDir + "/dev-ops"
)

type filesystemImpl struct{}

func (fs filesystemImpl)  getRelPathToDevOpDir(
devOpName string,
) (relPathToDevOpDir string, err error) {

  if ("" == devOpName) {
    err= errors.New("devOpName cannot be nil")
  }

  relPathToDevOpDir = path.Join(relPathToDevOpsDir, devOpName)

  return
}

func (fs filesystemImpl)  getRelPathToDevOpDockerComposeFile(
devOpName string,
) (relPathToDevOpDockerComposeFile string, err error) {

  var relPathToDevOpDir string
  relPathToDevOpDir, err= fs.getRelPathToDevOpDir(devOpName)
  if (nil != err) {
    return
  }

  relPathToDevOpDockerComposeFile = path.Join(relPathToDevOpDir, "docker-compose.yml")

  return

}

func (fs filesystemImpl)  saveDevOpDockerComposeFile(
devOpName string,
data []byte,
) (err error) {

  var relPathToDevOpDockerComposeFile string
  relPathToDevOpDockerComposeFile, err= fs.getRelPathToDevOpDockerComposeFile(devOpName)
  if (nil != err) {
    return
  }

  err= ioutil.WriteFile(
    relPathToDevOpDockerComposeFile,
    data,
    0777,
  )

  return

}