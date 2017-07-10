package pkg

import "github.com/opspec-io/sdk-golang/model"

func (this _Pkg) GetManifest(
	pkgHandle Handle,
) (
	*model.PkgManifest,
	error,
) {
	manifestReader, err := pkgHandle.GetContent(OpDotYmlFileName)
	if nil != err {
		return nil, err
	}

	return this.manifest.Unmarshal(manifestReader)
}
