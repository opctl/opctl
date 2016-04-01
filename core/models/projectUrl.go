package models

import (
  "net/url"
  "strings"
  "errors"
  "fmt"
)

func NewProjectUrl(
projectUrlStr string,
) (projectUrl *ProjectUrl, err error) {

  var parsedUrl *url.URL
  parsedUrl, err = url.Parse(projectUrlStr)
  if (nil != err) {
    return
  }

  projectUrl = &ProjectUrl{parsedUrl}

  if (parsedUrl.IsAbs() && strings.ToLower(parsedUrl.Scheme) != "file") {
    err = errors.New(
      fmt.Sprintf(
        "Remote projects are currently unsupported. Received non local projectUrl: `%v`",
        projectUrlStr,
      ),
    )
  }

  return

}

type ProjectUrl struct {
  *url.URL
}