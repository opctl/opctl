package docker

import (
	"context"
	"net"

	"github.com/areYouLazy/libhosty"
)

func unregisterDNSName(
	ctx context.Context,
	ipAddress string,
) error {
	hfl, err := libhosty.Init()
	if err != nil {
		return err
	}

	hfl.RemoveHostsFileLinesByIP(
		net.ParseIP(ipAddress),
	)

	return hfl.SaveHostsFile()
}
