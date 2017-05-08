package cliparamsatisfier

import (
	"github.com/golang-interfaces/iioutil"
	"gopkg.in/yaml.v2"
)

// NewYMLFileInputSrc constructs a new ymlFileInputSrc, ignoring any errors encountered
func NewYMLFileInputSrc(
	filePath string,
	ioutil iioutil.Iioutil,
) InputSrc {
	argMap := map[string]string{}

	// ignore errors; make best effort
	ymlBytes, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(ymlBytes, &argMap)

	return ymlFileInputSrc{argMap}
}

// ymlFileInputSrc implements InputSrc interface by sourcing inputs from a file containing a yml map
type ymlFileInputSrc struct {
	argMap map[string]string
}

func (this ymlFileInputSrc) Read(
	inputName string,
) *string {
	if inputValue, ok := this.argMap[inputName]; ok {
		// enforce read at most once.
		delete(this.argMap, inputName)

		return &inputValue
	}
	return nil
}
