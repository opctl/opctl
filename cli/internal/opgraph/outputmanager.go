package opgraph

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// OutputManager allows printing a "resettable" thing at the bottom of a stream
// of terminal output, when a tty is used
type OutputManager struct {
	out        io.Writer
	getWidth   func() (int, error)
	lastHeight int
}

// NewOutputManager returns a new OutputManager
func NewOutputManager() OutputManager {
	return OutputManager{
		getWidth: func() (int, error) {
			w, _, err := term.GetSize(int(os.Stdout.Fd()))
			return w, err
		},
		out: os.Stdout,
	}
}

// Clear clears the last thing printed by this object
func (o *OutputManager) Clear() error {
	w, err := o.getWidth()
	if err != nil {
		return err
	}
	// cursor to start of line (real big number) + clear line
	io.WriteString(o.out, fmt.Sprintf("\033[%dD\033[K", w))
	for i := 1; i < o.lastHeight; i++ {
		// move up one line + clear line
		io.WriteString(o.out, "\033[1A\033[K")
	}
	return nil
}

// Print prints the given string, with a preceding separator and width limited
// to the size of the terminal
func (o *OutputManager) Print(str string) error {
	w, err := o.getWidth()
	if err != nil {
		return err
	}
	lines := strings.Split(str, "\n")

	ruleWidth := 0
	for _, line := range lines {
		visualLen := countChars(stripAnsi(line))
		if visualLen > ruleWidth {
			ruleWidth = visualLen
		}
	}
	if ruleWidth > w {
		ruleWidth = w
	}

	io.WriteString(o.out, fmt.Sprintln(strings.Repeat("┄", ruleWidth)))

	for i, line := range lines {
		withoutAnsi := stripAnsi(line)
		if countChars(withoutAnsi) > w {
			// - append an ellipsis to make it more obvious the line's being truncated
			// - remove _two_ chars, not just one for the ellipsis, because the cursor
			//   will take up an additional space and cause output issues
			// - append a "reset" code to prevent color wrapping to next line
			io.WriteString(o.out, stripAnsiToLength(line, w-2)+"…\033[0m")
		} else {
			io.WriteString(o.out, line)
		}
		if i < len(lines)-1 {
			io.WriteString(o.out, "\n")
		}
	}

	o.lastHeight = len(lines) + 1
	return nil
}
