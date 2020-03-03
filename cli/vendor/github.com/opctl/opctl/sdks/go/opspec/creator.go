package opspec

import (
	"github.com/ghodss/yaml"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"path/filepath"
)

//counterfeiter:generate -o fakes/creator.go . Creator
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

	opFile := model.OpFile{
		Description: pkgDescription,
		Name:        pkgName,
	}

	opFileBytes, err := yaml.Marshal(&opFile)
	if nil != err {
		return err
	}

	return cr.ioUtil.WriteFile(
		filepath.Join(path, opfile.FileName),
		opFileBytes,
		0777,
	)

}
