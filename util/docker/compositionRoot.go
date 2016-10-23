package docker

type compositionRoot interface {
  EnsureEngineRunningUseCase() ensureEngineRunningUseCase
  GetEngineProtocolRelativeBaseUrlUseCase() getEngineProtocolRelativeBaseUrlUseCase
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
    getEngineProtocolRelativeBaseUrlUseCase:newGetEngineProtocolRelativeBaseUrlUseCase(),
  }

  return

}

type _compositionRoot struct {
  ensureEngineRunningUseCase ensureEngineRunningUseCase
  getEngineProtocolRelativeBaseUrlUseCase    getEngineProtocolRelativeBaseUrlUseCase
}

func (this _compositionRoot) EnsureEngineRunningUseCase(
) ensureEngineRunningUseCase {
  return this.ensureEngineRunningUseCase
}

func (this _compositionRoot) GetEngineProtocolRelativeBaseUrlUseCase(
) getEngineProtocolRelativeBaseUrlUseCase {
  return this.getEngineProtocolRelativeBaseUrlUseCase
}
