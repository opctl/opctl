package data

import (
	"context"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/types"
)

func newFSHandle(
	path string,
) types.DataHandle {
	return fsHandle{
		ioUtil: iioutil.New(),
		os:     ios.New(),
		path:   path,
	}
}

// fsHandle allows interacting w/ data sourced from the filesystem
type fsHandle struct {
	ioUtil iioutil.IIOUtil
	os     ios.IOS
	path   string
}

func (lh fsHandle) GetContent(
	ctx context.Context,
	contentPath string,
) (
	types.ReadSeekCloser,
	error,
) {
	return lh.os.Open(filepath.Join(lh.path, contentPath))
}

func (lh fsHandle) ListDescendants(
	ctx context.Context,
) (
	[]*types.DirEntry,
	error,
) {
	return lh.rListDescendants(lh.path)
}

// rListDescendants recursively lists descendants of the current data node
func (lh fsHandle) rListDescendants(
	path string,
) (
	[]*types.DirEntry,
	error,
) {
	childFileInfos, err := lh.ioUtil.ReadDir(path)
	if nil != err {
		return nil, err
	}

	var contents []*types.DirEntry
	for _, contentFileInfo := range childFileInfos {

		absContentPath := filepath.Join(path, contentFileInfo.Name())

		if contentFileInfo.IsDir() {
			// recurse into child dirs
			childContents, err := lh.rListDescendants(absContentPath)
			if nil != err {
				return nil, err
			}
			contents = append(contents, childContents...)
		}

		relContentPath, err := filepath.Rel(lh.path, absContentPath)
		if nil != err {
			return nil, err
		}
		contents = append(
			contents,
			&types.DirEntry{
				Mode: contentFileInfo.Mode(),
				Path: filepath.Join(string(os.PathSeparator), relContentPath),
				Size: contentFileInfo.Size(),
			},
		)

	}

	return contents, err
}

func (lh fsHandle) Path() *string {
	return &lh.path
}

func (lh fsHandle) Ref() string {
	return lh.path
}
