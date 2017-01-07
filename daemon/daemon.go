package daemon

import (
	"github.com/opspec-io/opctl/daemon/core"
	"github.com/opspec-io/opctl/daemon/tcp"
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
