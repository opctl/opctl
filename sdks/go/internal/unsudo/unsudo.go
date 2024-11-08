package unsudo

func getSudoUID() int {
	uid := 0
	if sudoUID := tryGetEnvInt("SUDO_UID"); sudoUID != nil {
		uid = *sudoUID
	}
	return uid
}

func getSudoGID() int {
	gid := 0
	if sudoGID := tryGetEnvInt("SUDO_GID"); sudoGID != nil {
		gid = *sudoGID
	}
	return gid
}
