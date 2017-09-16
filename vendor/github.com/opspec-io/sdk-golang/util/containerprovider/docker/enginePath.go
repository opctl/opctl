package docker

import (
	"path"
	"regexp"
	"strings"
)

var (
	windowsVolumeRegex = regexp.MustCompile(`([a-zA-Z]):(.*)`)
)

// converts a path on the local host to the path on the docker engine host using the conventions observed by
// Docker for Mac/Windows & Docker Machine
func (ctp _containerProvider) enginePath(localPath string) string {

	if ctp.runtime.GOOS() == "windows" {
		slashSeparatedPath := strings.Replace(localPath, `\`, `/`, -1)

		windowsVolumeRegexMatches := windowsVolumeRegex.FindStringSubmatch(slashSeparatedPath)
		if len(windowsVolumeRegexMatches) > 0 {
			return path.Join(`/`, strings.ToLower(windowsVolumeRegexMatches[1]), windowsVolumeRegexMatches[2])
		}
		return slashSeparatedPath
	}

	return localPath
}
