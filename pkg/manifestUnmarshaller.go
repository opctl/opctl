package pkg

//go:generate counterfeiter -o ./fakeManifestUnmarshaller.go --fake-name fakeManifestUnmarshaller ./ manifestUnmarshaller

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/virtual-go/vioutil"
	"gopkg.in/yaml.v2"
	"path"
)

type manifestUnmarshaller interface {
	Unmarshal(
		pkgRef string,
	) (*model.PkgManifest, error)
}

func newManifestUnmarshaller(
	ioUtil vioutil.VIOUtil,
	validator validator,
) manifestUnmarshaller {

	return _manifestUnmarshaller{
		ioUtil:    ioUtil,
		validator: validator,
	}

}

type _manifestUnmarshaller struct {
	ioUtil    vioutil.VIOUtil
	validator validator
}

func (this _manifestUnmarshaller) Unmarshal(
	pkgRef string,
) (*model.PkgManifest, error) {

	var err error

	// 1) ensure valid
	errs := this.validator.Validate(pkgRef)
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
	packageManifestPath := path.Join(pkgRef, ManifestFileName)

	packageManifestBytes, err := this.ioUtil.ReadFile(
		packageManifestPath,
	)
	if nil != err {
		return nil, err
	}

	pkgManifest := model.PkgManifest{}
	return &pkgManifest, yaml.Unmarshal(packageManifestBytes, &pkgManifest)

}
