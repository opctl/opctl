package nodeprovider

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/nodeProvider.go . NodeProvider
type NodeProvider interface {
	ListNodes() (nodes []NodeHandle, err error)
	CreateNodeIfNotExists() (nodeHandle NodeHandle, err error)
	KillNodeIfExists(nodeId string) (err error)
}
