package pkg

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
)

func (this _Pkg) Validate(
	opDirHandle model.DataHandle,
) []error {
	manifestReader, err := opDirHandle.GetContent(
		context.TODO(),
		OpDotYmlFileName,
	)
	if nil != err {
		return []error{err}
	}

	manifestBytes, err := this.ioUtil.ReadAll(manifestReader)
	if nil != err {
		return []error{err}
	}

	return this.manifest.Validate(manifestBytes)
}
