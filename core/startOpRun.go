package core

import (
  "github.com/opspec-io/sdk-golang/models"
  "time"
)

func (this _api) StartOpRun(
req models.StartOpRunReq,
) (
opRunId string,
err error,
) {

  opRunId = this.uniqueStringFactory.Construct()
  rootOpRunId := opRunId // this is root

  go func() {
    err = this.opRunner.Run(
      opRunId,
      req.Args,
      this.pathNormalizer.Normalize(req.OpUrl),
      rootOpRunId,
    )
    if (nil != err) {
      this.eventStream.Publish(
        models.Event{
          Timestamp: time.Now().UTC(),
          OpRunEncounteredError: &models.OpRunEncounteredErrorEvent{
            Msg: err.Error(),
            OpRunId:opRunId,
            RootOpRunId:rootOpRunId,
          },
        },
      )
    }
  }()

  return

}
