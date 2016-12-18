package model

// entry in a containers env; an env var
type ContainerEnvEntry struct {
  Bind string `yaml:"bind,omitempty"`
  // name of the env var in the container
  Name string `yaml:"name,omitempty"`
}

// entry in a containers fs; a file/directory
type ContainerFsEntry struct {
  Bind string `yaml:"bind,omitempty"`
  // path of the file/directory in the container
  Path string `yaml:"path"`
}

// entry in a containers network; a network socket
type ContainerNetEntry struct {
  Bind        string `yaml:"bind,omitempty"`
  // aliases to give the network socket host in the container
  HostAliases string `yaml:"hostAliases"`
}
