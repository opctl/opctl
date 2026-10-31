package node

import (
	"context"
	"fmt"
	"io"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

//counterfeiter:generate -o internal/fakes/pubSub.go ./pubsub PubSub

type stdoutCaptureKey struct{}

type stdoutCapture struct {
	writer      io.Writer
	size        int
	err         error
	diagnostics bootstrapDiagnostics
}

type bootstrapDiagnostics []byte

func (capture *bootstrapDiagnostics) Write(data []byte) (int, error) {
	n := len(data)
	if remaining := 64*1024 - len(*capture); len(data) > remaining {
		data = data[:remaining]
	}
	*capture = append(*capture, data...)
	return n, nil
}

func withStdoutCapture(ctx context.Context, writer io.Writer) context.Context {
	return context.WithValue(containerruntime.WithSeparateOutput(ctx), stdoutCaptureKey{}, &stdoutCapture{writer: writer})
}

func stdoutCaptureFrom(ctx context.Context) *stdoutCapture {
	capture, _ := ctx.Value(stdoutCaptureKey{}).(*stdoutCapture)
	return capture
}

func (capture *stdoutCapture) Write(data []byte) (int, error) {
	if capture.err != nil {
		return 0, capture.err
	}
	if len(data) > 32*1024*1024-capture.size {
		capture.err = fmt.Errorf("Go op bootstrap stdout exceeds 32 MiB")
		return 0, capture.err
	}
	n, err := capture.writer.Write(data)
	capture.size += n
	if err == nil && n != len(data) {
		err = io.ErrShortWrite
	}
	capture.err = err
	return n, err
}

// Writes go only to the capture; closing still ends the stdout log reader.
type teeCloser struct {
	io.Writer
	io.Closer
}

type bootstrapEventPublisher struct{ pubsub.EventPublisher }

func (publisher bootstrapEventPublisher) Publish(event model.Event) {
	// Image-pull progress bypasses the runtime's stdout writer.
	if event.ContainerStdOutWrittenTo != nil {
		diagnostic := model.ContainerStdErrWrittenTo(*event.ContainerStdOutWrittenTo)
		event.ContainerStdErrWrittenTo = &diagnostic
		event.ContainerStdOutWrittenTo = nil
	}
	publisher.EventPublisher.Publish(event)
}
