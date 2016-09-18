package core

//go:generate counterfeiter -o fakeContainerEngine.go --fake-name FakeContainerEngine ./ ContainerEngine

type ContainerEngine interface {
  StartContainer(
  opRunArgs map[string]string,
  opBundlePath string,
  opName string,
  opRunId string,
  eventPublisher EventPublisher,
  rootOpRunId string,
  ) (err error)

  EnsureContainerRemoved(
  opBundlePath string,
  opRunId string,
  )
}
