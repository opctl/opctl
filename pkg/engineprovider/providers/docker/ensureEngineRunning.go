package docker

func (this _engineProvider) EnsureEngineRunning(
) (err error) {

  // handle obsolete container
  this.containerRemover.RemoveIfExists(obsoleteContainerName)

  // if valid container running or error checking, return
  isValidContainerRunning, err := this.containerChecker.IsValidContainerRunning(imageRef)
  if (nil != err || isValidContainerRunning) {
    return
  }

  // cleanup invalid container
  this.containerRemover.RemoveIfExists(containerName)

  // start fresh container
  err = this.containerStarter.Start(imageRef)

  return
}
