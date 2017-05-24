package interpreter

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (itp _Interpreter) Interpret(
	name,
	value string,
	params map[string]*model.Param,
	scope map[string]*model.Data,
) (*model.Data, error) {

	var dcgValue *model.Data
	param, ok := params[name]
	if !ok {
		return nil, fmt.Errorf("Unable to bind to '%v'. '%v' is not a defined input", name, name)
	} else if "" == value {
		// implicit arg
		dcgValue, ok = scope[name]
		if !ok {
			return nil, fmt.Errorf("Unable to bind to '%v' via implicit ref. '%v' is not in scope", name, name)
		}
	} else if explicitRef := strings.TrimSuffix(strings.TrimPrefix(value, "$("), ")"); len(explicitRef) == (len(value) - 3) {
		// explicit arg
		dcgValue, ok = scope[explicitRef]
		if !ok {
			return nil, fmt.Errorf("Unable to bind '%v' to '%v' via explicit ref. '%v' is not in scope", name, explicitRef, explicitRef)
		}
	} else {
		interpolatedVal := itp.interpolater.Interpolate(value, scope)
		switch {
		// interpolated arg
		case nil != param.String:
			dcgValue = &model.Data{String: &interpolatedVal}
		case nil != param.Dir:
			interpolatedVal = itp.rootPath(interpolatedVal)
			dcgValue = &model.Data{Dir: &interpolatedVal}
		case nil != param.Number:
			floatVal, err := strconv.ParseFloat(interpolatedVal, 64)
			if nil != err {
				return nil, fmt.Errorf("Unable to bind '%v' to '%v' as number; error was: '%v'", name, interpolatedVal, err.Error())
			}
			dcgValue = &model.Data{Number: &floatVal}
		case nil != param.File:
			interpolatedVal = itp.rootPath(interpolatedVal)
			dcgValue = &model.Data{File: &interpolatedVal}
		case nil != param.Socket:
			return nil, fmt.Errorf("Unable to bind '%v' to '%v'; sockets must be passed by reference", name, interpolatedVal)
		}
	}

	if err := itp.validateParam(name, dcgValue, param); nil != err {
		return nil, err
	}

	return dcgValue, nil
}

// rootPath ensures paths are rooted (interpreted as having no parent) this ensures parent paths of input mounts aren't
// accessible (which would break encapsulation)
func (itp _Interpreter) rootPath(
	path string,
) string {
	path = strings.Replace(path, "../", string(os.PathSeparator), -1)
	path = strings.Replace(path, "..\\", string(os.PathSeparator), -1)
	return filepath.Clean(path)
}

func (itp _Interpreter) validateParam(
	name string,
	value *model.Data,
	param *model.Param,
) error {

	var (
		argDisplayValue string
	)

	switch {
	case nil != param.Dir:
		argDisplayValue = *value.Dir
	case nil != param.File:
		argDisplayValue = *value.File
	case nil != param.Number:
		if param.Number.IsSecret {
			argDisplayValue = "************"
		} else {
			argDisplayValue = strconv.FormatFloat(*value.Number, 'f', -1, 64)
		}
	case nil != param.Socket:
		argDisplayValue = *value.Socket
	case nil != param.String:
		if param.String.IsSecret {
			argDisplayValue = "************"
		} else {
			argDisplayValue = *value.String
		}
	}

	// validator
	validationErrors := itp.validator.Validate(value, param)

	messageBuffer := bytes.NewBufferString(``)

	if len(validationErrors) > 0 {
		messageBuffer.WriteString(fmt.Sprintf(`
  Name: %v
  Value: %v
  Error(s):`, name, argDisplayValue),
		)
		for _, validationError := range validationErrors {
			messageBuffer.WriteString(fmt.Sprintf(`
    - %v`, validationError.Error()))
		}
		messageBuffer.WriteString(`
`)
	}

	if messageBuffer.Len() > 0 {
		return fmt.Errorf(`
-
  validation of the following input failed:
%v
-`, messageBuffer.String())
	}
	return nil
}
