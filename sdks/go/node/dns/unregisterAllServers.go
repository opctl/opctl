package dns

import (
	"os"
	"regexp"
)

func unregisterAllServers() error {
	rc, err := os.ReadFile(etcResolvConfPath)
	if err != nil {
		return err
	}

	rcString := string(rc)

	m1 := regexp.MustCompile(`.*managed by opctl.*\n`)

	return os.WriteFile(
		etcResolvConfPath,
		[]byte(
			m1.ReplaceAllString(rcString, ""),
		),
		0600,
	)
}
