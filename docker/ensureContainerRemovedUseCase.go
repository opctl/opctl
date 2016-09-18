package docker

//go:generate counterfeiter -o ./fakeEnsureContainerRemovedUseCase.go --fake-name fakeEnsureContainerRemovedUseCase ./ ensureContainerRemovedUseCase

type ensureContainerRemovedUseCase interface {
  Execute(
  opBundlePath string,
  opRunId string,
  )
}

func newEnsureContainerRemovedUseCase(
containerRemover containerRemover,
filesys filesys,
) ensureContainerRemovedUseCase {

  return &_ensureContainerRemovedUseCase{
    containerRemover: containerRemover,
    filesys:filesys,
  }

}

type _ensureContainerRemovedUseCase struct {
  containerRemover containerRemover
  filesys          filesys
}

func (this _ensureContainerRemovedUseCase) Execute(
opBundlePath string,
opRunId string,
) {

  if (this.filesys.isDockerComposeFileExistent(opBundlePath)) {

    this.containerRemover.EnsureContainerRemoved(
      opBundlePath,
      opRunId,
    )

  }

  return

}
