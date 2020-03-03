package model

// dynamic call graph; see https://en.wikipedia.org/wiki/Call_graph
type DCG struct {
	Container *DCGContainerCall `json:"container,omitempty"`
	// id of call
	Id           string               `json:"id"`
	If           *bool                `json:"if,omitempty"`
	IsKilled     bool                 `json:"isKilled"`
	Op           *DCGOpCall           `json:"op,omitempty"`
	Parallel     []*SCG               `json:"parallel,omitempty"`
	ParallelLoop *DCGParallelLoopCall `json:"parallelLoop,omitempty"`
	// id of parent call
	ParentID   *string            `json:"parentId,omitempty"`
	Serial     []*SCG             `json:"serial,omitempty"`
	SerialLoop *DCGSerialLoopCall `json:"serialLoop,omitempty"`
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
	Src       *Value     `json:"src,omitempty"`
	Ref       *string    `json:"ref"`
	PullCreds *PullCreds `json:"pullCreds,omitempty"`
}

type DCGLoopVars struct {
	Index *string `json:"index,omitempty"`
	Key   *string `json:"key,omitempty"`
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

type DCGParallelLoopCall struct {
	// an array or object
	Range *Value       `json:"range,omitempty"`
	Run   DCG          `json:"run,omitempty"`
	Vars  *DCGLoopVars `json:"vars,omitempty"`
}

type DCGPredicate struct {
	Eq []*Value `json:"eq"`
	Ne []*Value `json:"ne"`
}

type DCGSerialLoopCall struct {
	// an array or object
	Range *Value       `json:"range,omitempty"`
	Run   DCG          `json:"run,omitempty"`
	Until *bool        `json:"until,omitempty"`
	Vars  *DCGLoopVars `json:"vars,omitempty"`
}
