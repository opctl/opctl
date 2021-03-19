package opfile

import (
	"bytes"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/sdks/go/model"
)

// Unmarshal validates and unmarshals an "op.yml" file
func Unmarshal(
	opFileBytes []byte,
) (*model.OpSpec, error) {

	var err error

	// 1) ensure valid
	errs := Validate(opFileBytes)
	if len(errs) > 0 {
		messageBuffer := bytes.NewBufferString(
			fmt.Sprint(`
-
  Error(s):`))
		for _, validationError := range errs {
			messageBuffer.WriteString(fmt.Sprintf(`
    - %v`, validationError.Error()))
		}
		err = fmt.Errorf(
			`%v
-`, messageBuffer.String())
	}
	if err != nil {
		return nil, err
	}

	// 2) build
	opFile := model.OpSpec{}
	return &opFile, yaml.Unmarshal(opFileBytes, &opFile)

}
