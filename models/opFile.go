package models

type OpFileSubOp struct {
  Url string `yaml:"url"`
}

type OpFileParam struct {
  IsSecret    bool `yaml:"isSecret"`
  Description string `yaml:"description"`
}

type OpFile struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  Params      map[string]OpFileParam `yaml:"params,omitempty"`
  SubOps      []OpFileSubOp `yaml:"subOps,omitempty"`
}
