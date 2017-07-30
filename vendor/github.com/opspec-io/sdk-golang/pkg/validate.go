package pkg

import "context"

func (this _Pkg) Validate(
	pkgHandle Handle,
) []error {
	manifestReader, err := pkgHandle.GetContent(
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
