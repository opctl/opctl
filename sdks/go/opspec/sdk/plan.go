package sdk

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type response struct {
	ContractVersion int    `json:"contractVersion"`
	OpspecSchema    int    `json:"opspecSchema"`
	Opspec          opSpec `json:"opspec"`
}

type opSpec struct {
	Name string    `json:"name"`
	Run  *callSpec `json:"run,omitempty"`
}

type callSpec struct {
	Container *containerSpec `json:"container"`
}

// Plan builds a Go op's op graph and returns its bootstrap response as JSON.
func Plan(fn func() Call) ([]byte, error) { return emit("plan", fn) }

func emit(phase string, fn func() Call) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	spec := opSpec{Name: filepath.Base(cwd)}
	if phase == "plan" {
		if fn == nil {
			return nil, errors.New("Go op plan needs a function returning a container call")
		}
		call := fn()
		if call.container.Image.Ref == "" {
			return nil, errors.New("Go op plan must return a container call with an image")
		}
		spec.Run = &callSpec{Container: &call.container}
	}
	data, err := json.Marshal(response{ContractVersion: 1, OpspecSchema: 1, Opspec: spec})
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}
