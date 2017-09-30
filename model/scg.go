package model

// static call graph; see https://en.wikipedia.org/wiki/Call_graph
type SCG struct {
	Container *SCGContainerCall `yaml:"container,omitempty"`
	Op        *SCGOpCall        `yaml:"op,omitempty"`
	Parallel  []*SCG            `yaml:"parallel,omitempty"`
	Serial    []*SCG            `yaml:"serial,omitempty"`
}

type SCGPullCreds struct {
	// will be interpolated
	Username string `yaml:"username"`
	// will be interpolated
	Password string `yaml:"password"`
}

type SCGContainerCall struct {
	// Cmd will be evaluated to strings
	Cmd []string `yaml:"cmd,omitempty"`
	// Dirs will be evaluated to files
	Dirs map[string]string `yaml:"dirs,omitempty"`

	// EnvVars will be evaluated to strings
	EnvVars map[string]interface{} `yaml:"envVars,omitempty"`

	// Files will be evaluated to files
	Files   map[string]interface{} `yaml:"files,omitempty"`
	Image   *SCGContainerCallImage `yaml:"image"`
	Sockets map[string]string      `yaml:"sockets,omitempty"`
	StdErr  map[string]string      `yaml:"stdErr,omitempty"`
	StdOut  map[string]string      `yaml:"stdOut,omitempty"`
	WorkDir string                 `yaml:"workDir,omitempty"`
	Name    string                 `yaml:"name,omitempty"`
	Ports   map[string]string      `yaml:"ports,omitempty"`
}

type SCGContainerCallImage struct {
	// will be interpolated
	Ref       string        `yaml:"ref"`
	PullCreds *SCGPullCreds `yaml:"pullCreds,omitempty"`
}

type SCGOpCall struct {
	Pkg *SCGOpCallPkg `yaml:"pkg"`
	// binds scope to inputs of referenced op
	Inputs map[string]interface{} `yaml:"inputs,omitempty"`
	// binds scope to outputs of referenced op
	Outputs map[string]string `yaml:"outputs,omitempty"`
}

type SCGOpCallPkg struct {
	// will be interpolated
	Ref       string        `yaml:"ref"`
	PullCreds *SCGPullCreds `yaml:"pullCreds,omitempty"`
}
