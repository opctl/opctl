package opspath

import (
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
)

func GetLocal() (
	string,
	error,
) {
	var (
		opspecDirInfo,
		opctlDirInfo os.FileInfo
		err error
	)
	if opspecDirInfo, err = os.Stat(
		model.DotOpspecDirName,
	); err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if opctlDirInfo, err = os.Stat(
		model.DotOpctlDirName,
	); err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if opctlDirInfo == nil && opspecDirInfo != nil {
		// support deprecated .opspec dir until we remove support
		return filepath.Abs(opspecDirInfo.Name())
	}

	return filepath.Abs(model.DotOpctlDirName)
}
