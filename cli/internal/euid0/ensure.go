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

// sudo-rs ignores -E, so vars are named instead, minus those classic sudo -E resets.
// Named vars override sudo-rs's own, so stale SUDO_* identity vars are dropped too.
func sudoArgs(sudoRs bool, environ []string, pathToOpctlBin string, args []string) []string {
	preserveEnv := "-E"
	if sudoRs {
		var names []string
		for _, kv := range environ {
			name, value, _ := strings.Cut(kv, "=")
			switch name {
			case "", "PATH", "USER", "LOGNAME",
				"SUDO_USER", "SUDO_UID", "SUDO_GID", "SUDO_HOME", "SUDO_COMMAND":
				continue
			}
			// --preserve-env= is comma-separated
			if strings.Contains(name, ",") {
				continue
			}
			// classic sudo -E drops these (exported bash functions)
			if strings.HasPrefix(value, "()") {
				continue
			}
			names = append(names, name)
		}
		preserveEnv = "--preserve-env=" + strings.Join(names, ",")
	}

	return append([]string{preserveEnv, pathToOpctlBin}, args...)
}
