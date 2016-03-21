package models

type PipelineView struct {
  Description string
  Name        string
  Stages      []PipelineStageView
}

func NewPipelineView(
description string,
name string,
stages []PipelineStageView,
) *PipelineView {

  return &PipelineView{
    Description:description,
    Name:name,
    Stages:stages,
  }

}