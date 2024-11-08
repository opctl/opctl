package main

import (
	"fmt"
	"os"
)

// ensureEuid0 or return error
func ensureEuid0() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("re-run command with sudo")
	}
	return nil
}
