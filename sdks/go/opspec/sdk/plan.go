package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

type response struct {
	ContractVersion int               `json:"contractVersion"`
	OpspecSchema    int               `json:"opspecSchema"`
	Opspec          opSpec            `json:"opspec"`
	OutputBindings  map[string]string `json:"outputBindings,omitempty"`
}

type opSpec struct {
	Name    string               `json:"name"`
	Inputs  map[string]paramSpec `json:"inputs,omitempty"`
	Outputs map[string]paramSpec `json:"outputs,omitempty"`
	Run     *callSpec            `json:"run,omitempty"`
}

type callSpec struct {
	Op        *opCallSpec    `json:"op,omitempty"`
	Container *containerSpec `json:"container,omitempty"`
	Serial    []*callSpec    `json:"serial,omitempty"`
	Parallel  []*callSpec    `json:"parallel,omitempty"`
}

// Plan builds a Go op's op graph and returns its bootstrap response as JSON.
func Plan[I any](fn func(I) Call, inputs map[string]any) ([]byte, error) {
	return emit("plan", withoutOutputs(fn), inputs)
}

// Describe returns the op interface without calling the plan function.
func Describe[I any]() ([]byte, error) { return emit[I, struct{}]("describe", nil, nil) }

func emit[I, O any](phase string, fn func(I) (Call, O), values map[string]any) (data []byte, err error) {
	planMu.Lock()
	defer planMu.Unlock()
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
	outputs, outputFields, err := parameterSpecs[O]("output")
	if err != nil {
		return nil, err
	}
	bindings := map[string]string{}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	spec := opSpec{Name: filepath.Base(cwd)}
	if phase == "describe" {
		spec.Inputs = inputs
		spec.Outputs = outputs
	}
	if phase == "plan" {
		if fn == nil {
			return nil, errors.New("Go op plan needs a function returning a call")
		}
		in, err := bindInputs[I](fields, values)
		if err != nil {
			return nil, err
		}
		activePlan = &outputPlan{handles: map[string]outputSite{}}
		defer func() { activePlan = nil }()
		call, out := fn(in)
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
		if err := activePlan.validate(spec.Run); err != nil {
			return nil, err
		}
		target := reflect.ValueOf(out)
		for _, field := range outputFields {
			handle := target.Field(field.index).Interface().(interface{ bindingValue() inputValue }).bindingValue()
			if handle.err != nil {
				return nil, fmt.Errorf("output %q: %w", field.name, handle.err)
			}
			if handle.name == "" && handle.raw == "" {
				if field.hasDefault {
					continue
				}
				return nil, fmt.Errorf("output %q must be bound to a handle or have a default", field.name)
			}
			ref := handle.handle()
			name, site, generated := activePlan.output(ref)
			if generated && !handle.isRaw() && !activePlan.visible[name] {
				return nil, fmt.Errorf("output %q uses %s created at %s without a visible producer", field.name, site.kind, site.site)
			}
			bindings[field.name] = ref
		}
	}
	data, err = json.Marshal(response{ContractVersion: 1, OpspecSchema: minimumSchema(spec.Run), Opspec: spec, OutputBindings: bindings})
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// PlanWithOutputs builds an op graph and binds its declared output handles.
func PlanWithOutputs[I, O any](fn func(I) (Call, O), inputs map[string]any) ([]byte, error) {
	return emit("plan", fn, inputs)
}

// DescribeWithOutputs returns input and output declarations without building the graph.
func DescribeWithOutputs[I, O any]() ([]byte, error) { return emit[I, O]("describe", nil, nil) }
func withoutOutputs[I any](fn func(I) Call) func(I) (Call, struct{}) {
	if fn == nil {
		return nil
	}
	return func(in I) (Call, struct{}) { return fn(in), struct{}{} }
}

func minimumSchema(call *callSpec) int {
	version, _ := schemaRequirement(call, "run")
	return version
}
