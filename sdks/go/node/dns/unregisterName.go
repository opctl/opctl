package dns

func UnregisterName(
	name,
	ipAddress string,
) error {
	ips, ok := ipsByHostname[name]
	if !ok {
		return nil
	}

	delete(ips, ipAddress)

	return nil

}
