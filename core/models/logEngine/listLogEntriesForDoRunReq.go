package logengine

func NewListLogEntriesForDoRunReq(
doRunId string,
) *ListLogEntriesForDoRunReq {

  return &ListLogEntriesForDoRunReq{
    DoRunId :doRunId,
  }

}

type ListLogEntriesForDoRunReq struct {
  DoRunId string
}