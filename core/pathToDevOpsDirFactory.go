package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToDevOpsDirFactory interface {
  Construct(
  projectUrl *models.ProjectUrl,
  ) (pathToDevOpsDir string)
}

func newPathToDevOpsDirFactory(
) (pathToDevOpsDirFactory pathToDevOpsDirFactory) {

  return &_pathToDevOpsDirFactory{}

}

type _pathToDevOpsDirFactory struct{}

func (this _pathToDevOpsDirFactory) Construct(
projectUrl *models.ProjectUrl,
) (pathToDevOpsDir string) {

  pathToDevOpsDir = path.Join(
    projectUrl.Path,
    ".dev-op-spec",
    "dev-ops",
  )

  return

}
