package docker

type compositionRoot interface {
  EnsureEngineRunningUseCase() ensureEngineRunningUseCase
  GetEngineBaseUrlUseCase() getEngineBaseUrlUseCase
}

func newCompositionRoot(
) (compositionRoot compositionRoot) {

  ensureEngineRunningUseCase := newEnsureEngineRunningUseCase(
    newContainerRemover(),
    newContainerStarter(newPathNormalizer()),
    newContainerChecker(),
  )

  compositionRoot = &_compositionRoot{
    ensureEngineRunningUseCase:ensureEngineRunningUseCase,
    getEngineBaseUrlUseCase:newGetEngineBaseUrlUseCase(),
  }

  return

}

type _compositionRoot struct {
  ensureEngineRunningUseCase ensureEngineRunningUseCase
  getEngineBaseUrlUseCase    getEngineBaseUrlUseCase
}

func (this _compositionRoot) EnsureEngineRunningUseCase(
) ensureEngineRunningUseCase {
  return this.ensureEngineRunningUseCase
}

func (this _compositionRoot) GetEngineBaseUrlUseCase(
) getEngineBaseUrlUseCase {
  return this.getEngineBaseUrlUseCase
}
