package node

import (
	"bufio"
	"io"
)

// readChunks reads from an io.Reader in chunks
func readChunks(
	stream io.Reader,
	onChunk func(line []byte),
) error {
	// support lines up to 4MB
	reader := bufio.NewReaderSize(stream, 4e6)
	var err error
	var b []byte

	for {
		// chunk on newlines
		if b, err = reader.ReadBytes('\n'); len(b) > 0 {
			// always call onChunk if len(bytes) read to ensure full stream sent; even under error conditions
			onChunk(b)
		}

		if err != nil {
			break
		}
	}

	if io.EOF == err {
		return nil
	}
	return err
}

// readRawChunks publishes terminal output as soon as the reader provides it.
func readRawChunks(stream io.Reader, onChunk func([]byte)) error {
	buf := make([]byte, 32*1024)
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			// Subscribers can read the event after the next read reuses buf.
			onChunk(append([]byte(nil), buf[:n]...))
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
