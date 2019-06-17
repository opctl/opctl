package iio

import (
	"io"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ IIO

// native io package exposed via interface
type IIO interface {
	// Pipe creates a synchronous in-memory pipe.
	// It can be used to connect code expecting an io.Reader
	// with code expecting an io.Writer.
	//
	// Reads and Writes on the pipe are matched one to one
	// except when multiple Reads are needed to consume a single Write.
	// That is, each Write to the PipeWriter blocks until it has satisfied
	// one or more Reads from the PipeReader that fully consume
	// the written data.
	// The data is copied directly from the Write to the corresponding
	// Read (or Reads); there is no internal buffering.
	//
	// It is safe to call Read and Write in parallel with each other or with Close.
	// Parallel calls to Read and parallel calls to Write are also safe:
	// the individual calls will be gated sequentially.
	Pipe() (*io.PipeReader, *io.PipeWriter)

	// Copy copies from src to dst until either EOF is reached
	// on src or an error occurs. It returns the number of bytes
	// copied and the first error encountered while copying, if any.
	//
	// A successful Copy returns err == nil, not err == EOF.
	// Because Copy is defined to read from src until EOF, it does
	// not treat an EOF from Read as an error to be reported.
	//
	// If src implements the WriterTo interface,
	// the copy is implemented by calling src.WriteTo(dst).
	// Otherwise, if dst implements the ReaderFrom interface,
	// the copy is implemented by calling dst.ReadFrom(src).
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
}

func New() IIO {
	return _IIO{}
}

type _IIO struct{}

func (iio _IIO) Pipe() (*io.PipeReader, *io.PipeWriter) {
	return io.Pipe()
}

func (iio _IIO) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}
