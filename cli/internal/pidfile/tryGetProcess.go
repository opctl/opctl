package pidfile

import (
	"context"
	"errors"
	"fmt"

	"github.com/shirou/gopsutil/v4/process"
)

func TryGetProcess(
	ctx context.Context,
	dirPath string,
) (*process.Process, error) {

	pID, err := getPIDFromFile(
		dirPath,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to read pid.lock: %w", err)
	}

	if pID <= 0 {
		// never started, already killed, or our mutex was manually removed
		return nil, nil
	}

	p, err := process.NewProcessWithContext(
		ctx,
		pID,
	)
	if err == nil {
		// running process
		return p, nil
	} else if errors.Is(err, process.ErrorProcessNotRunning) {
		// dead process
		return nil, nil
	}

	// unexpected retrieving process
	return nil, fmt.Errorf("unable to retrieve pid.lock process info: %w", err)
}
