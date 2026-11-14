package euid0

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Ensure EUID is 0 or return error
func Ensure() error {
	if os.Geteuid() != 0 {
		pathToOpctlBin, err := os.Executable()
		if err != nil {
			return err
		}

		pathToOpctlBin, err = filepath.EvalSymlinks(pathToOpctlBin)
		if err != nil {
			return err
		}

		cmd := exec.Command("sudo", sudoArgs(isSudoRs(), os.Environ(), pathToOpctlBin, os.Args[1:])...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		defer os.Stdin.Close()

		err = cmd.Run()
		if err != nil {
			return err
		}

		// exit since we re-exec'd successfully
		os.Exit(0)
	}
	return nil
}

func isSudoRs() bool {
	out, err := exec.Command("sudo", "--version").Output()
	return err == nil && bytes.HasPrefix(out, []byte("sudo-rs"))
}

// sudoArgs preserves the caller's environment across the sudo re-exec.
// sudo-rs (default sudo on Ubuntu 26.04+) ignores -E, so there each variable
// is named via --preserve-env=, minus those classic sudo -E resets anyway.
func sudoArgs(sudoRs bool, environ []string, pathToOpctlBin string, args []string) []string {
	preserveEnv := "-E"
	if sudoRs {
		var names []string
		for _, kv := range environ {
			name, _, _ := strings.Cut(kv, "=")
			switch name {
			case "", "PATH", "USER", "LOGNAME":
				continue
			}
			// --preserve-env= is comma separated, so such a name can't be expressed
			if strings.Contains(name, ",") {
				continue
			}
			names = append(names, name)
		}
		preserveEnv = "--preserve-env=" + strings.Join(names, ",")
	}

	return append([]string{preserveEnv, pathToOpctlBin}, args...)
}
