package pkg

//go:generate counterfeiter -o ./fakeManifestUnmarshaller.go --fake-name fakeManifestUnmarshaller ./ manifestUnmarshaller

import (
	"bytes"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	"github.com/opspec-io/sdk-golang/model"
	"gopkg.in/yaml.v2"
	"path"
)

type manifestUnmarshaller interface {
	Unmarshal(
		pkgRef string,
	) (*model.PkgManifest, error)
}

func newManifestUnmarshaller(
	ioUtil iioutil.Iioutil,
	manifestValidator manifestValidator,
) manifestUnmarshaller {

	return _manifestUnmarshaller{
		ioUtil:            ioUtil,
		manifestValidator: manifestValidator,
	}

}

type _manifestUnmarshaller struct {
	ioUtil            iioutil.Iioutil
	manifestValidator manifestValidator
}

func (this _manifestUnmarshaller) Unmarshal(
	pkgRef string,
) (*model.PkgManifest, error) {

	var err error

	// 1) ensure valid
	errs := this.manifestValidator.Validate(pkgRef)
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
	packageManifestPath := path.Join(pkgRef, OpDotYmlFileName)

	packageManifestBytes, err := this.ioUtil.ReadFile(
		packageManifestPath,
	)
	if nil != err {
		return nil, err
	}

	pkgManifest := model.PkgManifest{}
	return &pkgManifest, yaml.Unmarshal(packageManifestBytes, &pkgManifest)

}
