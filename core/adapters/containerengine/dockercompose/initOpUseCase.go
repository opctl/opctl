package dockercompose

type initOpUseCase interface {
  Execute(
  pathToOpDir string,
  opDir string,
  ) (err error)
}

func newInitOpUseCase(
fs filesystem,
yml yamlCodec,
) initOpUseCase {

  return &_initOpUseCase{
    fs:fs,
    yml:yml,
  }

}

type _initOpUseCase struct {
  fs  filesystem
  yml yamlCodec
}

func (this _initOpUseCase) Execute(
pathToOpDir string,
opName string,
) (err error) {

  var dockerComposeFile =
  dockerComposeFile{
    Version: "2",
    Services:map[string]dockerComposeFileService{
      opName:dockerComposeFileService{
        Image:"alpine:3.3",
        Entrypoint:"echo 'hello world'",
      },
    },
  }

  var dockerComposeFileBytes []byte
  dockerComposeFileBytes, err = this.yml.toYaml(&dockerComposeFile)

  err = this.fs.saveOpDockerComposeFile(
    pathToOpDir,
    dockerComposeFileBytes,
  )

  return

}
