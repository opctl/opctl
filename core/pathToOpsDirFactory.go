package core

import (
  "path"
  "github.com/open-devops/engine/core/models"
)

type pathToOpsDirFactory interface {
  Construct(
  projectUrl *models.Url,
  ) (pathToOpsDir string)
}

func newPathToOpsDirFactory(
) (pathToOpsDirFactory pathToOpsDirFactory) {

  return &_pathToOpsDirFactory{}

}

type _pathToOpsDirFactory struct{}

func (this _pathToOpsDirFactory) Construct(
projectUrl *models.Url,
) (pathToOpsDir string) {

  pathToOpsDir = path.Join(
    projectUrl.Path,
    ".open-devops",
    "ops",
  )

  return

}
