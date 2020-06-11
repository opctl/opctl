package hostruntime

import (
	"github.com/docker/docker/api/types/mount"
	"strings"
)

type HostPathMap map[string]string

func newHostPathMap(mounts []mount.Mount) HostPathMap {
	mappings := make(map[string]string)

	for _, m := range mounts {
		mappings[m.Source] = m.Target
	}

	return mappings
}

func (m HostPathMap) ToHostPath(localPath string) string {
	if len(m) <= 0 {
		return localPath
	}

	for host, local := range m {
		if strings.HasPrefix(localPath, local) {
			mapped := strings.Replace(localPath, local, host, 1)
			return mapped
		}
	}

	return localPath
}
