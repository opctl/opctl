package model

type Container struct {
  Cmd     string `yaml:"cmd,omitempty"`
  Env     []*ContainerEnvEntry `yaml:"env,omitempty"`
  Fs      []*ContainerFsEntry `yaml:"fs,omitempty"`
  Image   string `yaml:"image"`
  Net     []*ContainerNetEntry `yaml:"net,omitempty"`
  WorkDir string `yaml:"workDir,omitempty"`
}

type ContainerEnvEntry struct {
  Bind string `yaml:"bind,omitempty"`
  Name string `yaml:"name,omitempty"`
}

type ContainerFsEntry struct {
  Bind string `yaml:"bind,omitempty"`
  Path string `yaml:"path"`
}

type ContainerNetEntry struct {
  Bind string `yaml:"bind,omitempty"`
  Port int `yaml:"port,omitempty"`
}
