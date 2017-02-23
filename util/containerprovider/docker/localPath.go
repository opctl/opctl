package docker

import (
	"fmt"
	"path"
	"regexp"
)

// converts a path on the docker engine host to the path on the local host using the conventions observed by
// Docker for Mac/Windows & Docker Machine
func (this _containerProvider) localPath(enginePath string) string {

	if this.runtime.GOOS() == "windows" {
		driveRegex := regexp.MustCompile(`/([a-zA-Z])/(.*)`)
		driveRegexMatches := driveRegex.FindStringSubmatch(enginePath)

		if len(driveRegexMatches) > 0 {
			return path.Join(fmt.Sprintf("%v:", driveRegexMatches[1]), driveRegexMatches[2])
		}

	}

	return enginePath
}
