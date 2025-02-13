package dns

import (
	"context"
	"runtime"

	"strconv"

	"github.com/shirou/gopsutil/v4/net"
)

func shouldReusePort(
	ctx context.Context,
	port string,
) (bool, error) {
	if runtime.GOOS != "darwin" {
		return false, nil
	}

	portUInt, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return false, err
	}

	// Get all net connections
	conns, err := net.ConnectionsWithContext(
		ctx,
		"all",
	)
	if err != nil {
		return false, err
	}

	// Allow port reuse if processes bound to *:<port> (e.g. Apples MDNSResponder)
	for _, conn := range conns {
		if conn.Laddr.IP == "*" && conn.Laddr.Port == uint32(portUInt) && conn.Status == "LISTEN" {
			return true, nil
		}
	}
	return false, nil
}
