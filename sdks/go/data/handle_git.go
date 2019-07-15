package data

import (
	"context"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
)

func newGitHandle(
	path string,
	dataRef string,
) model.DataHandle {
	return gitHandle{
		ioUtil:  iioutil.New(),
		os:      ios.New(),
		path:    path,
		dataRef: dataRef,
	}
}

// gitHandle allows interacting w/ data sourced from git
type gitHandle struct {
	ioUtil  iioutil.IIOUtil
	os      ios.IOS
	path    string
	dataRef string
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

func (gh gitHandle) ListDescendants(
	ctx context.Context,
) (
	[]*model.DirEntry,
	error,
) {
	return gh.rListDescendants(gh.path)
}

// rListDescendants recursively lists descendants of the current data node
func (gh gitHandle) rListDescendants(
	path string,
) (
	[]*model.DirEntry,
	error,
) {
	childFileInfos, err := gh.ioUtil.ReadDir(path)
	if nil != err {
		return nil, err
	}

	var contents []*model.DirEntry
	for _, contentFileInfo := range childFileInfos {
		absContentPath := filepath.Join(path, contentFileInfo.Name())

		if contentFileInfo.IsDir() {
			// recurse into child dirs
			childContents, err := gh.rListDescendants(absContentPath)
			if nil != err {
				return nil, err
			}
			contents = append(contents, childContents...)
		}

		relContentPath, err := filepath.Rel(gh.path, absContentPath)
		if nil != err {
			return nil, err
		}

		contents = append(
			contents,
			&model.DirEntry{
				Mode: contentFileInfo.Mode(),
				Path: filepath.Join(string(os.PathSeparator), relContentPath),
				Size: contentFileInfo.Size(),
			},
		)

	}

	return contents, err
}

func (gh gitHandle) Path() *string {
	return &gh.path
}

func (gh gitHandle) Ref() string {
	return gh.dataRef
}
