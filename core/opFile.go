package core

type opFileSubOp struct {
  Url string `yaml:"url"`
}

type opFile struct {
  Name        *string `yaml:"name"`
  Description string `yaml:"description"`
  SubOps      []opFileSubOp `yaml:"subOps,omitempty"`
}