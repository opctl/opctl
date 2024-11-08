package pidfile

import (
	"context"
	"fmt"
	"os"
)

func TryGetLock(
	ctx context.Context,
	dirPath string,
) (bool, error) {
	file, err := os.OpenFile(
		constructPIDFilePath(dirPath),
		os.O_CREATE|os.O_EXCL|os.O_WRONLY,
		0744,
	)
	if err != nil {
		// unexpected error
		if !os.IsExist(err) {
			return false, err
		}

		nodeProcess, err := TryGetProcess(ctx, dirPath)
		if err != nil {
			return false, err
		}

		if nodeProcess != nil {
			return false, nil
		}

		// process doesn't exist anymore; we can overwrite...
	}
	defer file.Close()

	if err := os.WriteFile(
		constructPIDFilePath(dirPath),
		[]byte(fmt.Sprintf("%d", os.Getpid())),
		0744,
	); err != nil {
		return false, err
	}

	return true, nil
}
