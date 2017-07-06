package core

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// bindLines binds names to lines read from an io.Reader.
// bindings consist of linePrefix => name
// onBind(name, value) will be called for each line found starting w/ a bound linePrefix.
// values will be trimmed of bound linePrefix & EOL chars
func bindLines(
	stream io.Reader,
	bindings map[string]string,
	onBind func(name string, value *string),
) error {
	reader := bufio.NewReader(stream)

	var err error
	for {
		var buffer bytes.Buffer
		var l []byte
		var isPrefix bool

		for {
			// use ReadString NOT Scanner to support long lines
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			// If we've reached the end of the line, stop reading.
			if !isPrefix {
				break
			}

			// If we're just at the EOF, break
			if err != nil {
				break
			}
		}

		line := buffer.String()
		for boundPrefix, name := range bindings {
			trimmedLine := strings.TrimPrefix(line, boundPrefix)
			if trimmedLine != line {
				// if output trimming had effect we've got a match
				onBind(name, &trimmedLine)
			}
		}

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		return err
	}
	return nil
}

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
