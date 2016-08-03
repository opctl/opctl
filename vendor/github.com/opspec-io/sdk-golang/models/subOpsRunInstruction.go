package models

func NewSubOpsRunInstruction(
subOps []SubOpRunInstruction,
) *SubOpsRunInstruction {

  return &SubOpsRunInstruction{
    SubOps:subOps,
  }

}

type SubOpsRunInstruction struct {
  SubOps []SubOpRunInstruction `yaml:"subOps"`
}
