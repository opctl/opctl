package hostruntime

import (
	"strings"
)

type HostPathMap map[string]string

func newHostPathMap(binds []string) HostPathMap {
	mappings := make(map[string]string)

	for _, b := range binds {
		tokens := strings.Split(b, ":")
		mappings[tokens[0]] = tokens[1]
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
