package fs

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
)

func newHandle(
	path string,
) model.DataHandle {
	return handle{
		path: path,
	}
}

// handle allows interacting w/ data sourced from the filesystem
type handle struct {
	path string
}

func (lh handle) GetContent(
	ctx context.Context,
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {
	return os.Open(filepath.Join(lh.path, contentPath))
}

func (lh handle) ListDescendants(
	ctx context.Context,
) (
	[]*model.DirEntry,
	error,
) {
	return lh.rListDescendants(lh.path)
}

// rListDescendants recursively lists descendants of the current data node
func (lh handle) rListDescendants(
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

func (lh handle) Path() *string {
	return &lh.path
}

func (lh handle) Ref() string {
	return lh.path
}
