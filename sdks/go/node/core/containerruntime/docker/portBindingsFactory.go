package docker

//go:generate counterfeiter -o ./fakePortBindingsFactory.go --fake-name fakePortBindingsFactory ./ portBindingsFactory

import (
	"fmt"
	"github.com/docker/go-connections/nat"
)

type portBindingsFactory interface {
	Construct(
		containerCallPorts map[string]string,
	) (
		nat.PortMap,
		error,
	)
}

func newPortBindingsFactory() portBindingsFactory {
	return _portBindingsFactory{}
}

type _portBindingsFactory struct{}

func (pbf _portBindingsFactory) Construct(
	containerCallPorts map[string]string,
) (
	nat.PortMap,
	error,
) {
	portBindings := nat.PortMap{}
	for containerPort, hostPort := range containerCallPorts {
		portMappings, err := nat.ParsePortSpec(fmt.Sprintf("%v:%v", hostPort, containerPort))
		if nil != err {
			return nil, err
		}
		for _, portMapping := range portMappings {
			if _, ok := portBindings[portMapping.Port]; ok {
				portBindings[portMapping.Port] = append(portBindings[portMapping.Port], portMapping.Binding)
			} else {
				portBindings[portMapping.Port] = []nat.PortBinding{portMapping.Binding}
			}
		}
	}

	return portBindings, nil
}
