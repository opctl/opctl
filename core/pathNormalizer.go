package core

//go:generate counterfeiter -o ./fakePathNormalizer.go --fake-name fakePathNormalizer ./ pathNormalizer

import (
  "strings"
  "regexp"
  pathPkg "path"
  "os"
)

type pathNormalizer interface {
  Normalize(path string) string
}

func newPathNormalizer() pathNormalizer {
  return _pathNormalizer{
    driveRegex: regexp.MustCompile(`([a-zA-Z]):(.*)`),
  }
}

type _pathNormalizer struct {
  driveRegex *regexp.Regexp
}

func (this _pathNormalizer) Normalize(path string) string {
  backslashReplacedPath := strings.Replace(path, `\`, string(os.PathSeparator), -1)
  driveRegexMatches := this.driveRegex.FindStringSubmatch(backslashReplacedPath)
  if (len(driveRegexMatches) > 0) {
    return pathPkg.Join("/", strings.ToLower(driveRegexMatches[1]), driveRegexMatches[2])
  }
  return pathPkg.Clean(backslashReplacedPath)
}
