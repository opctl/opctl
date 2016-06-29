package models

type OpFileSubOp struct {
  Url string `yaml:"url"`
  IsParallel bool `yaml:"isParallel"`
}

type OpFileParam struct {
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
}

type OpFile struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  Params      map[string]OpFileParam `yaml:"params,omitempty"`
  SubOps      []OpFileSubOp `yaml:"subOps,omitempty"`
}
