package docker

//go:generate counterfeiter -o ./fakeFSPathConverter.go --fake-name fakeFSPathConverter ./ fsPathConverter

import (
	"github.com/opctl/sdk-golang/util/iruntime"
	"path"
	"regexp"
	"strings"
)

type fsPathConverter interface {
	LocalToEngine(localPath string) string
}

func newFSPathConverter() fsPathConverter {
	return _fsPathConverter{
		runtime: iruntime.New(),
	}
}

type _fsPathConverter struct {
	runtime iruntime.IRuntime
}

var (
	windowsVolumeRegex = regexp.MustCompile(`([a-zA-Z]):(.*)`)
)

// converts a path on the local host to the path on the docker engine host using the conventions observed by
// Docker for Mac/Windows & Docker Machine
func (npc _fsPathConverter) LocalToEngine(localPath string) string {

	if npc.runtime.GOOS() == "windows" {
		slashSeparatedPath := strings.Replace(localPath, `\`, `/`, -1)

		windowsVolumeRegexMatches := windowsVolumeRegex.FindStringSubmatch(slashSeparatedPath)
		if len(windowsVolumeRegexMatches) > 0 {
			return path.Join(`/`, strings.ToLower(windowsVolumeRegexMatches[1]), windowsVolumeRegexMatches[2])
		}
		return slashSeparatedPath
	}

	return localPath
}
