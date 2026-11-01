package docker

import (
	"context"
	"errors"
	"io"

	"github.com/docker/docker/api/types/container"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
)

var _ = Describe("output separation", func() {
	It("keeps ordinary privilege settings when only separate streams are requested", func() {
		for _, separate := range []bool{false, true} {
			client := new(FakeCommonAPIClient)
			client.ContainerCreateReturns(container.CreateResponse{}, errors.New("stop after config"))
			runtime := _runContainer{dockerClient: client}
			ctx := context.Background()
			if separate {
				ctx = containerruntime.WithSeparateOutput(ctx)
			}
			image := "alpine:3"
			_, err := runtime.RunContainer(ctx, &model.ContainerCall{Image: &model.ContainerCallImage{Ref: &image}}, "root", noOpEventPublisher{}, nopWriteCloser{io.Discard}, nopWriteCloser{io.Discard})
			Expect(err).To(HaveOccurred())
			_, config, _, _, _, _ := client.ContainerCreateArgsForCall(client.ContainerCreateCallCount() - 1)
			Expect(config.Tty).To(Equal(!separate))
			_, _, host, _, _, _ := client.ContainerCreateArgsForCall(client.ContainerCreateCallCount() - 1)
			Expect(host.Privileged).To(BeTrue())
			if separate {
				Expect(host.LogConfig.Type).To(Equal("none"))
			}
		}
	})
})
