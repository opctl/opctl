package model

// static call graph; see https://en.wikipedia.org/wiki/Call_graph
type Scg struct {
  Container *ContainerCall `yaml:"container,omitempty"`
  Op        *OpCall        `yaml:"op,omitempty"`
  Parallel  []*Scg   `yaml:"parallel,omitempty"`
  Serial    []*Scg   `yaml:"serial,omitempty"`
}

type ContainerCall struct {
  Cmd     []string             `yaml:"cmd,omitempty"`
  Env     []*ContainerCallEnvEntry `yaml:"env,omitempty"`
  Files   map[string]*ContainerFile `yaml:"files,omitempty"`
  Dirs    map[string]*ContainerDir `yaml:"dirs,omitempty"`
  Image   string               `yaml:"image"`
  Net     []*ContainerCallNetEntry `yaml:"net,omitempty"`
  WorkDir string               `yaml:"workDir,omitempty"`
}

// entry in a containers env; an env var
type ContainerCallEnvEntry struct {
  Bind string `yaml:"bind,omitempty"`
  // name of the env var in the container
  Name string `yaml:"name,omitempty"`
}

// file in a container
type ContainerFile struct {
  Bind string `yaml:"bind,omitempty"`
}

// dir in a container
type ContainerDir struct {
  Bind string `yaml:"bind,omitempty"`
}

// entry in a containers network; a network socket
type ContainerCallNetEntry struct {
  Bind        string `yaml:"bind,omitempty"`
  // aliases to give the network socket host in the container
  HostAliases []string `yaml:"hostAliases"`
}

type OpCall struct {
  Ref     string `yaml:"ref"`
  // binds in scope variables to inputs of referenced op
  Inputs  map[string]string `yaml:"inputs,omitempty"`
  // binds in scope variables to outputs of referenced op
  Outputs map[string]string `yaml:"outputs,omitempty"`
}
