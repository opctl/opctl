package model

type CallGraphDeclaration struct {
  Container *ContainerCallDeclaration `yaml:"container,omitempty"`
  Op        *OpCallDeclaration `yaml:"op:omitempty"`
  Parallel  *ParallelCallDeclaration `yaml:"parallel,omitempty"`
  Serial    *SerialCallDeclaration `yaml:"serial,omitempty"`
}

type ContainerCallDeclaration struct {
  Cmd     string `yaml:"cmd,omitempty"`
  Env     []*EnvEntryDeclaration `yaml:"env,omitempty"`
  Fs      []*FsEntryDeclaration `yaml:"fs,omitempty"`
  Image   string `yaml:"image"`
  Net     []*NetEntryDeclaration `yaml:"net,omitempty"`
  WorkDir string `yaml:"workDir,omitempty"`
}

// declaration of an entry in an env; an env var
type EnvEntryDeclaration struct {
  Bind string `yaml:"bind,omitempty"`
  Name string `yaml:"name,omitempty"`
}

// declaration of an entry in a fs; a file/directory
type FsEntryDeclaration struct {
  Bind string `yaml:"bind,omitempty"`
  Path string `yaml:"path"`
}

// declaration of an entry in a network (a.k.a a host)
type NetEntryDeclaration struct {
  Bind string `yaml:"bind,omitempty"`
  Port int `yaml:"port,omitempty"`
}

type OpCallDeclaration struct {
  Ref     string `yaml:ref`
  // binds inputs of referenced op to in scope variables
  Args    *map[string]string `yaml:args,omitempty`
  // binds in scope variables to outputs of referenced op
  Results *map[string]string `yaml:results:omitempty`
}

type ParallelCallDeclaration []CallGraphDeclaration

type SerialCallDeclaration []CallGraphDeclaration
