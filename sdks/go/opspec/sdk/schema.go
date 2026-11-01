package sdk

import (
	"encoding/json"
	"fmt"
)

const minSchemaVersion = 1
const schemaVersion = 4

type schemaRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Decode through pointers so a null field cannot retain an earlier duplicate's integer.
func (schema *schemaRange) UnmarshalJSON(data []byte) error {
	var fields struct {
		Min *int `json:"min"`
		Max *int `json:"max"`
	}
	if err := json.Unmarshal(data, &fields); err != nil || fields.Min == nil || fields.Max == nil {
		return fmt.Errorf("invalid opspec schema range in bootstrap request: min and max must be integers")
	}
	schema.Min, schema.Max = *fields.Min, *fields.Max
	return nil
}

func commonSchema(schema schemaRange) (int, error) {
	if schema.Min < 1 || schema.Max < schema.Min {
		return 0, fmt.Errorf("invalid opspec schema range in bootstrap request")
	}
	if schema.Min > schemaVersion {
		return 0, fmt.Errorf("no common opspec schema: node requires [%d,%d], Go op SDK supports [%d,%d]; upgrade the Go op SDK", schema.Min, schema.Max, minSchemaVersion, schemaVersion)
	}
	return min(schema.Max, schemaVersion), nil
}

func selectSchema(schema schemaRange, call *callSpec) (int, error) {
	target, err := commonSchema(schema)
	if err != nil {
		return 0, err
	}
	needed, path := schemaRequirement(call, "run")
	if needed > target {
		feature := "container stdout/stderr bindings"
		if needed == 3 {
			feature = "container outputs file"
		}
		if needed == 4 {
			feature = "typed container outputs"
		}
		return 0, fmt.Errorf("this op uses %s (opspec >= %d); upgrade opctl to a version supporting opspec %d (at %s)", feature, needed, needed, path)
	}
	return target, nil
}

func schemaRequirement(call *callSpec, path string) (int, string) {
	version, field := 1, ""
	if call == nil {
		return version, field
	}
	if c := call.Container; c != nil {
		if c.Stdout != "" || c.Stderr != "" {
			version, field = 2, path+".container.stdout/stderr"
		}
		if len(c.Outputs) > 0 {
			version, field = 3, path+".container.outputs"
		}
		if len(c.OutputParams) > 0 {
			version, field = 4, path+".container.outputParams"
		}
	}
	for i, group := range [][]*callSpec{call.Serial, call.Parallel} {
		for j, child := range group {
			needed, childField := schemaRequirement(child, fmt.Sprintf("%s.%s[%d]", path, []string{"serial", "parallel"}[i], j))
			if needed > version {
				version, field = needed, childField
			}
		}
	}
	return version, field
}
