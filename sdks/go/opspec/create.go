package opspec

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// Create an operation
func Create(
	path,
	pkgName,
	pkgDescription string,
) error {

	err := mkdirAll(path)
	if err != nil {
		return err
	}

	opFile := model.OpSpec{
		Description: pkgDescription,
		Name:        pkgName,
	}

	opFileBytes, err := yaml.Marshal(&opFile)
	if err != nil {
		return err
	}

	opFilePath := filepath.Join(path, opfile.FileName)
	if err := os.WriteFile(
		opFilePath,
		opFileBytes,
		0777,
	); err != nil {
		return err
	}

	return unsudo.Chown(opFilePath)
}

// mkdirAll hands each dir it creates to the user who ran sudo, because the CLI re-runs itself under sudo to start a node.
func mkdirAll(path string) error {
	var created []string
	for dir := filepath.Clean(path); ; dir = filepath.Dir(dir) {
		if _, err := os.Lstat(dir); err == nil {
			break
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		created = append(created, dir)
		if filepath.Dir(dir) == dir {
			break
		}
	}
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}
	for _, dir := range created {
		if err := unsudo.Chown(dir); err != nil {
			return err
		}
	}
	return nil
}
