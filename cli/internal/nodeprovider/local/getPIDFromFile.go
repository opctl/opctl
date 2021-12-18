package local

import (
	"os"
	"strconv"
)

func getPIDFromFile(
	pIDFilePath string,
) (int, error) {
	pIDFileBytes, err := os.ReadFile(pIDFilePath)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(pIDFileBytes))
}
