package pidfile

import (
	"os"
	"strconv"
)

func getPIDFromFile(
	dirPath string,
) (int32, error) {
	pIDFileBytes, err := os.ReadFile(
		constructPIDFilePath(dirPath),
	)
	if err != nil {
		return 0, err
	}

	pidInt, err := strconv.Atoi(string(pIDFileBytes))
	if err != nil {
		return 0, err
	}

	return int32(pidInt), nil
}
