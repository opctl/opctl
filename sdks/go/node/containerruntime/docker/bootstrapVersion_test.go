package docker

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
)

func TestBootstrapRequiresRecursiveReadOnlyAPI(t *testing.T) {
	for _, version := range []string{"1.43", "1.44", "1.45"} {
		t.Run(version, func(t *testing.T) {
			client := new(fakes.FakeCommonAPIClient)
			client.ClientVersionReturns(version)
			client.ContainerCreateReturns(container.CreateResponse{}, errors.New("stop after config"))
			image := "compiler"
			_, err := (_runContainer{dockerClient: client}).RunContainer(containerruntime.WithBootstrap(context.Background()), &model.ContainerCall{Image: &model.ContainerCallImage{Ref: &image, Prepared: true}}, "root", noOpEventPublisher{}, nopWriteCloser{io.Discard}, nopWriteCloser{io.Discard})
			if version == "1.43" {
				if err == nil || !strings.Contains(err.Error(), "Go bootstrap requires Docker API >= 1.44") || client.ContainerCreateCallCount() != 0 {
					t.Fatalf("older daemon reached bootstrap creation: %v", err)
				}
			} else if client.ContainerCreateCallCount() != 1 {
				t.Fatalf("supported daemon rejected: %v", err)
			}
		})
	}
}
