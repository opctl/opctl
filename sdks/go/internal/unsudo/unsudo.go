package unsudo

import "os"

func getSudoUID() int {
	uid := os.Geteuid()
	if sudoUID := tryGetEnvInt("SUDO_UID"); sudoUID != nil {
		uid = *sudoUID
	}
	return uid
}

func getSudoGID() int {
	gid := os.Getegid()
	if sudoGID := tryGetEnvInt("SUDO_GID"); sudoGID != nil {
		gid = *sudoGID
	}
	return gid
}
