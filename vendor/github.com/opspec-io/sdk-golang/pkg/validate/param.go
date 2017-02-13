package validate

import (
	"errors"
	"fmt"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

// validates an value against a parameter
func (this validate) Param(
	value *model.Data,
	param *model.Param,
) (errs []error) {
	if nil == param {
		// panic as errs represents validation errors not execution errors
		panic("param required")
	}

	switch {
	case nil != param.String:
		errs = this.stringParam(value, param.String)
	case nil != param.Socket:
		errs = this.socketParam(value, param.Socket)
	}
	return
}

// validates an value against a string parameter
func (this validate) stringParam(
	rawValue *model.Data,
	param *model.StringParam,
) (errs []error) {
	errs = []error{}

	// handle no value passed
	if nil == rawValue {
		errs = append(errs, errors.New("String required"))
		return
	}

	value := rawValue.String
	if "" == value && "" != param.Default {
		// apply default if value not set
		value = param.Default
	}

	// guard no constraints
	if nil == param.Constraints {
		return
	}

	constraintsJsonBytes, err := format.NewJsonFormat().From(param.Constraints)
	if err != nil {
		// panic as errs represents validation errors not execution errors
		panic(err.Error())
	}

	result, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(string(constraintsJsonBytes)),
		gojsonschema.NewStringLoader(fmt.Sprintf(`"%v"`, value)),
	)
	if err != nil {
		// panic as errs represents validation errors not execution errors
		panic(err.Error())
	}

	for _, errString := range result.Errors() {
		// enum validation errors include `(root) ` prefix we don't want
		errs = append(errs, errors.New(strings.TrimPrefix(errString.Description(), "(root) ")))
	}

	return
}

// validates an value against a network socket parameter
func (this validate) socketParam(
	rawValue *model.Data,
	param *model.SocketParam,
) (errs []error) {
	errs = []error{}

	// handle no value passed
	if nil == rawValue || "" == rawValue.Socket {
		errs = append(errs, errors.New("Socket required"))
	}
	return
}
