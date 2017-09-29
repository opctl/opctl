package files

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Files

import (
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/filecopier"
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
)

type Files interface {
	Interpret(
		pkgHandle model.PkgHandle,
		scope map[string]*model.Value,
		scgContainerCallFiles map[string]interface{},
		scratchDirPath string,
	) (map[string]string, error)
}

func New(
	rootFSPath string,
) Files {
	return _Files{
		expression: expression.New(),
		fileCopier: filecopier.New(),
		os:         ios.New(),
		rootFSPath: rootFSPath,
	}
}

type _Files struct {
	expression expression.Expression
	fileCopier filecopier.FileCopier
	os         ios.IOS
	rootFSPath string
}
