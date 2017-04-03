package pkg

//go:generate counterfeiter -o ./fakeViewFactory.go --fake-name fakeViewFactory ./ viewFactory

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/virtual-go/vioutil"
	"path"
)

type viewFactory interface {
	Construct(
		pkgRef string,
	) (
		packageView *model.PackageView,
		err error,
	)
}

func newViewFactory(
	ioUtil vioutil.VIOUtil,
	validator validator,
	yaml format.Format,
) viewFactory {

	return &_viewFactory{
		ioUtil:    ioUtil,
		validator: validator,
		yaml:      yaml,
	}

}

type _viewFactory struct {
	ioUtil    vioutil.VIOUtil
	validator validator
	yaml      format.Format
}

func (this _viewFactory) Construct(
	pkgRef string,
) (
	packageView *model.PackageView,
	err error,
) {

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
		return
	}

	// 2) build
	packageManifestPath := path.Join(pkgRef, NameOfPkgManifestFile)

	packageManifestBytes, err := this.ioUtil.ReadFile(
		packageManifestPath,
	)
	if nil != err {
		return
	}

	packageManifestView := model.PackageManifestView{}
	err = this.yaml.To(
		packageManifestBytes,
		&packageManifestView,
	)
	if nil != err {
		return
	}

	packageView = &model.PackageView{
		Description: packageManifestView.Description,
		Inputs:      packageManifestView.Inputs,
		Name:        packageManifestView.Name,
		Outputs:     packageManifestView.Outputs,
		Run:         packageManifestView.Run,
		Version:     packageManifestView.Version,
	}

	return

}
