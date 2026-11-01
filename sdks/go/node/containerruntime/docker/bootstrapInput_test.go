package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
)

type halfCloseConnection struct {
	net.Conn
	done chan struct{}
}

func (c halfCloseConnection) CloseWrite() error { close(c.done); return nil }
func TestPrivateBootstrapStdin(t *testing.T) {
	marker := strings.Repeat("synthetic-private-input", 4)
	client := new(fakes.FakeCommonAPIClient)
	done := make(chan struct{})
	reader, peer := net.Pipe()
	defer peer.Close()
	client.ContainerAttachReturns(types.NewHijackedResponse(halfCloseConnection{reader, done}, ""), nil)
	wait := make(chan container.WaitResponse)
	close(wait)
	client.ContainerWaitReturns(wait, nil)
	inputReceived := make(chan bool, 1)
	client.ContainerStartStub = func(context.Context, string, container.StartOptions) error {
		go func() {
			defer peer.Close()
			input := make([]byte, len(marker))
			_, err := io.ReadFull(peer, input)
			inputReceived <- err == nil && string(input) == marker
			stdcopy.NewStdWriter(peer, stdcopy.Stdout).Write([]byte("public result"))
		}()
		return nil
	}
	image := "sha256:prepared"
	req := &model.ContainerCall{Image: &model.ContainerCallImage{Ref: &image, Prepared: true}}
	_, err := (_runContainer{dockerClient: client}).RunContainer(containerruntime.WithBootstrapInput(context.Background(), strings.NewReader(marker)), req, "root", discardEvents{}, nopWriteCloser{io.Discard}, nopWriteCloser{io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if !<-inputReceived {
		t.Fatal("private bytes did not arrive on stdin")
	}
	select {
	case <-done:
	default:
		t.Fatal("stdin was not half-closed")
	}
	_, config, host, _, _, _ := client.ContainerCreateArgsForCall(0)
	serialized, _ := json.Marshal([]interface{}{req, config, host})
	if bytes.Contains(serialized, []byte(marker)) {
		t.Fatal("private input entered serialized create/call config")
	}
	if !config.OpenStdin || !config.AttachStdin || !config.StdinOnce || host.Tmpfs["/run/opctl-auth"] == "" {
		t.Fatal("missing private stdin/tmpfs configuration")
	}
	_, _, options := client.ContainerAttachArgsForCall(0)
	if !options.Stdin {
		t.Fatal("stdin was not attached")
	}
	if client.ContainerLogsCallCount() != 0 || host.LogConfig.Type != "none" {
		t.Fatal("private bootstrap enabled persistent logs")
	}
}

func TestPrivateInputLimit(t *testing.T) {
	for _, size := range []int{containerruntime.BootstrapInputLimit - 1, containerruntime.BootstrapInputLimit, containerruntime.BootstrapInputLimit + 1} {
		var delivered bytes.Buffer
		err := streamBootstrapInput(&delivered, strings.NewReader(strings.Repeat("x", size)))
		if (err != nil) != (size > containerruntime.BootstrapInputLimit) {
			t.Fatalf("size %d: %v", size, err)
		}
		if delivered.Len() > containerruntime.BootstrapInputLimit {
			t.Fatal("transport exceeded archive bound")
		}
	}
}
