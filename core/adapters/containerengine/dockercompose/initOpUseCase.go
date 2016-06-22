package dockercompose

//go:generate counterfeiter -o ./fakeInitOpUseCase.go --fake-name fakeInitOpUseCase ./ initOpUseCase


type initOpUseCase interface {
  Execute(
  opBundlePath string,
  opDir string,
  ) (err error)
}

func newInitOpUseCase(
filesys filesys,
yamlCodec yamlCodec,
) initOpUseCase {

  return &_initOpUseCase{
    filesys:filesys,
    yamlCodec:yamlCodec,
  }

}

type _initOpUseCase struct {
  filesys   filesys
  yamlCodec yamlCodec
}

func (this _initOpUseCase) Execute(
opBundlePath string,
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
  dockerComposeFileBytes, err = this.yamlCodec.toYaml(&dockerComposeFile)

  err = this.filesys.saveOpDockerComposeFile(
    opBundlePath,
    dockerComposeFileBytes,
  )

  return

}
