package euid0

import (
	"os"
	"os/exec"
)

// Ensure EUID is 0 or return error
func Ensure() error {
	if os.Geteuid() != 0 {
		cmd := exec.Command("sudo", append([]string{"-E", os.Args[0]}, os.Args[1:]...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		defer os.Stdin.Close()

		err := cmd.Run()
		if err != nil {
			return err
		}

		// exit since we re-exec'd successfully
		os.Exit(0)
	}
	return nil
}
