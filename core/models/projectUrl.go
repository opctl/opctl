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

  var url *url.URL
  url, err = url.Parse(projectUrlStr)

  if (url.IsAbs() && strings.ToLower(url.Scheme) != "file") {
    err = errors.New(
      fmt.Sprintf(
        "Remote projects are currently unsupported. Received non local projectUrl: `%v`",
        projectUrlStr,
      ),
    )
  }

  projectUrl = &ProjectUrl{url}

  return

}

type ProjectUrl struct {
  *url.URL
}