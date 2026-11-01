package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/opctl/opctl/sdks/go/model"
	boxed "github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"
)

const containerResultLimit = 1024 * 1024

type resultStream struct {
	sync.Mutex
	name   string
	cancel func()
	data   []byte
	err    error
}

func newResultStream(name string, cancel func()) *resultStream {
	return &resultStream{name: name, cancel: cancel}
}
func (b *resultStream) Write(data []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	if b.err != nil {
		return 0, b.err
	}
	if len(data) > containerResultLimit-len(b.data) {
		b.err = fmt.Errorf("container %s exceeds 1 MiB", b.name)
		b.cancel()
		return 0, b.err
	}
	b.data = append(b.data, data...)
	return len(data), nil
}
func (b *resultStream) value() (*model.Value, error) {
	b.Lock()
	defer b.Unlock()
	if b.err != nil {
		return nil, b.err
	}
	if !utf8.Valid(b.data) {
		return nil, fmt.Errorf("container %s is not valid UTF-8", b.name)
	}
	text := strings.TrimSuffix(string(b.data), "\n")
	if len(text) != len(b.data) {
		text = strings.TrimSuffix(text, "\r")
	}
	return &model.Value{String: &text}, nil
}

func readContainerOutputs(path string) (map[string]*model.Value, error) {
	f, err := openContainerOutputs(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open outputs file; write a readable regular values.json file in $OPCTL_OUTPUTS")
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		return nil, fmt.Errorf("outputs file must be a readable regular file")
	}
	if stat.Size() > containerResultLimit {
		return nil, fmt.Errorf("outputs file exceeds 1 MiB")
	}
	data, err := io.ReadAll(io.LimitReader(f, containerResultLimit+1))
	if err != nil {
		return nil, fmt.Errorf("unable to read outputs file; check file permissions")
	}
	return parseContainerOutputs(data)
}

func parseContainerOutputs(data []byte) (map[string]*model.Value, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("outputs file is empty; write one JSON object to $OPCTL_OUTPUTS")
	}
	if len(data) > containerResultLimit {
		return nil, fmt.Errorf("outputs file exceeds 1 MiB")
	}
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("outputs file is not valid UTF-8")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := readOutputJSON(decoder, 0)
	if err != nil {
		return nil, err
	}
	object, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("outputs file must contain one JSON object")
	}
	if _, err = decoder.Token(); err != io.EOF {
		return nil, fmt.Errorf("outputs file has trailing data at offset %d", decoder.InputOffset())
	}
	results := make(map[string]*model.Value, len(object))
	for name, value := range object {
		results[name], err = boxed.Construct(value)
		if err != nil {
			return nil, fmt.Errorf("unsupported outputs file value")
		}
	}
	return results, nil
}

// Token decoding detects duplicate keys before a map can discard them.
func readOutputJSON(d *json.Decoder, depth int) (interface{}, error) {
	token, err := d.Token()
	if err != nil {
		return nil, fmt.Errorf("invalid outputs file JSON at offset %d", d.InputOffset())
	}
	switch t := token.(type) {
	case nil:
		return nil, fmt.Errorf("outputs file contains null at offset %d", d.InputOffset())
	case json.Number:
		number, err := strconv.ParseFloat(string(t), 64)
		if err != nil || math.IsInf(number, 0) || math.IsNaN(number) {
			return nil, fmt.Errorf("outputs file number exceeds float64 range at offset %d", d.InputOffset())
		}
		return number, nil
	case json.Delim:
		if depth >= 100 {
			return nil, fmt.Errorf("outputs file nesting exceeds 100 containers")
		}
		switch t {
		case '{':
			object := map[string]interface{}{}
			for d.More() {
				key, err := d.Token()
				if err != nil {
					return nil, fmt.Errorf("invalid outputs file object at offset %d", d.InputOffset())
				}
				name, ok := key.(string)
				if !ok {
					return nil, fmt.Errorf("invalid outputs file object key")
				}
				if _, exists := object[name]; exists {
					return nil, fmt.Errorf("duplicate outputs file object key at offset %d", d.InputOffset())
				}
				value, err := readOutputJSON(d, depth+1)
				if err != nil {
					return nil, err
				}
				object[name] = value
			}
			if end, err := d.Token(); err != nil || end != json.Delim('}') {
				return nil, fmt.Errorf("invalid outputs file object at offset %d", d.InputOffset())
			}
			return object, nil
		case '[':
			array := []interface{}{}
			for d.More() {
				value, err := readOutputJSON(d, depth+1)
				if err != nil {
					return nil, err
				}
				array = append(array, value)
			}
			if end, err := d.Token(); err != nil || end != json.Delim(']') {
				return nil, fmt.Errorf("invalid outputs file array at offset %d", d.InputOffset())
			}
			return array, nil
		}
		return nil, fmt.Errorf("invalid outputs file delimiter at offset %d", d.InputOffset())
	default:
		return token, nil
	}
}
