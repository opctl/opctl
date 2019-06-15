package core

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

		if nil != err {
			break
		}
	}

	if io.EOF == err {
		return nil
	}
	return err
}
