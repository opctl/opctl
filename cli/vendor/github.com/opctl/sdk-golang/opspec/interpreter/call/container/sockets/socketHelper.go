package sockets

import "strings"

func isUnixSocketAddress(address string) bool {
	const unixSocketAddressDiscriminationChars = `/\`
	// note: this mechanism for determining the type of socket is naive; higher level of sophistication may be required
	return strings.ContainsAny(address, unixSocketAddressDiscriminationChars)
}
