package op

//go:generate counterfeiter -o ./fakeCreator.go --fake-name FakeCreator ./ Creator

import (
	"github.com/ghodss/yaml"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/opfile"
	"path/filepath"
)

type Creator interface {
	// Create creates an operation
	Create(
		path,
		pkgName,
		pkgDescription string,
	) error
}

// NewCreator returns an initialized Creator instance
func NewCreator() Creator {
	return _creator{
		os:     ios.New(),
		ioUtil: iioutil.New(),
	}
}

type _creator struct {
	os     ios.IOS
	ioUtil iioutil.IIOUtil
}

func (cr _creator) Create(
	path,
	pkgName,
	pkgDescription string,
) error {

	err := cr.os.MkdirAll(
		path,
		0777,
	)
	if nil != err {
		return err
	}

	opDotYml := model.OpDotYml{
		Description: pkgDescription,
		Name:        pkgName,
	}

	opDotYmlBytes, err := yaml.Marshal(&opDotYml)
	if nil != err {
		return err
	}

	return cr.ioUtil.WriteFile(
		filepath.Join(path, dotyml.FileName),
		opDotYmlBytes,
		0777,
	)

}
