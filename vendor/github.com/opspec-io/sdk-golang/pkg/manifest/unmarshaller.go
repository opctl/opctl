package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/opspec-io/sdk-golang/model"
)

func (this _Manifest) Unmarshal(
	manifestBytes []byte,
) (*model.PkgManifest, error) {

	var err error

	// 1) ensure valid
	errs := this.Validate(manifestBytes)
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
	pkgManifest := model.PkgManifest{}

	manifestJSONBytes, err := yaml.YAMLToJSON([]byte(manifestBytes))
	if nil != err {
		return nil, err
	}
	return &pkgManifest, json.Unmarshal(manifestJSONBytes, &pkgManifest)

}
