package pkg

import (
	"context"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

func newGitHandle(
	path string,
	pkgRef string,
) Handle {
	return gitHandle{
		ioUtil: iioutil.New(),
		os:     ios.New(),
		path:   path,
		pkgRef: pkgRef,
	}
}

// gitHandle allows interacting w/ a package sourced from the filesystem
type gitHandle struct {
	ioUtil iioutil.IIOUtil
	os     ios.IOS
	path   string
	pkgRef string
}

func (gh gitHandle) GetContent(
	ctx context.Context,
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {
	return gh.os.Open(filepath.Join(gh.path, contentPath))
}

func (gh gitHandle) ListContents(
	ctx context.Context,
) (
	[]*model.PkgContent,
	error,
) {
	return gh.rListContents(gh.path)
}

// rListContents recursively lists pkg contents at path
func (gh gitHandle) rListContents(
	path string,
) (
	[]*model.PkgContent,
	error,
) {
	childFileInfos, err := gh.ioUtil.ReadDir(path)
	if nil != err {
		return nil, err
	}

	var contents []*model.PkgContent
	for _, contentFileInfo := range childFileInfos {
		absContentPath := filepath.Join(path, contentFileInfo.Name())

		if contentFileInfo.IsDir() {
			// recurse into child dirs
			childContents, err := gh.rListContents(absContentPath)
			if nil != err {
				return nil, err
			}
			contents = append(contents, childContents...)
		} else {
			contents = append(
				contents,
				&model.PkgContent{
					Path: strings.TrimPrefix(absContentPath, gh.path),
					Size: contentFileInfo.Size(),
				},
			)
		}

	}

	return contents, err
}

func (gh gitHandle) Ref() string {
	return gh.pkgRef
}
