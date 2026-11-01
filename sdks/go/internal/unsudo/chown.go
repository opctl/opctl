package unsudo

import "os"

// Chown hands fsPath back to the user & group who ran sudo
func Chown(
	fsPath string,
) error {
	return os.Chown(fsPath, getSudoUID(), getSudoGID())
}
