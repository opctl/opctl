package dockercompose

type initDevOpUcExecuter interface {
  Execute(
  devOpName string,
  ) (err error)
}

func newInitDevOpUcExecuter(
fs filesystem,
yml yamlCodec,
) initDevOpUcExecuter {

  return &initDevOpUcExecuterImpl{
    fs:fs,
    yml:yml,
  }

}

type initDevOpUcExecuterImpl struct {
  fs  filesystem
  yml yamlCodec
}

func (uc initDevOpUcExecuterImpl) Execute(
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
  dockerComposeFileBytes, err= uc.yml.toYaml(&dockerComposeFile)

  return uc.
  fs.
  saveDevOpDockerComposeFile(
    devOpName,
    dockerComposeFileBytes,
  )

}
