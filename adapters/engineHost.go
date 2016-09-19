package adapters

//go:generate counterfeiter -o fakeEngineHost.go --fake-name FakeEngineHost ./ EngineHost

type EngineHost interface {
  EnsureEngineRunning(
  image string,
  ) (err error)

  GetEngineProtocolRelativeBaseUrl(
  ) (
  protocolRelativeBaseUrl string,
  err error,
  )
}
