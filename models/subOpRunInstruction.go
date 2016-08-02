package models

func NewSubOpRunInstruction(
url string,
isParallel bool,
) *SubOpRunInstruction {

  return &SubOpRunInstruction{
    Url:url,
    IsParallel:isParallel,
  }

}

type SubOpRunInstruction struct {
  Url        string `yaml:"url"`
  IsParallel bool `yaml:"isParallel,omitempty"`
}
