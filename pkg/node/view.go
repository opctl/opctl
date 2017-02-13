package node

type CapabilityView struct {
	LinuxContainerProvider bool
}

type InfoView struct {
	// will be formatted as multi-addresses (see https://github.com/multiformats/multiaddr)
	Addresses    []string
	Capabilities *CapabilityView
}
