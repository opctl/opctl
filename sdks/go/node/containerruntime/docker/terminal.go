package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/opctl/opctl/sdks/go/model"
)

func (_containerRuntime) SupportsInteractive() bool { return true }

type terminalAttachment struct {
	input, output types.HijackedResponse
	stop          func() bool
}

func (cr _runContainer) attachTerminal(ctx context.Context, name string) (*terminalAttachment, error) {
	input, err := cr.dockerClient.ContainerAttach(ctx, name, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		// The CLI handles Ctrl-] itself. Preserve Docker's default Ctrl-P Ctrl-Q for the container.
		DetachKeys: "ctrl-],ctrl-]",
	})
	if err != nil {
		return nil, fmt.Errorf("attach container stdin: %w", err)
	}

	// Docker logs buffer partial lines. Use live TTY output so prompts
	// arrive immediately, on a separate connection that survives stdin EOF.
	output, err := cr.dockerClient.ContainerAttach(ctx, name, container.AttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		input.Close()
		return nil, fmt.Errorf("attach container output: %w", err)
	}
	tty := &terminalAttachment{input: input, output: output}
	tty.stop = context.AfterFunc(ctx, tty.closeStreams)
	return tty, nil
}

func (tty *terminalAttachment) close() {
	tty.stop()
	tty.closeStreams()
}

func (tty *terminalAttachment) closeStreams() {
	tty.input.Close()
	tty.output.Close()
}

func copyTerminalInput(tty types.HijackedResponse, stdin io.Reader) {
	defer tty.CloseWrite()
	if _, err := io.Copy(tty.Conn, stdin); err != nil {
		if pipe, ok := stdin.(interface{ CloseWithError(error) error }); ok {
			pipe.CloseWithError(err)
		}
		return
	}
	// Docker leaves the TTY master open. Send Ctrl-D twice so a pending
	// canonical line can be read before the next read gets EOF with the
	// default EOF key. Raw-mode programs receive two literal bytes and
	// may keep running until the op is killed.
	tty.Conn.Write([]byte{4, 4})
}

func (cr _runContainer) resizeTerminal(ctx context.Context, name string, resize <-chan model.TermSize) {
	for {
		select {
		case <-ctx.Done():
			return
		case size, ok := <-resize:
			if !ok {
				return
			}
			cr.dockerClient.ContainerResize(ctx, name, container.ResizeOptions{
				Height: size.Rows,
				Width:  size.Cols,
			})
		}
	}
}
