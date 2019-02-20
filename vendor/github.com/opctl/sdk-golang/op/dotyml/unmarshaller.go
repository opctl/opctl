package dotyml

//go:generate counterfeiter -o ./fakeUnmarshaller.go --fake-name FakeUnmarshaller ./ Unmarshaller

import (
	"bytes"
	"fmt"
	"github.com/opctl/sdk-golang/model"
	"gopkg.in/yaml.v2"
)

// @TODO make private
type Unmarshaller interface {
	// Unmarshal validates and unmarshals an "op.yml" file
	Unmarshal(
		manifestBytes []byte,
	) (*model.OpDotYml, error)
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
	manifestBytes []byte,
) (*model.OpDotYml, error) {

	var err error

	// 1) ensure valid
	errs := uml.validator.Validate(manifestBytes)
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
	opDotYml := model.OpDotYml{}
	return &opDotYml, yaml.Unmarshal(manifestBytes, &opDotYml)

}
