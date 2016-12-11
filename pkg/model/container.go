package model

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
