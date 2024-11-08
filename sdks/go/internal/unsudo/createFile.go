package unsudo

import (
	"os"
	"path/filepath"
)

// CreateFile as the user & group who ran sudo
func CreateFile(
	fsPath string,
	data []byte,
) error {
	if err := CreateDir(filepath.Dir(fsPath)); err != nil {
		return err
	}

	if err := os.WriteFile(fsPath, []byte{}, 0700); err != nil {
		return err
	}

	return os.Chown(fsPath, getSudoUID(), getSudoGID())
}
