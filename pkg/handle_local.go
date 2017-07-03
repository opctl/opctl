package pkg

import (
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func newHandleLocal(
	path string,
) handle {
	return handleLocal{
		ioUtil: iioutil.New(),
		os:     ios.New(),
		path:   path,
	}
}

type handleLocal struct {
	ioUtil iioutil.Iioutil
	os     ios.IOS
	path   string
}

func (lh handleLocal) ListContents() (
	[]*model.PkgContent,
	error,
) {
	return lh.listContents(lh.path)
}

// listContents is same as ListContents but requires path as a param
func (lh handleLocal) listContents(
	path string,
) (
	[]*model.PkgContent,
	error,
) {
	childFileInfos, err := lh.ioUtil.ReadDir(path)
	if nil != err {
		return nil, err
	}

	var contents []*model.PkgContent
	for _, contentFileInfo := range childFileInfos {

		contentPath := filepath.Join(path, contentFileInfo.Name())

		if contentFileInfo.IsDir() {
			// recurse into child dirs
			childContents, err := lh.listContents(contentPath)
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

func (lh handleLocal) GetContent(
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {
	return lh.os.Open(filepath.Join(lh.path, contentPath))
}
