package dockercompose

import (
  "os"
  "os/exec"
)

type devOpResourceFlusher interface {
  flush(devOpName string) (err error)
}

func newDevOpResourceFlusher(
fs filesystem,
) devOpResourceFlusher {

  return &devOpResourceFlusherImpl{
    fs:fs,
  }

}

type devOpResourceFlusherImpl struct {
  fs filesystem
}

func (f devOpResourceFlusherImpl) flush(
devOpName string,
) (err error) {

  var relPathToDevOpDockerComposeFile string
  relPathToDevOpDockerComposeFile, err = f.fs.getRelPathToDevOpDockerComposeFile(devOpName)
  if (nil != err) {
    return
  }

  // down
  dockerComposeDownCmd :=
  exec.Command(
    "docker-compose",
    "-f",
    relPathToDevOpDockerComposeFile,
    "down",
    "--rmi",
    "local",
    "-v",
  )
  dockerComposeDownCmd.Stdout = os.Stdout
  dockerComposeDownCmd.Stderr = os.Stderr
  err = dockerComposeDownCmd.Run()

  return

}
