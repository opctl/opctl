package git

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
)

func newHandle(
	path string,
	dataRef string,
) model.DataHandle {
	return handle{
		path:    path,
		dataRef: dataRef,
	}
}

// handle allows interacting w/ data sourced from git
type handle struct {
	path    string
	dataRef string
}

func (gh handle) GetContent(
	ctx context.Context,
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {
	return os.Open(filepath.Join(gh.path, contentPath))
}

func (gh handle) ListDescendants(
	ctx context.Context,
) (
	[]*model.DirEntry,
	error,
) {
	return gh.rListDescendants(gh.path)
}

// rListDescendants recursively lists descendants of the current data node
func (gh handle) rListDescendants(
	path string,
) (
	[]*model.DirEntry,
	error,
) {
	childFileInfos, err := ioutil.ReadDir(path)
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

func (gh handle) Path() *string {
	return &gh.path
}

func (gh handle) Ref() string {
	return gh.dataRef
}
