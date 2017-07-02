package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func (this _Pkg) ListContents(
	pkgRef string,
) (
	[]*model.PkgContent,
	error,
) {
	childFileInfos, err := this.ioUtil.ReadDir(pkgRef)
	if nil != err {
		return nil, err
	}

	var contents []*model.PkgContent
	for _, contentFileInfo := range childFileInfos {

		contentPath := filepath.Join(pkgRef, contentFileInfo.Name())

		if contentFileInfo.IsDir() {
			// recurse into child dirs
			childContents, err := this.ListContents(contentPath)
			if nil != err {
				return nil, err
			}
			contents = append(contents, childContents...)
		} else {
			contents = append(
				contents,
				&model.PkgContent{
					Path: contentPath,
					Size: contentFileInfo.Size(),
				},
			)
		}

	}

	return contents, err
}
