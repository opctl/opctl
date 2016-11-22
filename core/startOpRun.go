package core

import (
  "github.com/opspec-io/sdk-golang/pkg/model"
)

func (this _core) StartOpRun(
req model.StartOpRunReq,
) (
opRunId string,
err error,
) {

  opRunId = this.uniqueStringFactory.Construct()

  go func() {
    err = this.opRunner.Run(
      opRunId,
      req.Args,
      this.pathNormalizer.Normalize(req.OpUrl),
      "",
    )
  }()

  return

}
