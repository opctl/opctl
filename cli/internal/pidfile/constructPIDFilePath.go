package pidfile

import "path/filepath"

func constructPIDFilePath(
	dirPath string,
) string {
	return filepath.Join(
		dirPath,
		"pid.lock",
	)
}
