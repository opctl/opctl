package files

import (
	"context"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (f _Files) Interpret(
	pkgHandle model.PkgHandle,
	scope map[string]*model.Value,
	scgContainerCallFiles map[string]string,
	scratchDirPath string,
) (map[string]string, error) {
	dcgContainerCallFiles := map[string]string{}
fileLoop:
	for scgContainerFilePath, scgContainerFileBind := range scgContainerCallFiles {

		if "" == scgContainerFileBind {
			// bound implicitly
			scgContainerFileBind = scgContainerFilePath
		}

		isBoundToPkgContent := strings.HasPrefix(scgContainerFileBind, "/")
		value, isBoundToScope := scope[scgContainerFileBind]

		var contentSrc io.Reader
		var err error
		switch {
		case isBoundToPkgContent:
			// bound to pkg file
			dcgContainerCallFiles[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)

			pkgContentReadSeekCloser, err := pkgHandle.GetContent(context.TODO(), scgContainerFileBind)
			defer pkgContentReadSeekCloser.Close()
			if nil != err {
				return nil, fmt.Errorf(
					"unable to bind file '%v' to pkg content '%v'; error was: %v",
					scgContainerFilePath,
					scgContainerFileBind,
					err.Error(),
				)
			}

			pkgContentsList, err := pkgHandle.ListContents(context.TODO())
			if nil != err {
				return nil, fmt.Errorf(
					"unable to bind file '%v' to pkg content '%v'; error was: %v",
					scgContainerFilePath,
					scgContainerFileBind,
					err.Error(),
				)
			}

			// @TODO: return mode from GetContent so this isn't needed
			var contentMode os.FileMode
			for _, pkgContent := range pkgContentsList {
				if pkgContent.Path == scgContainerFileBind {
					contentMode = pkgContent.Mode
					break
				}
			}

			containerFileWriter, err := f.os.Create(dcgContainerCallFiles[scgContainerFilePath])
			defer containerFileWriter.Close()
			if nil != err {
				return nil, fmt.Errorf(
					"unable to bind file '%v' to pkg content '%v'; error was: %v",
					scgContainerFilePath,
					scgContainerFileBind,
					err.Error(),
				)
			}

			err = f.os.Chmod(dcgContainerCallFiles[scgContainerFilePath], contentMode)
			if nil != err {
				return nil, fmt.Errorf(
					"unable to bind file '%v' to pkg content '%v'; error was: %v",
					scgContainerFilePath,
					scgContainerFileBind,
					err.Error(),
				)
			}

			_, err = f.io.Copy(containerFileWriter, pkgContentReadSeekCloser)
			if nil != err {
				return nil, fmt.Errorf(
					"unable to bind file '%v' to pkg content '%v'; error was: %v",
					scgContainerFilePath,
					scgContainerFileBind,
					err.Error(),
				)
			}

			containerFileWriter.Close()
			pkgContentReadSeekCloser.Close()

			continue fileLoop
		case isBoundToScope:
			switch {
			case nil == value:
				return nil, fmt.Errorf(
					"unable to bind file '%v' to '%v'; '%v' null",
					scgContainerFilePath,
					scgContainerFileBind,
					scgContainerFileBind,
				)
			case nil != value.File:
				if !strings.HasPrefix(*value.File, f.rootFSPath) {
					// bound to non rootFS file
					dcgContainerCallFiles[scgContainerFilePath] = *value.File
					continue fileLoop
				}

				// bound to rootFS file
				contentSrc, err = f.os.Open(*value.File)
				if nil != err {
					return nil, fmt.Errorf(
						"unable to bind file '%v' to '%v'; error was: %v",
						scgContainerFilePath,
						scgContainerFileBind,
						err.Error(),
					)
				}
			default:
				content, err := f.data.CoerceToString(value)
				if nil != err {
					return nil, fmt.Errorf(
						"unable to bind file '%v' to '%v'; error was: %v",
						scgContainerFilePath,
						scgContainerFileBind,
						err.Error(),
					)
				}
				contentSrc = strings.NewReader(content)
			}
		default:
			// unbound
			contentSrc = strings.NewReader("")
		}
		dcgContainerCallFiles[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)

		// create file
		if err := f.os.MkdirAll(
			path.Dir(dcgContainerCallFiles[scgContainerFilePath]),
			0700,
		); nil != err {
			return nil, err
		}
		outputFile, err := f.os.Create(dcgContainerCallFiles[scgContainerFilePath])
		if nil != err {
			return nil, err
		}

		// copy content to file
		_, err = f.io.Copy(outputFile, contentSrc)
		outputFile.Close()
		if nil != err {
			return nil, err
		}

	}
	return dcgContainerCallFiles, nil
}
