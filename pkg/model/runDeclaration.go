package model

type RunDeclaration struct {
  Op OpRunDeclaration `yaml:"op,omitempty"`
  Parallel *ParallelRunDeclaration `yaml:"parallel,omitempty"`
  Serial   *SerialRunDeclaration `yaml:"serial,omitempty"`
}

type OpRunDeclaration string

type ParallelRunDeclaration []RunDeclaration

type SerialRunDeclaration []RunDeclaration
