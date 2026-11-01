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
	Name   string               `json:"name"`
	Inputs map[string]paramSpec `json:"inputs,omitempty"`
	Run    *callSpec            `json:"run,omitempty"`
}

type callSpec struct {
	Container *containerSpec `json:"container,omitempty"`
	Serial    []*callSpec    `json:"serial,omitempty"`
	Parallel  []*callSpec    `json:"parallel,omitempty"`
}

// Plan builds a Go op's op graph and returns its bootstrap response as JSON.
func Plan[I any](fn func(I) Call, inputs map[string]any) ([]byte, error) {
	return emit("plan", fn, inputs)
}

// Describe returns the op interface without calling the plan function.
func Describe[I any]() ([]byte, error) { return emit[I]("describe", nil, nil) }

func emit[I any](phase string, fn func(I) Call, values map[string]any) (data []byte, err error) {
	defer func() {
		if failure := recover(); failure != nil {
			if input, ok := failure.(inputError); ok {
				data = nil
				err = input.error
			} else {
				panic(failure)
			}
		}
	}()
	inputs, fields, err := inputSpecs[I]()
	if err != nil {
		return nil, err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	spec := opSpec{Name: filepath.Base(cwd)}
	if phase == "describe" {
		spec.Inputs = inputs
	}
	if phase == "plan" {
		if fn == nil {
			return nil, errors.New("Go op plan needs a function returning a call")
		}
		in, err := bindInputs[I](fields, values)
		if err != nil {
			return nil, err
		}
		call := fn(in)
		if call == nil {
			return nil, errors.New("Go op plan must return a call")
		}
		spec.Run, err = call.render()
		if err != nil {
			return nil, err
		}
		if spec.Run == nil {
			return nil, errors.New("op graph has no calls")
		}
	}
	data, err = json.Marshal(response{ContractVersion: 1, OpspecSchema: 1, Opspec: spec})
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}
