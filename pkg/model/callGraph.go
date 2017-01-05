package model

type CallGraph struct {
	Container *ContainerCall `yaml:"container,omitempty"`
	Op        *OpCall        `yaml:"op,omitempty"`
	Parallel  []*CallGraph   `yaml:"parallel,omitempty"`
	Serial    []*CallGraph   `yaml:"serial,omitempty"`
}

type ContainerCall struct {
	Cmd     []string             `yaml:"cmd,omitempty"`
	Env     []*ContainerEnvEntry `yaml:"env,omitempty"`
	Fs      []*ContainerFsEntry  `yaml:"fs,omitempty"`
	Image   string               `yaml:"image"`
	Net     []*ContainerNetEntry `yaml:"net,omitempty"`
	WorkDir string               `yaml:"workDir,omitempty"`
}

type OpCall struct {
	Ref string `yaml:"ref"`
	// binds in scope variables to inputs of referenced op
	Inputs map[string]string `yaml:"inputs,omitempty"`
	// binds in scope variables to outputs of referenced op
	Outputs map[string]string `yaml:"outputs,omitempty"`
}
