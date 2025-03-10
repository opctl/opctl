package dns

import (
	"context"
	"fmt"
	"github.com/goodhosts/hostsfile"
)

const managedByOpctlComment = "managed by opctl"

func RegisterName(
	ctx context.Context,
	name,
	ipAddress string,
) error {
	h, err := hostsfile.NewHosts()
	if err != nil {
		return err
	}

	err = h.AddRaw(
		fmt.Sprintf(
			"%s %s # %s",
			ipAddress,
			name,
			managedByOpctlComment,
		),
	)
	if err != nil {
		return err
	}

	return h.Flush()
}
