package model

// dynamic call graph; see https://en.wikipedia.org/wiki/Call_graph
type Dcg struct {
	Id        string        `json:"id"`
	Container *DcgContainer `json:"container,omitempty"`
	Op        *DcgOp        `json:"op,omitempty"`
	OpGraphId string        `json:"opGraphId"`
	OpRef     string        `json:"opRef"`
}

type DcgContainer struct{}

type DcgContainerCall struct {
	Cmd []string `json:"cmd"`
	// format: containerPath => hostPath
	Dirs map[string]string `json:"dirs"`
	// format: name => value
	EnvVars map[string]string `json:"envVars"`
	// format: containerPath => hostPath
	Files map[string]string `json:"files"`
	Image string            `json:"image"`
	// format: containerSocket => hostSocket
	Sockets   map[string]string `json:"sockets"`
	WorkDir   string            `json:"workDir"`
	IpAddress string            `json:"ipAddress"`
}

type DcgOp struct{}
