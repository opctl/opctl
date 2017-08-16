package pkg

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
)

func (this _Pkg) GetManifest(
	pkgHandle model.PkgHandle,
) (
	*model.PkgManifest,
	error,
) {
	manifestReader, err := pkgHandle.GetContent(context.TODO(), OpDotYmlFileName)
	if nil != err {
		return nil, err
	}
	defer manifestReader.Close()

	manifestBytes, err := this.ioUtil.ReadAll(manifestReader)
	if nil != err {
		return nil, err
	}

	return this.manifest.Unmarshal(manifestBytes)
}
