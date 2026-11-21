package attach

import (
	"io"
)

// NewlineWriter supplies the carriage returns normally added by the terminal driver.
type NewlineWriter struct {
	io.Writer
	IsRaw  func() bool
	lastCR bool
}

func (w *NewlineWriter) Write(data []byte) (int, error) {
	if !w.IsRaw() {
		w.lastCR = false
		return w.Writer.Write(data)
	}

	converted := make([]byte, 0, len(data))
	lastCR := w.lastCR
	for _, b := range data {
		if b == '\n' && !lastCR {
			converted = append(converted, '\r')
		}
		converted = append(converted, b)
		lastCR = b == '\r'
	}
	w.lastCR = lastCR
	_, err := w.Writer.Write(converted)
	return len(data), err
}
