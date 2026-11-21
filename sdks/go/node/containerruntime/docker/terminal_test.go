package docker

import (
	"bufio"
	"context"
	"io"
	"net"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Describe("Docker TTY input", func() {
	It("unblocks pipe writers when the Docker input connection fails", func() {
		conn, peer := net.Pipe()
		defer conn.Close()
		Expect(peer.Close()).To(Succeed())
		reader, writer := io.Pipe()
		defer reader.Close()
		defer writer.Close()

		go copyTerminalInput(types.HijackedResponse{Conn: conn}, reader)
		_, err := io.WriteString(writer, "first")
		Expect(err).NotTo(HaveOccurred())
		result := make(chan error, 1)
		go func() {
			_, err := io.WriteString(writer, "second")
			result <- err
		}()
		Eventually(result).Should(Receive(MatchError(io.ErrClosedPipe)))
	})
})

var _ = Describe("real Docker TTY", func() {
	It("sets initial size, streams prompts, forwards input and resize, and ends a canonical read", func() {
		host := os.Getenv("OPCTL_TEST_DOCKER_HOST")
		if host == "" {
			Skip("set OPCTL_TEST_DOCKER_HOST to run against Docker with alpine:3 available")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		docker, err := dockerclient.NewClientWithOpts(dockerclient.WithHost(host), dockerclient.WithAPIVersionNegotiation())
		Expect(err).NotTo(HaveOccurred())
		defer docker.Close()
		config := &container.Config{
			Image: "alpine:3", Tty: true, OpenStdin: true, AttachStdin: true, StdinOnce: true,
			Cmd: []string{"sh", "-c", `
				stty size
				printf 'Name? '
				IFS= read -r x
				printf 'got=%s rc=%s\n' "$x" "$?"
				for want in '31 91' '42 102'; do
					until [ "$(stty size 2>/dev/null)" = "$want" ]; do sleep 0.02; done
					stty size
				done
				printf 'Again? '
				IFS= read -r x
				printf 'got=%s rc=%s\n' "$x" "$?"
			`},
		}
		hostConfig := constructHostConfig(nil, nil, nil, nil, false, model.TermSize{Rows: 25, Cols: 75})
		hostConfig.Privileged = false
		hostConfig.NetworkMode = "none"
		created, err := docker.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(docker.ContainerRemove(ctx, created.ID, container.RemoveOptions{Force: true})).To(Succeed())
		}()

		runtime := _runContainer{dockerClient: docker}
		tty, err := runtime.attachTerminal(ctx, created.ID)
		Expect(err).NotTo(HaveOccurred())
		defer tty.close()
		Expect(docker.ContainerStart(ctx, created.ID, container.StartOptions{})).To(Succeed())

		reader, writer := io.Pipe()
		defer reader.Close()
		defer writer.Close()
		go copyTerminalInput(tty.input, reader)
		resize := make(chan model.TermSize, 1)
		go runtime.resizeTerminal(ctx, created.ID, resize)
		output := bufio.NewReader(tty.output.Reader)

		expectPrompt := func(want string) {
			got := make([]byte, len(want))
			_, err := io.ReadFull(output, got)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(got)).To(Equal(want))
		}
		expectLine := func(want string) {
			got, err := output.ReadString('\n')
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(want + "\r\n"))
		}

		expectLine("25 75")
		expectPrompt("Name? ")
		_, err = io.WriteString(writer, "hello\r")
		Expect(err).NotTo(HaveOccurred())
		expectLine("hello")
		expectLine("got=hello rc=0")
		resize <- model.TermSize{Rows: 31, Cols: 91}
		expectLine("31 91")
		resize <- model.TermSize{Rows: 42, Cols: 102}
		expectLine("42 102")
		expectPrompt("Again? ")
		Expect(writer.Close()).To(Succeed())
		expectLine("got= rc=1")
		wait, waitErr := docker.ContainerWait(ctx, created.ID, container.WaitConditionNotRunning)
		select {
		case result := <-wait:
			Expect(result.StatusCode).To(BeZero())
		case err := <-waitErr:
			Expect(err).NotTo(HaveOccurred())
			Fail("the Docker wait channel closed without an exit status")
		case <-ctx.Done():
			Fail("the container did not exit after EOF")
		}
	})
})
