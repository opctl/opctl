package dockercompose

type initDevOpUseCase interface {
  Execute(
  devOpName string,
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
devOpName string,
) (err error) {

  var dockerComposeFile =
  dockerComposeFile{
    Version: "2",
    Services:map[string]dockerComposeFileService{
      devOpName:dockerComposeFileService{
        Image:"alpine:3.3",
      },
    },
  }

  var dockerComposeFileBytes []byte
  dockerComposeFileBytes, err= this.yml.toYaml(&dockerComposeFile)

  return this.
  fs.
  saveDevOpDockerComposeFile(
    devOpName,
    dockerComposeFileBytes,
  )

}
