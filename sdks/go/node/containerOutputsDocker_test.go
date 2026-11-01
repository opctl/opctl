package node

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
)

// Opt in on a Docker host as the non-root node user with socket access.
func TestContainerOutputsDockerCleanup(t *testing.T) {
	if os.Getenv("OPCTL_OUTPUTS_DOCKER_TEST") != "1" {
		t.Skip("set OPCTL_OUTPUTS_DOCKER_TEST=1 on a Docker host")
	}
	if os.Geteuid() == 0 {
		t.Fatal("this regression must run as a non-root node user")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	runtime, err := docker.New(ctx, "unix:///var/run/docker.sock")
	if err != nil {
		t.Fatal(err)
	}
	for _, user := range []string{"root", "nobody"} {
		t.Run(user, func(t *testing.T) {
			dir := t.TempDir()
			image := "alpine:3"
			call := &model.ContainerCall{ContainerID: "outputs-cleanup-" + time.Now().Format("150405.000000000"), Image: &model.ContainerCallImage{Ref: &image}, OutputsFile: filepath.Join(dir, "outputs", "values.json"), Cmd: []string{"sh", "-ce", `mkdir -p /opctl/outputs/leftover; printf private > /opctl/outputs/leftover/value; printf '{"count":3}' > "$OPCTL_OUTPUTS"`}}
			if user == "nobody" {
				call.Cmd = []string{"su", "-s", "/bin/sh", "nobody", "-c", `printf '{"count":3}' > "$OPCTL_OUTPUTS"`}
			}
			values, err := newContainerCaller(runtime, new(fakes.FakePubSub), &bootstrapSnapshotStore{}).Call(ctx, call, nil, &model.ContainerCallSpec{Outputs: map[string]string{"count": "$(count)"}}, "cleanup-root")
			if err != nil {
				t.Fatal(err)
			}
			if values["count"] == nil || values["count"].Number == nil || *values["count"].Number != 3 {
				t.Fatal("count not collected")
			}
			if _, err := os.Stat(filepath.Dir(call.OutputsFile)); !os.IsNotExist(err) {
				t.Fatalf("outputs scratch was not removed: %v", err)
			}
		})
	}
}
