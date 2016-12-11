package model

type CallGraph struct {
  Container *ContainerCall `yaml:"container,omitempty"`
  Op        *OpCall `yaml:"op:omitempty"`
  Parallel  *ParallelCall `yaml:"parallel,omitempty"`
  Serial    *SerialCall `yaml:"serial,omitempty"`
}

type ContainerCall struct {
  Cmd     string `yaml:"cmd,omitempty"`
  Env     []*EnvEntry `yaml:"env,omitempty"`
  Fs      []*FsEntry `yaml:"fs,omitempty"`
  Image   string `yaml:"image"`
  Net     []*NetEntry `yaml:"net,omitempty"`
  WorkDir string `yaml:"workDir,omitempty"`
}

// entry in an env; an env var
type EnvEntry struct {
  Bind string `yaml:"bind,omitempty"`
  Name string `yaml:"name,omitempty"`
}

// entry in a fs; a file/directory
type FsEntry struct {
  Bind string `yaml:"bind,omitempty"`
  Path string `yaml:"path"`
}

// entry in a network; a host
type NetEntry struct {
  Bind string `yaml:"bind,omitempty"`
  Port int `yaml:"port,omitempty"`
}

type OpCall struct {
  Ref     string `yaml:ref`
  // binds inputs of referenced op to in scope variables
  Args    *map[string]string `yaml:args,omitempty`
  // binds in scope variables to outputs of referenced op
  Results *map[string]string `yaml:results:omitempty`
}

type ParallelCall []CallGraph

type SerialCall []CallGraph
