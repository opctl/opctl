package docker

import (
	"context"
	"github.com/areYouLazy/libhosty"
)

func registerDNSName(
	ctx context.Context,
	dnsName,
	ip string,
) error {
	hfl, err := libhosty.Init()
	if err != nil {
		return err
	}

	hfl.AddHostsFileLine(ip, dnsName, "managed by opctl")

	return hfl.SaveHostsFile()
}
