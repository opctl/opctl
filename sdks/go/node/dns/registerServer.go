package dns

import (
	"fmt"
	"os"
	"strings"
)

var etcResolvConfPath = "/etc/resolv.conf"

func registerServer(
	ipAddress string,
) error {
	rc, err := os.ReadFile(etcResolvConfPath)
	if err != nil {
		return err
	}

	nsPrefix := fmt.Sprintf("nameserver %s", ipAddress)

	rcString := string(rc)
	if strings.HasPrefix(rcString, nsPrefix) {
		return os.WriteFile(
			etcResolvConfPath,
			[]byte(fmt.Sprintf("%s # do not edit; managed by opctl\n%s", nsPrefix, rcString)),
			0600,
		)
	}
	return nil
}
