package docker

type ensureEngineRunningUseCase interface {
  Execute(
  image string,
  ) (err error)
}

func newEnsureEngineRunningUseCase(
containerRemover           containerRemover,
containerStarter           containerStarter,
containerChecker  containerChecker,
) (ensureEngineRunningUseCase ensureEngineRunningUseCase) {

  ensureEngineRunningUseCase = &_ensureEngineRunningUseCase{
    containerRemover:containerRemover,
    containerStarter:containerStarter,
    containerChecker:containerChecker,
  }

  return

}

type _ensureEngineRunningUseCase struct {
  containerRemover containerRemover
  containerStarter containerStarter
  containerChecker containerChecker
}

func (this _ensureEngineRunningUseCase) Execute(
image string,
) (err error) {

  // handle obsolete container
  this.containerRemover.RemoveIfExists(obsoleteContainerName)

  // if valid container running or error checking, return
  isValidContainerRunning, err := this.containerChecker.IsValidContainerRunning(image)
  if (nil != err || isValidContainerRunning) {
    return
  }

  // cleanup invalid container
  this.containerRemover.RemoveIfExists(containerName)

  // start fresh container
  err = this.containerStarter.Start(image)

  return
}
