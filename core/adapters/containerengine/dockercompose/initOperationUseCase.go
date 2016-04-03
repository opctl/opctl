package dockercompose

import "path"

type initOperationUseCase interface {
  Execute(
  pathToOperationDir string,
  ) (err error)
}

func newInitOperationUseCase(
fs filesystem,
yml yamlCodec,
) initOperationUseCase {

  return &_initOperationUseCase{
    fs:fs,
    yml:yml,
  }

}

type _initOperationUseCase struct {
  fs  filesystem
  yml yamlCodec
}

func (this _initOperationUseCase) Execute(
pathToOperationDir string,
) (err error) {

  operationName := path.Base(pathToOperationDir)

  var dockerComposeFile =
  dockerComposeFile{
    Version: "2",
    Services:map[string]dockerComposeFileService{
      operationName:dockerComposeFileService{
        Image:"alpine:3.3",
        Entrypoint:"echo 'hello world'",
      },
    },
  }

  var dockerComposeFileBytes []byte
  dockerComposeFileBytes, err = this.yml.toYaml(&dockerComposeFile)

  err = this.fs.saveOperationDockerComposeFile(
    pathToOperationDir,
    dockerComposeFileBytes,
  )

  return

}
