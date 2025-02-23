package dns

func UnregisterNames(
	names []string,
	ipAddress string,
) error {
	for _, name := range names {
		ips, ok := ipsByHostname[name]
		if !ok {
			return nil
		}

		delete(ips, ipAddress)
	}

	return nil

}
