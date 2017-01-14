package model

// static call graph; see https://en.wikipedia.org/wiki/Call_graph
type Scg struct {
  Container *ScgContainer `yaml:"container,omitempty"`
  Op        *ScgOp        `yaml:"op,omitempty"`
  Parallel  []*Scg   `yaml:"parallel,omitempty"`
  Serial    []*Scg   `yaml:"serial,omitempty"`
}

type ScgContainer struct {
  Cmd     []string             `yaml:"cmd,omitempty"`
  Env     []*ScgContainerEnvEntry `yaml:"env,omitempty"`
  Files   map[string]*ScgContainerFile `yaml:"files,omitempty"`
  Dirs    map[string]*ScgContainerDir `yaml:"dirs,omitempty"`
  Image   string               `yaml:"image"`
  Net     []*ScgContainerNetEntry `yaml:"net,omitempty"`
  WorkDir string               `yaml:"workDir,omitempty"`
}

// entry in a containers env; an env var
type ScgContainerEnvEntry struct {
  Bind string `yaml:"bind,omitempty"`
  // name of the env var in the container
  Name string `yaml:"name,omitempty"`
}

// file in a container
type ScgContainerFile struct {
  Bind string `yaml:"bind,omitempty"`
}

// dir in a container
type ScgContainerDir struct {
  Bind string `yaml:"bind,omitempty"`
}

// entry in a containers network; a network socket
type ScgContainerNetEntry struct {
  Bind        string `yaml:"bind,omitempty"`
  // aliases to give the network socket host in the container
  HostAliases []string `yaml:"hostAliases"`
}

type ScgOp struct {
  Ref     string `yaml:"ref"`
  // binds in scope variables to inputs of referenced op
  Inputs  map[string]string `yaml:"inputs,omitempty"`
  // binds in scope variables to outputs of referenced op
  Outputs map[string]string `yaml:"outputs,omitempty"`
}
