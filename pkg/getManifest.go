package pkg

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
)

func (this _Pkg) GetManifest(
	opDirHandle model.DataHandle,
) (
	*model.PkgManifest,
	error,
) {
	manifestReader, err := opDirHandle.GetContent(context.TODO(), OpDotYmlFileName)
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
