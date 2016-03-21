package ports
import "github.com/dev-op-spec/engine/core/models"


type ContainerEngine interface {
  InitDevOp(
  devOpName string,
  ) (err error)

  RunDevOp(
  devOpName string,
  ) (devOpRun models.DevOpRunView, err error)
}
