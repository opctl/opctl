package docker

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeFSPathConverter.go --fake-name fakeFSPathConverter ./ fsPathConverter

import (
	"github.com/opctl/opctl/sdks/go/internal/iruntime"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker/hostruntime"
	"github.com/pkg/errors"
	"path"
	"regexp"
	"strings"
)

type fsPathConverter interface {
	LocalToEngine(localPath string) string
}

func newFSPathConverter(cli hostruntime.ContainerInspector) (fsPathConverter, error) {
	hr, err := hostruntime.New(cli)
	if err != nil {
		return _fsPathConverter{}, errors.Wrap(err, "error detecting docker host runtime")
	}
	pc := _fsPathConverter{
		runtime:     iruntime.New(),
		hostRuntime: hr,
	}
	return pc, nil
}

type _fsPathConverter struct {
	runtime     iruntime.IRuntime
	hostRuntime hostruntime.RuntimeInfo
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

	if npc.runtime.GOOS() == "linux" && npc.hostRuntime.InAContainer {
		return npc.hostRuntime.HostPathMap.ToHostPath(localPath)
	}

	return localPath
}
