package opfile

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/sdks/go/model"
)

// @TODO make private
//counterfeiter:generate -o fakes/unmarshaller.go . Unmarshaller
type Unmarshaller interface {
	// Unmarshal validates and unmarshals an "op.yml" file
	Unmarshal(
		opFileBytes []byte,
	) (*model.OpFile, error)
}

// NewUnmarshaller returns an initialized Unmarshaller instance
func NewUnmarshaller() Unmarshaller {
	return _unmarshaller{
		validator: newValidator(),
	}
}

type _unmarshaller struct {
	validator validator
}

func (uml _unmarshaller) Unmarshal(
	opFileBytes []byte,
) (*model.OpFile, error) {

	var err error

	// 1) ensure valid
	errs := uml.validator.Validate(opFileBytes)
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
	if nil != err {
		return nil, err
	}

	// 2) build
	opFile := model.OpFile{}
	return &opFile, yaml.Unmarshal(opFileBytes, &opFile)

}
