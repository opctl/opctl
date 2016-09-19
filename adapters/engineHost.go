package adapters

//go:generate counterfeiter -o fakeEngineHost.go --fake-name FakeEngineHost ./ EngineHost

type EngineHost interface {
  EnsureEngineRunning(
  image string,
  ) (err error)

  GetEngineBaseUrl(
  ) (
  baseUrl string,
  err error,
  )
}
