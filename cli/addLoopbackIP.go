package main

import (
	"errors"
	"net"
	"os/exec"
	"runtime"
)

// addLoopbackIP adds a new loopback IP using `ifconfig`
func addLoopbackIP(address string) error {

	ip, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "linux" {
		cmd = exec.Command("ifconfig", "lo", "add", ip, "netmask", "255.0.0.0")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("ifconfig", "lo0", "alias", ip, "255.0.0.0")
	} else {
		return errors.New("unsupported operating system")
	}

	return cmd.Run()
}

