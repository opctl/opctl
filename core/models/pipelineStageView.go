package models

type PipelineStageView struct {
  Name string
  Type string
}

func NewPipelineStageView(
name string,
_type string,
) *PipelineStageView {

  return &PipelineStageView{
    Name:name,
    Type:_type,
  }

}