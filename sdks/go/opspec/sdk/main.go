package sdk

import (
	"encoding/json"
	"fmt"
	"os"
)

type request struct {
	ContractVersion int    `json:"contractVersion"`
	Phase           string `json:"phase"`
	OpspecSchema    struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"opspecSchema"`
	Inputs       map[string]json.RawMessage `json:"inputs"`
	KnownOutputs map[string]json.RawMessage `json:"knownOutputs"`
}

// Main serves the describe or plan bootstrap phase requested by opctl.
// Diagnostics go to stderr; stdout contains only the JSON response.
func Main(fn func() Call) {
	if err := run(fn); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fn func() Call) error {
	if len(os.Args) != 2 {
		return fmt.Errorf("expected bootstrap phase: describe or plan")
	}
	phase := os.Args[1]
	if phase != "describe" && phase != "plan" {
		return fmt.Errorf("unsupported bootstrap phase %q", phase)
	}
	data, err := os.ReadFile(os.Getenv("OPCTL_BOOTSTRAP_REQUEST"))
	if err != nil {
		return fmt.Errorf("read bootstrap request: %w", err)
	}
	var req request
	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("decode bootstrap request: %w", err)
	}
	if req.ContractVersion != 1 {
		return fmt.Errorf("unsupported bootstrap contract version %d", req.ContractVersion)
	}
	if req.Phase != phase {
		return fmt.Errorf("bootstrap request phase %q differs from argv phase %q", req.Phase, phase)
	}
	if req.OpspecSchema.Min > 1 || req.OpspecSchema.Max < 1 || req.OpspecSchema.Min > req.OpspecSchema.Max {
		return fmt.Errorf("Go op requires opspec schema 1; upgrade opctl to a version supporting it")
	}
	if len(req.Inputs) != 0 {
		return fmt.Errorf("Go op inputs are not supported yet")
	}
	if len(req.KnownOutputs) != 0 {
		return fmt.Errorf("Go op known outputs are not supported yet")
	}
	response, err := emit(phase, fn)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(response)
	return err
}
