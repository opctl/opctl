package manifest

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"gopkg.in/yaml.v2"
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
	return &pkgManifest, yaml.Unmarshal(manifestBytes, &pkgManifest)

}
