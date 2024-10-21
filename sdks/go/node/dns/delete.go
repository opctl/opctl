package dns

func Delete() error {
	if err := unregisterAllServers(); err != nil {
		return err
	}

	return deleteAllResolvers()
}
