package resolvercfg

import (
	"context"
	"os"
	"regexp"
)

// Delete modifications to the current system
func Delete(
	ctx context.Context,
) error {
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
