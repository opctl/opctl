package node

import (
	"github.com/opspec-io/opctl/pkg/node/core"
	"github.com/opspec-io/opctl/pkg/node/tcp"
	"github.com/opspec-io/opctl/util/containerprovider/docker"
)

func New() {
	containerProvider, err := docker.New()
	if nil != err {
		panic(err)
	}

	tcp.New(
		core.New(
			containerProvider,
		),
	).Start()

}
