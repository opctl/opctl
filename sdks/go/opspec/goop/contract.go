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

// NewRequest requests a plan without inputs or known outputs.
func NewRequest() Request {
	return Request{ContractVersion: 1, Phase: "plan", OpspecSchema: SchemaRange{Min: 1, Max: opfile.SchemaVersion}, Inputs: map[string]interface{}{}, KnownOutputs: map[string]interface{}{}}
}

// ParseResponse validates a single bootstrap response using the op file schema.
func ParseResponse(data []byte) (*model.OpSpec, error) {
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
	if len(spec.Inputs) != 0 || len(spec.Outputs) != 0 {
		return nil, fmt.Errorf("Go op inputs and outputs are not supported yet")
	}
	if spec.Run == nil {
		return nil, fmt.Errorf("Go op plan must include a run call")
	}
	return spec, nil
}
