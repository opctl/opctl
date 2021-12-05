package docker

import (
	"fmt"

	"github.com/docker/go-connections/nat"
)

func constructPortBindings(
	containerCallPorts map[string]string,
) (
	nat.PortMap,
	error,
) {
	portBindings := nat.PortMap{}
	for containerPort, hostPort := range containerCallPorts {
		portMappings, err := nat.ParsePortSpec(fmt.Sprintf("%v:%v", hostPort, containerPort))
		if err != nil {
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
