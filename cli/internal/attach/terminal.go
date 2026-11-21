package attach

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

const switchKey = 0x1d

type target struct {
	label string
	input node.ContainerInput
}

// Terminal manages container focus and owns stdin's raw mode.
type Terminal struct {
	ctx      context.Context
	node     node.Node
	stdin    *os.File
	stdout   *os.File
	stderr   *os.File
	notice   io.Writer
	state    *term.State
	original *unix.Termios
	focus    focus
	targets  map[string]target
	active   string
	closed   bool
	buf      [4096]byte
}

func NewTerminal(ctx context.Context, node node.Node, stdin, stdout, stderr *os.File) *Terminal {
	t := &Terminal{
		ctx: ctx, node: node, stdin: stdin, stdout: stdout, stderr: stderr,
		targets: make(map[string]target),
	}
	t.notice = t.MessageWriter(stderr)
	return t
}

// Focused reports whether a container owns the terminal.
func (t *Terminal) Focused() bool { return t.active != "" }

func (t *Terminal) MessageWriter(file *os.File) io.Writer {
	if !term.IsTerminal(int(file.Fd())) {
		return file
	}
	return &NewlineWriter{Writer: file, IsRaw: func() bool { return t.state != nil }}
}

func (t *Terminal) Started(event *model.ContainerStarted) error {
	if t.closed {
		return nil
	}
	id := event.ContainerID
	if _, ok := t.targets[id]; ok {
		return nil
	}
	label := event.OpRef
	if event.Name != nil && *event.Name != "" {
		label = *event.Name
	}
	if label == "" {
		label = id
	}
	t.targets[id] = target{label: label}
	t.focus.Start(id)
	return t.updateFocus()
}

func (t *Terminal) Ended(id string) error {
	t.remove(id)
	return t.updateFocus()
}

func (t *Terminal) remove(id string) {
	if input := t.targets[id].input; input != nil {
		input.Close()
	}
	delete(t.targets, id)
	t.focus.End(id)
}

func (t *Terminal) updateFocus() error {
	var attachErr error
	for {
		id := t.focus.Current()
		if id == t.active {
			return attachErr
		}
		if id == "" {
			t.active = ""
			err := t.restore()
			fmt.Fprintln(t.notice, "input → none (Ctrl-] to focus a container; Ctrl-C stops the op)")
			return errors.Join(attachErr, err)
		}
		target := t.targets[id]
		if target.input == nil {
			ctx, cancel := context.WithTimeout(t.ctx, 2*time.Second)
			input, err := t.node.AttachContainer(ctx, id)
			cancel()
			if err != nil {
				t.remove(id)
				if !errors.Is(err, node.ErrContainerInputNotFound) {
					attachErr = errors.Join(attachErr, fmt.Errorf("input for %s: %w", target.label, err))
				}
				continue
			}
			target.input = input
			t.targets[id] = target
		}
		select {
		case <-target.input.Done():
			t.remove(id)
			continue
		default:
		}
		if err := t.makeRaw(); err != nil {
			return errors.Join(attachErr, err, t.Close())
		}
		fmt.Fprintf(t.notice, "input → %s (Ctrl-] to switch)\n", target.label)
		// Input still queued when the notice is printed belongs to the old
		// focus, including anything typed while the new attachment opened.
		if err := flushInput(int(t.stdin.Fd())); err != nil {
			return errors.Join(attachErr, err, t.Close())
		}
		t.active = id
		return errors.Join(attachErr, t.Resize())
	}
}

func (t *Terminal) makeRaw() error {
	if t.state != nil {
		return nil
	}
	fd := int(t.stdin.Fd())
	if t.original == nil {
		original, err := unix.IoctlGetTermios(fd, getTermios)
		if err != nil {
			return err
		}
		// In cooked mode, have Ctrl-] complete a line without Enter.
		// Canonical input and signal handling (including Ctrl-C) stay enabled.
		cooked := *original
		cooked.Lflag |= unix.ICANON | unix.ISIG
		cooked.Cc[unix.VEOL] = switchKey
		if err := unix.IoctlSetTermios(fd, setTermios, &cooked); err != nil {
			return err
		}
		t.original = original
	}
	state, err := term.MakeRaw(fd)
	t.state = state
	return err
}

func (t *Terminal) Resize() error {
	input := t.targets[t.active].input
	if input == nil || t.state == nil {
		return nil
	}
	cols, rows, err := term.GetSize(int(t.stdout.Fd()))
	if err != nil || rows == 0 || cols == 0 {
		return err
	}
	return input.Resize(uint(rows), uint(cols))
}

// PollInput removes disconnected targets before reading the keyboard.
// beforeFocus clears the live graph before an unfocused terminal becomes raw.
func (t *Terminal) PollInput(beforeFocus func()) error {
	for id, target := range t.targets {
		if target.input != nil {
			select {
			case <-target.input.Done():
				t.remove(id)
			default:
			}
		}
	}
	if t.active != t.focus.Current() {
		return t.updateFocus()
	}
	if len(t.targets) == 0 {
		return nil
	}

	fd := int(t.stdin.Fd())
	poll := [1]unix.PollFd{{Fd: int32(fd), Events: unix.POLLIN}}
	n, err := unix.Poll(poll[:], 0)
	if err == unix.EINTR || (err == nil && n == 0) {
		return nil
	}
	if err != nil || poll[0].Revents&unix.POLLIN == 0 {
		return errors.Join(err, t.Close())
	}
	n, err = unix.Read(fd, t.buf[:])
	if err == unix.EINTR || err == unix.EAGAIN {
		return nil
	}
	if err != nil {
		return errors.Join(err, t.Close())
	}
	if n == 0 {
		if t.Focused() {
			return t.Close()
		}
		return nil // Ctrl-D on an empty cooked line.
	}

	data := t.buf[:n]
	index := bytes.IndexByte(data, switchKey)
	keys := data
	if index >= 0 {
		keys = data[:index]
	}
	if t.Focused() && len(keys) > 0 {
		if _, err := t.targets[t.active].input.Write(keys); err != nil {
			return errors.Join(err, t.Ended(t.active))
		}
	}
	if index < 0 {
		return nil
	}

	if !t.Focused() && beforeFocus != nil {
		beforeFocus()
	}
	// Bytes after the switch key were read before the next focus notice.
	t.focus.Cycle()
	return t.updateFocus()
}

// WriteOutput writes focused container bytes directly to the terminal.
func (t *Terminal) WriteOutput(event *model.Event) bool {
	if t.active == "" {
		return false
	}
	if out := event.ContainerStdOutWrittenTo; out != nil && out.ContainerID == t.active {
		t.stdout.Write(out.Data)
		return true
	}
	if out := event.ContainerStdErrWrittenTo; out != nil && out.ContainerID == t.active {
		t.stderr.Write(out.Data)
		return true
	}
	return false
}

func (t *Terminal) restore() error {
	if t.state == nil {
		return nil
	}
	fd := int(t.stdin.Fd())
	if err := term.Restore(fd, t.state); err != nil {
		return err
	}
	t.state = nil
	// Discard input left in raw mode before the cooked terminal can replay it.
	return flushInput(fd)
}

func (t *Terminal) Close() error {
	t.closed = true
	err := t.restore()
	if t.original != nil {
		if restoreErr := unix.IoctlSetTermios(int(t.stdin.Fd()), setTermios, t.original); restoreErr != nil {
			err = errors.Join(err, restoreErr)
		} else {
			t.original = nil
			t.state = nil
		}
	}
	for id := range t.targets {
		t.remove(id)
	}
	t.active = ""
	return err
}
