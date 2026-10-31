// Package goop prepares and validates the container bootstrap for Go ops.
package goop

import (
	"encoding/json"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// SchemaRange is the inclusive range of opspec schemas the node can read.
type SchemaRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Request is the language-neutral bootstrap request.
type Request struct {
	ContractVersion int                    `json:"contractVersion"`
	Phase           string                 `json:"phase"`
	OpspecSchema    SchemaRange            `json:"opspecSchema"`
	Inputs          map[string]interface{} `json:"inputs"`
	KnownOutputs    map[string]interface{} `json:"knownOutputs"`
}

// Response contains a Go op's declaration as OpSpec JSON.
type Response struct {
	ContractVersion int             `json:"contractVersion"`
	OpspecSchema    int             `json:"opspecSchema"`
	Opspec          json.RawMessage `json:"opspec"`
}

// NewRequest requests a bootstrap phase with public scalar or structured inputs.
func NewRequest(phase string, inputs map[string]interface{}) Request {
	if inputs == nil {
		inputs = map[string]interface{}{}
	}
	return Request{ContractVersion: 1, Phase: phase, OpspecSchema: SchemaRange{Min: 1, Max: opfile.SchemaVersion}, Inputs: inputs, KnownOutputs: map[string]interface{}{}}
}

func parseResponse(data []byte) (*model.OpSpec, error) {
	var response Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("decode Go op bootstrap response: %w", err)
	}
	if response.ContractVersion != 1 {
		return nil, fmt.Errorf("unsupported Go op bootstrap contract version %d", response.ContractVersion)
	}
	if response.OpspecSchema < 1 || response.OpspecSchema > opfile.SchemaVersion {
		return nil, fmt.Errorf("Go op requires opspec schema %d; upgrade opctl to a version supporting it", response.OpspecSchema)
	}
	spec, err := opfile.Unmarshal(response.Opspec)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// ParseDescribe validates an interface without a run graph or outputs.
func ParseDescribe(data []byte) (*model.OpSpec, error) {
	spec, err := parseResponse(data)
	if err != nil {
		return nil, err
	}
	if spec.Run != nil || len(spec.Outputs) != 0 {
		return nil, fmt.Errorf("Go op describe must not include run or outputs")
	}
	return spec, nil
}

// ParsePlan validates a run graph without input or output declarations.
func ParsePlan(data []byte) (*model.OpSpec, error) {
	spec, err := parseResponse(data)
	if err != nil {
		return nil, err
	}
	if len(spec.Inputs) != 0 || len(spec.Outputs) != 0 {
		return nil, fmt.Errorf("Go op plan must not include inputs or outputs")
	}
	if spec.Run == nil {
		return nil, fmt.Errorf("Go op plan must include a run call")
	}
	return spec, nil
}

// PlanInputs keeps secrets and filesystem values out of bootstrap requests.
func PlanInputs(params map[string]*model.ParamSpec, values map[string]*model.Value) (map[string]interface{}, error) {
	inputs := map[string]interface{}{}
	for name, param := range params {
		if param == nil {
			return nil, fmt.Errorf("input %q has no declaration", name)
		}
		if param.Dir != nil || param.File != nil || param.Socket != nil || (param.String != nil && param.String.IsSecret) || (param.Number != nil && param.Number.IsSecret) || (param.Array != nil && param.Array.IsSecret) || (param.Object != nil && param.Object.IsSecret) {
			continue
		}
		value := values[name]
		if value == nil {
			return nil, fmt.Errorf("missing resolved input %q", name)
		}
		unboxed, err := value.Unbox()
		if err != nil {
			return nil, fmt.Errorf("unable to unbox input %q", name)
		}
		inputs[name] = unboxed
	}
	return inputs, nil
}
