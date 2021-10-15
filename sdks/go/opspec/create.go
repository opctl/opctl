package opspec

import (
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// Create an operation
func Create(
	path,
	pkgName,
	pkgDescription string,
) error {

	err := os.MkdirAll(
		path,
		0777,
	)
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

	return os.WriteFile(
		filepath.Join(path, opfile.FileName),
		opFileBytes,
		0777,
	)

}
