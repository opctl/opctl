package node

import (
	"github.com/opspec-io/opctl/node/core"
	"github.com/opspec-io/opctl/node/tcp"
	"github.com/opspec-io/opctl/pkg/containerengine/engines/docker"
)

func New() {
	containerEngine, err := docker.New()
	if nil != err {
		panic(err)
	}

	tcp.New(
		core.New(
			containerEngine,
		),
	).Start()

}
