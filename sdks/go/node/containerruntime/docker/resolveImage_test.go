package docker

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
)

type discardEvents struct{}

func (discardEvents) Publish(model.Event) {}

func TestResolveBootstrapImage(t *testing.T) {
	for _, offline := range []bool{false, true} {
		client := new(fakes.FakeCommonAPIClient)
		client.ImagePullReturns(io.NopCloser(strings.NewReader("")), nil)
		if offline {
			client.ImagePullReturns(nil, errors.New("registry unavailable"))
		}
		client.ImageInspectWithRawReturns(types.ImageInspect{ID: "sha256:exact", Os: "linux", Architecture: "amd64", Variant: "v1"}, nil, nil)
		ref := "golang:1.25"
		result, err := (_containerRuntime{dockerClient: client}).ResolveImage(context.Background(), &model.ContainerCallImage{Ref: &ref}, "root", discardEvents{})
		if err != nil || result.Ref != "sha256:exact" || result.ID != "sha256:exact" || result.OS != "linux" || result.Arch != "amd64" || result.Variant != "v1" {
			t.Fatalf("unexpected identity: %+v, %v", result, err)
		}
		if client.ImagePullCallCount() != 0 {
			t.Fatal("cached versioned image contacted the registry")
		}
		// An unavailable tag can still resolve through the inspectable offline image.
		client.ImageInspectWithRawReturnsOnCall(client.ImageInspectWithRawCallCount(), types.ImageInspect{}, nil, errors.New("tag unavailable"))
		_, err = (_containerRuntime{dockerClient: client}).ResolveImage(context.Background(), &model.ContainerCallImage{Ref: &ref}, "root", discardEvents{})
		if err != nil {
			t.Fatal(err)
		}
	}
	client := new(fakes.FakeCommonAPIClient)
	client.ImagePullReturns(nil, errors.New("private diagnostic"))
	client.ImageInspectWithRawReturns(types.ImageInspect{}, nil, errors.New("private diagnostic"))
	ref := "golang:missing"
	_, err := (_containerRuntime{dockerClient: client}).ResolveImage(context.Background(), &model.ContainerCallImage{Ref: &ref}, "root", discardEvents{})
	if err == nil || strings.Contains(err.Error(), "private diagnostic") {
		t.Fatal("resolution must fail without leaking registry diagnostics")
	}
}
