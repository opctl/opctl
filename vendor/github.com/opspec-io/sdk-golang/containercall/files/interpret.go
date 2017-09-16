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
	for scgContainerFilePath, fileExpression := range scgContainerCallFiles {

		if "" == fileExpression {
			// bound implicitly
			fileExpression = scgContainerFilePath
		}

		isBoundToPkgContent := strings.HasPrefix(fileExpression, "/")
		value, isBoundToScope := scope[fileExpression]

		var contentSrc io.Reader
		contentFileMode := os.FileMode(0666)
		var err error
		switch {
		case isBoundToPkgContent:
			// bound to pkg file
			pkgContentReadSeekCloser, err := pkgHandle.GetContent(context.TODO(), fileExpression)
			if nil != err {
				return nil, fmt.Errorf(
					"unable to bind file '%v' to pkg content '%v'; error was: %v",
					scgContainerFilePath,
					fileExpression,
					err.Error(),
				)
			}
			defer pkgContentReadSeekCloser.Close()
			contentSrc = pkgContentReadSeekCloser

			pkgContentsList, err := pkgHandle.ListContents(context.TODO())
			if nil != err {
				return nil, fmt.Errorf(
					"unable to bind file '%v' to pkg content '%v'; error was: %v",
					scgContainerFilePath,
					fileExpression,
					err.Error(),
				)
			}

			// @TODO: return mode from GetContent so this isn't needed
			for _, pkgContent := range pkgContentsList {
				if pkgContent.Path == fileExpression {
					contentFileMode = pkgContent.Mode
					break
				}
			}
		case isBoundToScope:
			if nil == value {
				return nil, fmt.Errorf(
					"unable to bind file '%v' to '%v'; '%v' null",
					scgContainerFilePath,
					fileExpression,
					fileExpression,
				)
			} else if nil != value.File {
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
						fileExpression,
						err.Error(),
					)
				}
			} else {
				content, err := f.data.CoerceToString(value)
				if nil != err {
					return nil, fmt.Errorf(
						"unable to bind file '%v' to '%v'; error was: %v",
						scgContainerFilePath,
						fileExpression,
						err.Error(),
					)
				}
				contentSrc = strings.NewReader(*content.String)
			}
		default:
			// unbound
			contentSrc = strings.NewReader("")
		}
		dcgContainerCallFiles[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)

		// create file
		if err := f.os.MkdirAll(
			path.Dir(dcgContainerCallFiles[scgContainerFilePath]),
			0777,
		); nil != err {
			return nil, fmt.Errorf(
				"unable to mkdir for bound file '%v'; error was: %v",
				scgContainerFilePath,
				err.Error(),
			)
		}

		outputFile, err := f.os.OpenFile(dcgContainerCallFiles[scgContainerFilePath], os.O_RDWR|os.O_CREATE, contentFileMode)
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
