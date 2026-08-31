package pidfile

import (
	"os"
	"strconv"
	"strings"
)

// getPIDFromFile returns the PID recorded in dirPath's pid file, or 0 if
// there's no usable PID.
//
// A file which is missing, empty, or garbled tells us nothing about a running
// node, so it's reported as 0 rather than as an error; we read the file
// unlocked, so catching a holder mid-write is expected.
func getPIDFromFile(
	dirPath string,
) (int32, error) {
	pIDFileBytes, err := os.ReadFile(
		constructPIDFilePath(dirPath),
	)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}

		return 0, err
	}

	pidInt, err := strconv.Atoi(
		strings.TrimSpace(string(pIDFileBytes)),
	)
	if err != nil {
		return 0, nil
	}

	return int32(pidInt), nil
}
