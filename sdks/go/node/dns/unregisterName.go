package dns

func UnregisterName(
	name,
	ipAddress string,
) error {
	ipsByHostnameMu.Lock()
	defer ipsByHostnameMu.Unlock()

	ips, ok := ipsByHostname[name]
	if !ok {
		return nil
	}

	delete(ips, ipAddress)

	return nil

}
