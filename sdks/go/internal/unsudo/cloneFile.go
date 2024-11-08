package unsudo

import (
	"io"
	"os"
	"path/filepath"
)

// CloneFile as the user & group who ran sudo
func CloneFile(
	srcPath,
	dstPath string,
) error {
	if err := CreateDir(
		filepath.Dir(dstPath),
	); err != nil {
		return err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return os.Chown(dstPath, getSudoUID(), getSudoGID())
}
