package pathnormalizer

import (
  "strings"
  "regexp"
  pathPkg "path"
  "os"
)

type PathNormalizer interface {
  Normalize(path string) string
}

func NewPathNormalizer() PathNormalizer {
  return pathNormalizer{
    driveRegex: regexp.MustCompile(`([a-zA-Z]):(.*)`),
  }
}

type pathNormalizer struct {
  driveRegex *regexp.Regexp
}

func (this pathNormalizer) Normalize(path string) string {
  backslashReplacedPath := strings.Replace(path, `\`, string(os.PathSeparator), -1)
  driveRegexMatches := this.driveRegex.FindStringSubmatch(backslashReplacedPath)
  if (len(driveRegexMatches) > 0) {
    return pathPkg.Join("/", strings.ToLower(driveRegexMatches[1]), driveRegexMatches[2])
  }
  return pathPkg.Clean(backslashReplacedPath)
}
