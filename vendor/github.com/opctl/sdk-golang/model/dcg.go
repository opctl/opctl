package model

// dynamic call graph; see https://en.wikipedia.org/wiki/Call_graph
type DCG struct {
	Container *DCGContainer `json:"container,omitempty"`
	Op        *DCGOp        `json:"op,omitempty"`
}

type DCGBaseCall struct {
	RootOpID string `json:"rootOpId"`
	OpHandle DataHandle
}

type DCGContainer struct{}

type DCGContainerCall struct {
	DCGBaseCall
	ContainerID string   `json:"containerId"`
	Cmd         []string `json:"cmd"`
	// format: containerPath => hostPath
	Dirs map[string]string `json:"dirs"`
	// format: name => value
	EnvVars map[string]string `json:"envVars"`
	// format: containerPath => hostPath
	Files map[string]string      `json:"files"`
	Image *DCGContainerCallImage `json:"image"`
	// format: containerSocket => hostSocket
	Sockets map[string]string `json:"sockets"`
	WorkDir string            `json:"workDir"`
	Name    string            `json:"name,omitempty"`
	Ports   map[string]string `json:"ports,omitempty"`
}

type DCGContainerCallImage struct {
	Ref       string     `yaml:"ref"`
	PullCreds *PullCreds `yaml:"pullCreds,omitempty"`
}

type DCGOp struct{}

type DCGOpCall struct {
	DCGBaseCall
	OpID         string            `json:"opId"`
	Inputs       map[string]*Value `json:"inputs"`
	ChildCallSCG *SCG              `json:"childCallScg"`
	ChildCallID  string            `json:"childCallId"`
}

type DCGOpCallPkg struct {
	Ref       string     `yaml:"ref"`
	PullCreds *PullCreds `yaml:"pullCreds,omitempty"`
}
