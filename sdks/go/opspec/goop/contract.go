// Package goop prepares and validates the container bootstrap for Go ops.
package goop

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
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
	ContractVersion int               `json:"contractVersion"`
	OpspecSchema    int               `json:"opspecSchema"`
	Opspec          json.RawMessage   `json:"opspec"`
	OutputBindings  map[string]string `json:"outputBindings,omitempty"`
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
	if needed := opspec.MinimumSchema(spec.Run); needed > response.OpspecSchema {
		return nil, fmt.Errorf("Go op graph requires opspec schema %d but response declares schema %d", needed, response.OpspecSchema)
	}
	return spec, nil
}

// ParseDescribe validates an interface without a run graph or runtime bindings.
func ParseDescribe(data []byte) (*model.OpSpec, error) {
	spec, err := parseResponse(data)
	if err != nil {
		return nil, err
	}
	var response map[string]json.RawMessage
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	if spec.Run != nil || response["outputBindings"] != nil {
		return nil, fmt.Errorf("Go op describe must not include run or output bindings")
	}
	return spec, nil
}

// ParsePlan validates a run graph without input or output declarations.
func ParsePlan(data []byte) (*Plan, error) {
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
	var response Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	for name, ref := range response.OutputBindings {
		if !regexp.MustCompile(`^[-_a-zA-Z0-9]+$`).MatchString(name) || strings.HasPrefix(name, "opctl_") {
			return nil, fmt.Errorf("invalid Go output binding name %q", name)
		}
		if !wholeBinding(ref) {
			return nil, fmt.Errorf("Go output binding %q must be a whole reference", name)
		}
	}
	return &Plan{OpSpec: spec, OutputBindings: response.OutputBindings}, nil
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

// Plan keeps runtime aliases separate from the op interface and graph.
type Plan struct {
	*model.OpSpec
	OutputBindings map[string]string
}

func wholeBinding(ref string) bool {
	if !strings.HasPrefix(ref, "$(") || !strings.HasSuffix(ref, ")") {
		return false
	}
	depth := 0
	for i := 0; i < len(ref); i++ {
		if ref[i] == '(' {
			depth++
		}
		if ref[i] == ')' {
			depth--
			if depth == 0 && i != len(ref)-1 {
				return false
			}
		}
		if depth < 0 {
			return false
		}
	}
	return depth == 0 && len(ref) > 3
}

// HasDefault reports an authored parameter default, including false and empty strings.
func HasDefault(param *model.ParamSpec) bool {
	if param == nil {
		return false
	}
	fields := reflect.ValueOf(param).Elem()
	for i := 0; i < fields.NumField(); i++ {
		body := fields.Field(i)
		if body.Kind() != reflect.Pointer || body.IsNil() {
			continue
		}
		value := body.Elem().FieldByName("Default")
		if value.IsValid() && !value.IsNil() {
			return true
		}
	}
	return false
}

// ValidateOutputs rejects incomplete or undeclared aliases before running the graph.
func (plan *Plan) ValidateOutputs(params map[string]*model.ParamSpec) error {
	names := make([]string, 0, len(plan.OutputBindings))
	for name := range plan.OutputBindings {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if _, ok := params[name]; !ok {
			return fmt.Errorf("undeclared Go output binding %q", name)
		}
	}
	names = nil
	for name := range params {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if _, ok := plan.OutputBindings[name]; !ok && !HasDefault(params[name]) {
			return fmt.Errorf("missing Go output binding %q", name)
		}
	}
	return nil
}
