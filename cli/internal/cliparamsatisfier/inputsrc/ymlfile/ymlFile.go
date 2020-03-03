package ymlfile

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc"
)

// New constructs a new ymlFileInputSrc
func New(
	filePath string,
) (inputsrc.InputSrc, error) {

	_, err := os.Stat(filePath)
	if nil != err {
		if !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
		// if file doesn't exist, treat as empty
		return ymlFileInputSrc{map[string]interface{}{}}, nil
	}

	ymlBytes, err := ioutil.ReadFile(filePath)
	if nil != err {
		return nil, err
	}

	// convert YAML to JSON before unmarshalling;
	// otherwise nested objects cause : "json: unsupported type: map[interface {}]interface {}"
	jsonBytes, err := yaml.YAMLToJSON(ymlBytes)
	if nil != err {
		return nil, err
	}

	argMap := map[string]interface{}{}

	if err := json.Unmarshal(jsonBytes, &argMap); nil != err {
		return nil, err
	}

	return ymlFileInputSrc{argMap}, nil
}

// ymlFileInputSrc implements InputSrc interface by sourcing inputs from a file containing a yml map
type ymlFileInputSrc struct {
	argMap map[string]interface{}
}

func (this ymlFileInputSrc) ReadString(
	inputName string,
) (*string, bool) {
	if inputValue, ok := this.argMap[inputName]; ok {
		// enforce read at most once.
		delete(this.argMap, inputName)
		switch inputValue := inputValue.(type) {
		case string:
			return &inputValue, true
		default:
			inputValueBytes, err := json.Marshal(inputValue)
			if nil != err {
				return nil, false
			}

			inputValueString := string(inputValueBytes)

			return &inputValueString, true
		}
	}
	return nil, false
}
