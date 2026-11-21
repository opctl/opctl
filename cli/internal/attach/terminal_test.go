package attach

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/fakes"
	"golang.org/x/sys/unix"
)

type recordingInput struct {
	bytes.Buffer
	done chan struct{}
}

func (*recordingInput) Close() error            { return nil }
func (*recordingInput) Resize(uint, uint) error { return nil }
func (i *recordingInput) Done() <-chan struct{} { return i.done }

func TestTerminalFocus(t *testing.T) {
	for _, scenario := range []string{"switches back from cooked mode", "discards input on automatic focus changes", "removes disconnected attachments"} {
		t.Run(scenario, func(t *testing.T) {
			master, slave, err := pty.Open()
			if err != nil {
				t.Fatal(err)
			}
			defer master.Close()
			defer slave.Close()
			fd := int(slave.Fd())
			original, err := unix.IoctlGetTermios(fd, getTermios)
			if err != nil {
				t.Fatal(err)
			}
			core := new(fakes.FakeNode)
			first := &recordingInput{done: make(chan struct{})}
			second := &recordingInput{done: make(chan struct{})}
			core.AttachContainerReturnsOnCall(0, first, nil)
			core.AttachContainerReturnsOnCall(1, second, nil)
			terminal := NewTerminal(context.Background(), core, slave, slave, slave)
			defer terminal.Close()
			var notice bytes.Buffer
			terminal.notice = &notice

			start := func(id string) {
				t.Helper()
				if err := terminal.Started(&model.ContainerStarted{ContainerID: id}); err != nil {
					t.Fatal(err)
				}
			}
			typeKeys := func(keys string) {
				t.Helper()
				if _, err := master.WriteString(keys); err != nil {
					t.Fatal(err)
				}
			}
			pollUntil := func(condition func() bool) {
				t.Helper()
				deadline := time.Now().Add(time.Second)
				for !condition() && time.Now().Before(deadline) {
					if err := terminal.PollInput(nil); err != nil {
						t.Fatal(err)
					}
					time.Sleep(time.Millisecond)
				}
				if !condition() {
					t.Fatalf("input %q / %q, focus notices %q", first.String(), second.String(), notice.String())
				}
			}

			start("first")
			switch scenario {
			case "switches back from cooked mode":
				typeKeys("\x1d")
				pollUntil(func() bool {
					return bytes.Contains(notice.Bytes(), []byte("input → none (Ctrl-] to focus a container; Ctrl-C stops the op)"))
				})
				const cookedFlags = unix.ICANON | unix.ISIG
				settings, err := unix.IoctlGetTermios(fd, getTermios)
				if err != nil {
					t.Fatal(err)
				}
				if terminal.Focused() || settings.Lflag&cookedFlags != cookedFlags {
					t.Fatal("unfocused input must be canonical and handle Ctrl-C as a signal")
				}
				typeKeys("ignored\x1d")
				pollUntil(func() bool { return bytes.Count(notice.Bytes(), []byte("input → first (Ctrl-] to switch)")) == 2 })
				typeKeys("first")
				pollUntil(func() bool { return first.String() == "first" })
			case "discards input on automatic focus changes":
				typeKeys("old input")
				start("second")
				typeKeys("second")
				pollUntil(func() bool { return second.String() == "second" })
				typeKeys("more old input")
				if err := terminal.Ended("second"); err != nil {
					t.Fatal(err)
				}
				typeKeys("first")
				pollUntil(func() bool { return first.String() == "first" })
			case "removes disconnected attachments":
				start("second")
				typeKeys("old input")
				close(second.done)
				pollUntil(func() bool { return bytes.Count(notice.Bytes(), []byte("input → first")) == 2 })
				typeKeys("first")
				pollUntil(func() bool { return first.String() == "first" })
				if second.Len() != 0 {
					t.Fatalf("disconnected attachment got %q", second.String())
				}
			}
			if err := terminal.Close(); err != nil {
				t.Fatal(err)
			}
			restored, err := unix.IoctlGetTermios(fd, getTermios)
			if err != nil {
				t.Fatal(err)
			}
			if *restored != *original {
				t.Fatal("closing the terminal must restore its original settings")
			}
		})
	}
}
