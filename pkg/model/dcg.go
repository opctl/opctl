package model

// dynamic call graph; see https://en.wikipedia.org/wiki/Call_graph
type Dcg struct {
	Container *DcgContainer `json:"container,omitempty"`
	Op        *DcgOp        `json:"op,omitempty"`
}

type DcgContainer struct {
}

type DcgContainerCall struct {
	ContainerId string   `json:"containerId"`
	Cmd         []string `json:"cmd"`
	// format: containerPath => hostPath
	Dirs map[string]string `json:"dirs"`
	// format: name => value
	EnvVars map[string]string `json:"envVars"`
	// format: containerPath => hostPath
	Files     map[string]string      `json:"files"`
	Image     *DcgContainerCallImage `json:"image"`
	IpAddress string                 `json:"ipAddress"`
	RootOpId  string                 `json:"rootOpId"`
	PkgRef    string                 `json:"pkgRef"`
	// format: containerSocket => hostSocket
	Sockets map[string]string `json:"sockets"`
	WorkDir string            `json:"workDir"`
}

type DcgContainerCallImage struct {
	Ref          string `json:"ref"`
	PullIdentity string `json:"pullIdentity,omitempty"`
	PullSecret   string `json:"pullSecret,omitempty"`
}

type DcgOp struct{}
