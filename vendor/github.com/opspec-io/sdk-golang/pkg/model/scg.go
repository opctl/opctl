package model

// static call graph; see https://en.wikipedia.org/wiki/Call_graph
type Scg struct {
	Container *ScgContainerCall `yaml:"container,omitempty"`
	Op        *ScgOpCall        `yaml:"op,omitempty"`
	Parallel  []*Scg            `yaml:"parallel,omitempty"`
	Serial    []*Scg            `yaml:"serial,omitempty"`
}

type ScgContainerCall struct {
	// each entry of cmd will be interpolated
	Cmd  []string          `yaml:"cmd,omitempty"`
	Dirs map[string]string `yaml:"dirs,omitempty"`

	// each env var value will be interpolated
	EnvVars map[string]string      `yaml:"envVars,omitempty"`
	Files   map[string]string      `yaml:"files,omitempty"`
	Image   *ScgContainerCallImage `yaml:"image"`
	Sockets map[string]string      `yaml:"sockets,omitempty"`
	StdErr  map[string]string      `yaml:"stdErr,omitempty"`
	StdOut  map[string]string      `yaml:"stdOut,omitempty"`
	WorkDir string                 `yaml:"workDir,omitempty"`
}

type ScgContainerCallImage struct {
	// will be interpolated
	Ref string `yaml:"ref"`
	// will be interpolated
	PullIdentity string `yaml:"pullIdentity,omitempty"`
	// will be interpolated
	PullSecret string `yaml:"pullSecret,omitempty"`
}

type ScgOpCall struct {
	Ref string `yaml:"ref"`
	// binds in scope variables to inputs of referenced op
	Inputs map[string]string `yaml:"inputs,omitempty"`
	// binds in scope variables to outputs of referenced op
	Outputs map[string]string `yaml:"outputs,omitempty"`
}
