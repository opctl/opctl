package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
)

var oneMB int64 = 1024 * 1024
var maxEmbedBytes int64 = 40 * oneMB

func getMaxEmbedMB() float64 {
	return float64(maxEmbedBytes) / float64(oneMB)
}

//embedLocalFilesAndDirs mutates args by embedding any referenced file &/or dirs as objects.
// This makes the args location independent and therefore runnable on a remote node.
// note1: this loads referenced files/dirs into memory so we limit max combined embed to maxEmbedBytes to reduce the chances of memory exhaustion.
// note2: this approach is primitive; we will need to support de-dupe and chunking in the future.
func embedLocalFilesAndDirs(
	args map[string]*model.Value,
) error {
	var totalSize int64
	for key, val := range args {
		var fileOrDirPath string
		if nil != val.Dir {
			fileOrDirPath = *val.Dir
		} else if nil != val.File {
			fileOrDirPath = *val.File
		} else {
			continue
		}

		fileObj, size, err := fileOrDirPathToObject(fileOrDirPath)
		if nil != err {
			return err
		}

		totalSize += size
		if totalSize > maxEmbedBytes {
			return fmt.Errorf("embedding failed: combined size of files/dirs cannot exceed %gMb", getMaxEmbedMB())
		}

		args[key] = &model.Value{Object: &fileObj}
	}
	return nil
}

// fileOrDirPathToObject recursively serializes a file or directory to an object
func fileOrDirPathToObject(
	path string,
) (map[string]interface{}, int64, error) {
	info, err := os.Stat(path)
	if nil != err {
		return nil, 0, err
	}

	if !info.IsDir() {
		if info.Size() >= maxEmbedBytes {
			return nil, 0, fmt.Errorf("%s is %gMb but cannot be bigger than %gMb", path, float64(info.Size())/float64(oneMB), getMaxEmbedMB())
		}

		body, err := ioutil.ReadFile(path)
		if nil != err {
			return nil, 0, err
		}

		return map[string]interface{}{
			"data": string(body),
		}, info.Size(), nil
	}

	childFileInfos, err := ioutil.ReadDir(path)
	if nil != err {
		return nil, 0, err
	}

	fileObject := map[string]interface{}{}
	var totalSize int64
	for _, childFileInfo := range childFileInfos {
		if childFileInfo.Size()+totalSize > maxEmbedBytes {
			return nil, 0, fmt.Errorf("embedding failed: %s cannot exceed %gMb", path, getMaxEmbedMB())
		}

		childFileObject, childSize, err := fileOrDirPathToObject(filepath.Join(path, childFileInfo.Name()))
		if nil != err {
			return nil, 0, err
		}
		fileObject[fmt.Sprintf("/%s", childFileInfo.Name())] = childFileObject
		totalSize += childSize
	}
	return fileObject, totalSize, nil
}
