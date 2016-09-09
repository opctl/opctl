package models

type RunDeclaration struct {
  Op OpRunDeclaration `yaml:"op,omitempty"`
  Parallel *ParallelRunDeclaration `yaml:"parallel,omitempty"`
  Serial   *SerialRunDeclaration `yaml:"serial,omitempty"`
}
