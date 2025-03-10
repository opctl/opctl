package dns

import (
	"github.com/goodhosts/hostsfile"
)

func UnregisterName(
	name,
	ipAddress string,
) error {
	h, err := hostsfile.NewHosts()
	if err != nil {
		return err
	}

	err = h.Remove(
		ipAddress,
		name,
	)
	if err != nil {
		return err
	}

	return h.Flush()
}
