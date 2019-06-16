package model

// dynamic call graph; see https://en.wikipedia.org/wiki/Call_graph
type DCG struct {
	Container *DCGContainerCall `json:"container,omitempty"`
	// id of call
	Id       string     `json:"id"`
	If       *bool      `json:"if,omitempty"`
	IsKilled bool       `json:"isKilled"`
	Loop     *DCGLoop   `json:"loop,omitempty"`
	Op       *DCGOpCall `json:"op,omitempty"`
	Parallel []*SCG     `json:"parallel,omitempty"`
	// id of parent call
	ParentID *string `json:"parentId,omitempty"`
	Serial   []*SCG  `json:"serial,omitempty"`
}

type DCGBaseCall struct {
	RootOpID string `json:"rootOpId"`
	OpHandle DataHandle
}

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
	Name    *string           `json:"name,omitempty"`
	Ports   map[string]string `json:"ports,omitempty"`
}

type DCGContainerCallImage struct {
	Ref       string     `json:"ref"`
	PullCreds *PullCreds `json:"pullCreds,omitempty"`
}

type DCGLoop struct {
	For   *DCGLoopFor `json:"for,omitempty"`
	Index *string     `json:"index,omitempty"`
	Until *bool       `json:"until,omitempty"`
}

type DCGLoopFor struct {
	// an array or object
	Each  *Value  `json:"each"`
	Value *string `json:"value,omitempty"`
}

type DCGOpCall struct {
	DCGBaseCall
	OpID         string            `json:"opId"`
	Inputs       map[string]*Value `json:"inputs"`
	ChildCallSCG *SCG              `json:"childCallScg"`
	ChildCallID  string            `json:"childCallId"`
}

type DCGOpCallPkg struct {
	Ref       string     `json:"ref"`
	PullCreds *PullCreds `json:"pullCreds,omitempty"`
}

type DCGPredicate struct {
	Eq []*Value `json:"eq"`
	Ne []*Value `json:"ne"`
}
