package model

// static call graph; see https://en.wikipedia.org/wiki/Call_graph
type Scg struct {
	Container *ScgContainerCall `yaml:"container,omitempty"`
	Op        *ScgOpCall        `yaml:"op,omitempty"`
	Parallel  []*Scg            `yaml:"parallel,omitempty"`
	Serial    []*Scg            `yaml:"serial,omitempty"`
}

type ScgContainerCall struct {
	// cmd strings will be interpolated
	Cmd     []string                       `yaml:"cmd,omitempty"`
	EnvVars map[string]*ScgContainerEnvVar `yaml:"envVars,omitempty"`
	Files   map[string]*ScgContainerFile   `yaml:"files,omitempty"`
	Dirs    map[string]*ScgContainerDir    `yaml:"dirs,omitempty"`
	Image   string                         `yaml:"image"`
	Sockets map[string]*ScgContainerSocket `yaml:"sockets,omitempty"`
	WorkDir string                         `yaml:"workDir,omitempty"`
}

// entry in a containers env; an env var
type ScgContainerEnvVar struct {
	Bind string `yaml:"bind,omitempty"`
	// value string will be interpolated
	Value string `yaml:"value,omitempty"`
}

// file in a container
type ScgContainerFile struct {
	Bind string `yaml:"bind,omitempty"`
}

// dir in a container
type ScgContainerDir struct {
	Bind string `yaml:"bind,omitempty"`
}

// socket in a container
type ScgContainerSocket struct {
	Bind string `yaml:"bind,omitempty"`
}

type ScgOpCall struct {
	Ref string `yaml:"ref"`
	// binds in scope variables to inputs of referenced op
	Inputs map[string]string `yaml:"inputs,omitempty"`
	// binds in scope variables to outputs of referenced op
	Outputs map[string]string `yaml:"outputs,omitempty"`
}
