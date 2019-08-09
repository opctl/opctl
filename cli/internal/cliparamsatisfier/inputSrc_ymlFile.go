package cliparamsatisfier

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"os"
)

// NewYMLFileInputSrc constructs a new ymlFileInputSrc
func (isf _InputSrcFactory) NewYMLFileInputSrc(
	filePath string,
) (InputSrc, error) {

	_, err := isf.os.Stat(filePath)
	if nil != err {
		if !os.IsNotExist(err) {
			// return actual errors
			return nil, err
		}
		// if file doesn't exist, treat as empty
		return ymlFileInputSrc{map[string]interface{}{}}, nil
	}

	ymlBytes, err := isf.ioutil.ReadFile(filePath)
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

	if err := isf.json.Unmarshal(jsonBytes, &argMap); nil != err {
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
