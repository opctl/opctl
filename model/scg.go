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
	// Cmd entries will be evaluated to strings
	Cmd []interface{} `yaml:"cmd,omitempty"`

	// Dirs entries will be evaluated to files
	Dirs map[string]string `yaml:"dirs,omitempty"`

	// EnvVars entries will be evaluated to strings
	EnvVars map[string]interface{} `yaml:"envVars,omitempty"`

	// Files entries will be evaluated to files
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
	// Ref represents a references to the op; will be interpolated
	Ref string `yaml:"ref"`
	// PullCreds represent creds for pulling the op from a provider
	PullCreds *SCGPullCreds `yaml:"pullCreds,omitempty"`
	// binds scope to inputs of referenced op
	Inputs map[string]interface{} `yaml:"inputs,omitempty"`
	// binds scope to outputs of referenced op
	Outputs map[string]string `yaml:"outputs,omitempty"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface to handle deprecated properties gracefully in one place
func (soc *SCGOpCall) UnmarshalYAML(
	unmarshal func(interface{}) error,
) error {

	if err := unmarshal(soc); nil != err {
		return err
	}

	// handle deprecated property
	deprecated := struct {
		Pkg *SCGOpCallPkg `yaml:"pkg"`
	}{}
	if err := unmarshal(&deprecated); nil != err {
		return err
	}

	if nil != deprecated.Pkg {
		soc.Ref = deprecated.Pkg.Ref
		soc.PullCreds = deprecated.Pkg.PullCreds
	}
	return nil
}

// SCGOpCallPkg deprecated
type SCGOpCallPkg struct {
	// will be interpolated
	Ref       string        `yaml:"ref"`
	PullCreds *SCGPullCreds `yaml:"pullCreds,omitempty"`
}
