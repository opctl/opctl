package engineprovider

//go:generate counterfeiter -o providers/fake/engineProvider.go --fake-name EngineProvider ./ EngineProvider

type EngineProvider interface {
  EnsureEngineRunning(
  ) (err error)

  GetEngineProtocolRelativeBaseUrl(
  ) (
  protocolRelativeBaseUrl string,
  err error,
  )
}
