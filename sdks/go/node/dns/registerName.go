package dns

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func RegisterName(
	name,
	ipAddress string,
) error {

	ips := ipsByHostname[name]
	if ips == nil {
		ipsByHostname[name] = map[string]interface{}{}
	}

	ipsByHostname[name][ipAddress] = nil

	if runtime.GOOS == "darwin" {
		nameParts := strings.Split(name, ".")
		err := setResolver(nameParts[len(nameParts)-1])
		if err != nil {
			return err
		}

		cmd := exec.Command("dscacheutil", "-flushcache")

		outputBytes, err := cmd.CombinedOutput()

		if err != nil {
			return fmt.Errorf("failed to flush dns cache: %w, %s", err, string(outputBytes))
		}

		cmd = exec.Command("killall", "-HUP", "mDNSResponder")
		outputBytes, err = cmd.CombinedOutput()

		if err != nil {
			return fmt.Errorf("failed to flush mdns cache: %w, %s", err, string(outputBytes))
		}
	}

	return nil
}
