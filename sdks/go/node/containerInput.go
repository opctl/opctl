package node

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/opctl/opctl/sdks/go/model"
)

var (
	ErrContainerInputNotFound = errors.New("interactive container not found")
	ErrContainerInputAttached = errors.New("interactive container already attached")
)

// MaxContainerInputMessageSize is the maximum stdin WebSocket payload in bytes.
const MaxContainerInputMessageSize = 64 * 1024

// ContainerInput is the input side of a running container's terminal.
type ContainerInput interface {
	io.WriteCloser
	Resize(rows, cols uint) error
	// Done closes when the attachment closes or the container exits.
	Done() <-chan struct{}
}

type containerInputs struct {
	mu       sync.Mutex
	sessions map[string]*containerInputSession
}

type containerInputSession struct {
	reader    *io.PipeReader
	writer    *io.PipeWriter
	resize    chan model.TermSize
	done      chan struct{}
	attached  bool // guarded by the registry mutex
	closeOnce sync.Once
}

func (inputs *containerInputs) register(containerID string) (*containerInputSession, func()) {
	reader, writer := io.Pipe()
	session := &containerInputSession{
		reader: reader, writer: writer,
		resize: make(chan model.TermSize, 1),
		done:   make(chan struct{}),
	}
	inputs.mu.Lock()
	if inputs.sessions == nil {
		inputs.sessions = make(map[string]*containerInputSession)
	}
	inputs.sessions[containerID] = session
	inputs.mu.Unlock()
	return session, sync.OnceFunc(func() {
		inputs.mu.Lock()
		delete(inputs.sessions, containerID)
		inputs.mu.Unlock()
		session.Close()
		reader.Close()
	})
}

func (c core) AttachContainer(ctx context.Context, containerID string) (ContainerInput, error) {
	input, _, err := c.ClaimContainerInput(ctx, containerID)
	return input, err
}

func (c core) ClaimContainerInput(ctx context.Context, containerID string) (ContainerInput, func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	inputs := c.containerInputs
	if inputs == nil {
		return nil, nil, ErrContainerInputNotFound
	}
	inputs.mu.Lock()
	defer inputs.mu.Unlock()
	session, ok := inputs.sessions[containerID]
	if !ok {
		return nil, nil, ErrContainerInputNotFound
	}
	if session.attached {
		return nil, nil, ErrContainerInputAttached
	}
	session.attached = true
	return session, sync.OnceFunc(func() {
		inputs.mu.Lock()
		session.attached = false
		inputs.mu.Unlock()
	}), nil
}

func (s *containerInputSession) Write(data []byte) (int, error) {
	return s.writer.Write(data)
}

func (s *containerInputSession) Resize(rows, cols uint) error {
	select {
	case <-s.done:
		return io.ErrClosedPipe
	case s.resize <- model.TermSize{Rows: rows, Cols: cols}:
		return nil
	}
}

func (s *containerInputSession) Done() <-chan struct{} { return s.done }

func (s *containerInputSession) Close() error {
	s.closeOnce.Do(func() {
		close(s.done)
		s.writer.Close()
	})
	return nil
}
