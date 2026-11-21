package node

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	"github.com/opctl/opctl/sdks/go/opspec"
)

//counterfeiter:generate -o internal/fakes/containerCaller.go . containerCaller
type containerCaller interface {
	// Executes a container call
	Call(
		ctx context.Context,
		containerCall *model.ContainerCall,
		inboundScope map[string]*model.Value,
		containerCallSpec *model.ContainerCallSpec,
		rootCallID string,
	) (
		map[string]*model.Value,
		error,
	)
}

func newContainerCaller(
	containerRuntime containerruntime.ContainerRuntime,
	pubSub pubsub.PubSub,
	stateStore stateStore,
) containerCaller {

	return _containerCaller{
		containerRuntime: containerRuntime,
		pubSub:           pubSub,
		stateStore:       stateStore,
	}

}

type _containerCaller struct {
	containerRuntime containerruntime.ContainerRuntime
	pubSub           pubsub.PubSub
	stateStore       stateStore
}

func (cc _containerCaller) Call(
	ctx context.Context,
	containerCall *model.ContainerCall,
	inboundScope map[string]*model.Value,
	containerCallSpec *model.ContainerCallSpec,
	rootCallID string,
) (
	map[string]*model.Value,
	error,
) {
	outputs := map[string]*model.Value{}
	var exitCode int64

	if containerCall.Image.Ref != nil && containerCall.Image.PullCreds == nil {
		if auth := cc.stateStore.TryGetAuth(*containerCall.Image.Ref); auth != nil {
			containerCall.Image.PullCreds = &auth.Creds
		}
	}

	stdout, stderr := cc.newLogWriters(containerCall, rootCallID)

	outputs = cc.interpretOutputs(
		containerCallSpec,
		containerCall,
	)

	rawExitCode, err := cc.containerRuntime.RunContainer(
		ctx,
		containerCall,
		rootCallID,
		cc.pubSub,
		stdout,
		stderr,
	)

	// @TODO: handle no exit code
	if rawExitCode != nil {
		exitCode = *rawExitCode
	}

	if exitCode != 0 {
		err = fmt.Errorf("nonzero container exit code: %d", exitCode)
	}

	return outputs, err
}

// newLogWriters returns stdout and stderr writers that publish each line as
// an event from within Write. They share a logOrder, so events keep the order the
// container runtime wrote them in.
func (this _containerCaller) newLogWriters(
	containerCall *model.ContainerCall,
	rootCallID string,
) (stdout io.WriteCloser, stderr io.WriteCloser) {
	order := new(logOrder)
	stdout = &lineEventWriter{order: order, publish: func(line []byte) {
		this.pubSub.Publish(model.Event{
			Timestamp: time.Now().UTC(),
			ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
				Data:        line,
				ContainerID: containerCall.ContainerID,
				OpRef:       containerCall.OpPath,
				RootCallID:  rootCallID,
			},
		})
	}}
	stderr = &lineEventWriter{order: order, publish: func(line []byte) {
		this.pubSub.Publish(model.Event{
			Timestamp: time.Now().UTC(),
			ContainerStdErrWrittenTo: &model.ContainerStdErrWrittenTo{
				Data:        line,
				ContainerID: containerCall.ContainerID,
				OpRef:       containerCall.OpPath,
				RootCallID:  rootCallID,
			},
		})
	}}
	return stdout, stderr
}

// logOrder is shared by a container's stdout and stderr writers.
type logOrder struct {
	sync.Mutex
	// partial is the writer holding an unterminated line, if any
	partial *lineEventWriter
}

// lineEventWriter publishes complete lines. A partial line is published when
// the other stream writes, or on Close, so fragments keep their order too.
type lineEventWriter struct {
	order   *logOrder
	buf     []byte
	closed  bool
	publish func(line []byte)
}

func (w *lineEventWriter) Write(p []byte) (int, error) {
	w.order.Lock()
	defer w.order.Unlock()

	if w.closed {
		// like the io.Pipe this replaced: nothing is published after Close
		return 0, io.ErrClosedPipe
	}
	if other := w.order.partial; other != nil && other != w {
		other.flush()
	}

	n := len(p)
	// scan only the new bytes, so a long line arriving in many writes stays linear
	for {
		i := bytes.IndexByte(p, '\n')
		if i < 0 {
			w.buf = append(w.buf, p...)
			break
		}
		w.publish(append(bytes.Clone(w.buf), p[:i+1]...))
		w.buf = w.buf[:0]
		p = p[i+1:]
	}
	if len(w.buf) > 0 {
		w.order.partial = w
	} else if w.order.partial == w {
		w.order.partial = nil
	}
	return n, nil
}

// flush publishes any partial line; callers hold the order lock.
func (w *lineEventWriter) flush() {
	if len(w.buf) > 0 {
		w.publish(bytes.Clone(w.buf))
		w.buf = nil
	}
	if w.order.partial == w {
		w.order.partial = nil
	}
}

func (w *lineEventWriter) Close() error {
	w.order.Lock()
	defer w.order.Unlock()

	w.flush()
	w.closed = true
	return nil
}

func (this _containerCaller) interpretOutputs(
	containerCallSpec *model.ContainerCallSpec,
	containerCall *model.ContainerCall,
) map[string]*model.Value {
	outputs := map[string]*model.Value{}

	for socketAddr, name := range containerCallSpec.Sockets {
		// add socket outputs
		if "0.0.0.0" == socketAddr {
			outputs[name] = &model.Value{Socket: &containerCall.ContainerID}
		}
	}
	for callSpecContainerFilePath, mountSrc := range containerCallSpec.Files {
		mountSrcStr, ok := mountSrc.(string)
		if !ok {
			continue
		}

		if mountSrcStr == "" {
			// skip embedded files
			continue
		}

		// add file outputs
		for callContainerFilePath, callHostFilePath := range containerCall.Files {
			if callSpecContainerFilePath == callContainerFilePath {
				// copy callHostFilePath before taking address; range vars have same address for every iteration
				value := callHostFilePath
				outputs[opspec.RefToName(mountSrcStr)] = &model.Value{File: &value}
			}
		}
	}
	for callSpecContainerDirPath, mountSrc := range containerCallSpec.Dirs {
		mountSrcStr, ok := mountSrc.(string)
		if !ok {
			continue
		}

		if mountSrcStr == "" {
			// skip embedded dirs
			continue
		}

		// add dir outputs
		for callContainerDirPath, callHostDirPath := range containerCall.Dirs {
			if callSpecContainerDirPath == callContainerDirPath {
				// copy callHostDirPath before taking address; range vars have same address for every iteration
				value := callHostDirPath
				outputs[opspec.RefToName(mountSrcStr)] = &model.Value{Dir: &value}
			}
		}
	}

	return outputs
}
