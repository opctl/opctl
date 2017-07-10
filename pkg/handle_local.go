package pkg

import (
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func newLocalHandle(
	path string,
) Handle {
	return localHandle{
		ioUtil: iioutil.New(),
		os:     ios.New(),
		path:   path,
	}
}

// localHandle is a provider which operates on a relative pkg ref
type localHandle struct {
	ioUtil iioutil.IIOUtil
	os     ios.IOS
	path   string
}

func (lh localHandle) ListContents() (
	[]*model.PkgContent,
	error,
) {
	return lh.rListContents(lh.path)
}

// rListContents recursively lists pkg contents at path
func (lh localHandle) rListContents(
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
			childContents, err := lh.rListContents(contentPath)
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

func (lh localHandle) GetContent(
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {
	return lh.os.Open(filepath.Join(lh.path, contentPath))
}

func (lh localHandle) Ref() string {
	return lh.Ref()
}
