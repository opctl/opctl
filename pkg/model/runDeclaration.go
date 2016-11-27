package model

type RunDeclaration struct {
  Container *ContainerRunDeclaration `yaml:"container,omitempty"`
  Op        *OpRunDeclaration `yaml:"op:omitempty"`
  Parallel  *ParallelRunDeclaration `yaml:"parallel,omitempty"`
  Serial    *SerialRunDeclaration `yaml:"serial,omitempty"`
}

type ContainerRunDeclaration Container

type OpRunDeclaration struct {
  Ref     string `yaml:ref`
  // binds inputs of referenced op to in scope variables
  Args    *map[string]string `yaml:args,omitempty`
  // binds in scope variables to outputs of referenced op
  Results *map[string]string `yaml:results:omitempty`
}

type ParallelRunDeclaration []RunDeclaration

type SerialRunDeclaration []RunDeclaration
