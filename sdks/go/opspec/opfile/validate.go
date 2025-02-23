package opfile

import (
	"encoding/json"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/opspec/opfile"
	"github.com/opctl/opctl/sdks/go/internal/jsonschema"
)

// Validate validates an "op.yml"
func Validate(
	opFileBytes []byte,
) []error {

	var unmarshalledSchema map[string]interface{}
	err := json.Unmarshal(opfile.JsonSchemaBytes, &unmarshalledSchema)
	if err != nil {
		return []error{err}
	}

	var unmarshalledOp map[string]interface{}
	err = yaml.Unmarshal(opFileBytes, &unmarshalledOp)
	if err != nil {
		// handle syntax errors specially
		return []error{err}
	}

	return jsonschema.Validate(unmarshalledOp, unmarshalledSchema)
}
