package model

type CallGraph struct {
  Container *ContainerCall `yaml:"container,omitempty"`
  Op        *OpCall `yaml:"op:omitempty"`
  Parallel  *ParallelCall `yaml:"parallel,omitempty"`
  Serial    *SerialCall `yaml:"serial,omitempty"`
}

type ContainerCall struct {
  Cmd     []string `yaml:"cmd,omitempty"`
  Env     []*ContainerEnvEntry `yaml:"env,omitempty"`
  Fs      []*ContainerFsEntry `yaml:"fs,omitempty"`
  Image   string `yaml:"image"`
  Net     []*ContainerNetEntry `yaml:"net,omitempty"`
  WorkDir string `yaml:"workDir,omitempty"`
}

type OpCall struct {
  Ref     string `yaml:"ref"`
  // binds inputs of referenced op to in scope variables
  Args    map[string]string `yaml:"args,omitempty"`
  // binds in scope variables to outputs of referenced op
  Results map[string]string `yaml:"results,omitempty"`
}

type ParallelCall []CallGraph

type SerialCall []CallGraph
