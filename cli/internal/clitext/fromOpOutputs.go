package clitext

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"strconv"

	"github.com/ghodss/yaml"
)

// fromValue formats a value to CLI text
func fromValue(value model.Value) (interface{}, error) {
	if nil != value.Array {
		formattedArray := []interface{}{}
		for itemKey, itemValue := range *value.Array {
			switch typedItemValue := itemValue.(type) {
			case model.Value:
				formattedValue, err := fromValue(typedItemValue)
				if nil != err {
					return "", fmt.Errorf("unable to stringify item '%v' from array: %w", itemKey, err)
				}

				formattedArray = append(formattedArray, formattedValue)
			default:
				formattedArray = append(formattedArray, itemValue)
			}
		}
		return formattedArray, nil
	} else if nil != value.Boolean {
		return strconv.FormatBool(*value.Boolean), nil
	} else if nil != value.Dir {
		return *value.Dir, nil
	} else if nil != value.File {
		return *value.File, nil
	} else if nil != value.Number {
		return fmt.Sprintf("%f", *value.Number), nil
	} else if nil != value.Object {
		formattedMap := map[string]interface{}{}
		for propKey, propValue := range *value.Object {
			switch typedPropValue := propValue.(type) {
			case model.Value:
				var err error
				if formattedMap[propKey], err = fromValue(typedPropValue); nil != err {
					return "", fmt.Errorf("unable to stringify property '%v' from object: %w", propKey, err)
				}
			default:
				formattedMap[propKey] = propValue
			}
		}
		return formattedMap, nil
	} else if nil != value.Socket {
		return *value.Socket, nil
	} else if nil != value.String {
		return *value.String, nil
	}
	return "", fmt.Errorf("unable to stringify value '%+v'", value)
}

// FromOpOutputs formats op outputs to CLI text
func FromOpOutputs(outputs map[string]*model.Value) (string, error) {
	if len(outputs) == 0 {
		return "outputs: none\n", nil
	}

	formattedValues := map[string]interface{}{}
	for name, value := range outputs {
		description, err := fromValue(*value)
		if err != nil {
			return "", err
		}
		formattedValues[name] = description
	}
	bytes, err := yaml.Marshal(
		map[string]interface{}{"outputs": formattedValues},
	)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
