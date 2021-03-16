package git

import (
	"context"
	"io/fs"
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
	childFiles, err := os.ReadDir(path)
	if nil != err {
		return nil, err
	}

	rootFS := os.DirFS(path)

	var contents []*model.DirEntry
	for _, childFile := range childFiles {
		absContentPath := filepath.Join(path, childFile.Name())

		if childFile.IsDir() {
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

		fileInfo, err := fs.Stat(rootFS, childFile.Name())
		if err != nil {
			return nil, err
		}

		contents = append(
			contents,
			&model.DirEntry{
				Mode: fileInfo.Mode(),
				Path: filepath.Join(string(os.PathSeparator), relContentPath),
				Size: fileInfo.Size(),
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
