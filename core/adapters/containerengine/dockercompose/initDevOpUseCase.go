package dockercompose

import "path"

type initDevOpUseCase interface {
  Execute(
  pathToDevOpDir string,
  ) (err error)
}

func newInitDevOpUseCase(
fs filesystem,
yml yamlCodec,
) initDevOpUseCase {

  return &_initDevOpUseCase{
    fs:fs,
    yml:yml,
  }

}

type _initDevOpUseCase struct {
  fs  filesystem
  yml yamlCodec
}

func (this _initDevOpUseCase) Execute(
pathToDevOpDir string,
) (err error) {

  devOpName := path.Base(pathToDevOpDir)

  var dockerComposeFile =
  dockerComposeFile{
    Version: "2",
    Services:map[string]dockerComposeFileService{
      devOpName:dockerComposeFileService{
        Image:"alpine:3.3",
        Entrypoint:"echo 'hello world'",
      },
    },
  }

  var dockerComposeFileBytes []byte
  dockerComposeFileBytes, err = this.yml.toYaml(&dockerComposeFile)

  err = this.fs.saveDevOpDockerComposeFile(
    pathToDevOpDir,
    dockerComposeFileBytes,
  )

  return

}
