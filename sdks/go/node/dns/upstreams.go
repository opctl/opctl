package dns

import (
	"net"
	"os"
	"strings"
)

var resolvConfPath = "/etc/resolv.conf"

/*
snapshotUpstreams parses nameserver entries from resolvConfPath into "host:port"
upstream addresses.

It excludes the node's own listen address (a wildcard listen address matches any
loopback nameserver) and entries tagged "managed by opctl" (leftover from a
previous node on the same host) to prevent forwarding loops. An unreadable
resolv.conf yields no upstreams.
*/
func snapshotUpstreams(
	resolvConfPath,
	ownIPAddress,
	ownPort string,
) []string {
	rc, err := os.ReadFile(resolvConfPath)
	if err != nil {
		return nil
	}

	ownIP := net.ParseIP(ownIPAddress)

	upstreams := []string{}
	for _, line := range strings.Split(string(rc), "\n") {
		if strings.Contains(line, "managed by opctl") {
			continue
		}

		if commentStart := strings.IndexAny(line, "#;"); commentStart != -1 {
			line = line[:commentStart]
		}

		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] != "nameserver" {
			continue
		}

		ip := net.ParseIP(fields[1])
		if ip == nil {
			continue
		}

		// resolv.conf nameservers always speak on port 53
		if ownPort == "53" &&
			(ip.Equal(ownIP) || (ownIP != nil && ownIP.IsUnspecified() && ip.IsLoopback())) {
			continue
		}

		upstreams = append(
			upstreams,
			net.JoinHostPort(fields[1], "53"),
		)
	}

	return upstreams
}
