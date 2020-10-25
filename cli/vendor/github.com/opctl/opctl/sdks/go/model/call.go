package model

//Call is a node of a call graph; see https://en.wikipedia.org/wiki/Call_graph
type Call struct {
	Container *ContainerCall `json:"container,omitempty"`
	// id of call
	Id           string            `json:"id"`
	If           *bool             `json:"if,omitempty"`
	IsKilled     bool              `json:"isKilled"`
	Name         *string           `json:"name,omitempty"`
	Needs        []string          `json:"needs,omitempty"`
	Op           *OpCall           `json:"op,omitempty"`
	Parallel     []*CallSpec       `json:"parallel,omitempty"`
	ParallelLoop *ParallelLoopCall `json:"parallelLoop,omitempty"`
	// id of parent call
	ParentID   *string         `json:"parentId,omitempty"`
	Serial     []*CallSpec     `json:"serial,omitempty"`
	SerialLoop *SerialLoopCall `json:"serialLoop,omitempty"`
}

type BaseCall struct {
	RootOpID string `json:"rootOpId"`
	OpPath   string `json:"opPath"`
}

//ContainerCall is a call of a container
type ContainerCall struct {
	BaseCall
	ContainerID string   `json:"containerId"`
	Cmd         []string `json:"cmd"`
	// format: containerPath => hostPath
	Dirs map[string]string `json:"dirs"`
	// format: name => value
	EnvVars map[string]string `json:"envVars"`
	// format: containerPath => hostPath
	Files map[string]string   `json:"files"`
	Image *ContainerCallImage `json:"image"`
	// format: containerSocket => hostSocket
	Sockets map[string]string `json:"sockets"`
	WorkDir string            `json:"workDir"`
	Name    *string           `json:"name,omitempty"`
	Ports   map[string]string `json:"ports,omitempty"`
}

//ContainerCallImage is the image used when calling a container
type ContainerCallImage struct {
	Src *Value `json:"src,omitempty"`
	// @TODO: deprecate in favor of Src
	Ref       *string    `json:"ref"`
	PullCreds *PullCreds `json:"pullCreds,omitempty"`
}

//LoopVars is a loops vars
type LoopVars struct {
	Index *string `json:"index,omitempty"`
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

//OpCall is a call of an op
type OpCall struct {
	BaseCall
	OpID              string            `json:"opId"`
	Inputs            map[string]*Value `json:"inputs"`
	ChildCallCallSpec *CallSpec         `json:"childCallScg"`
	ChildCallID       string            `json:"childCallId"`
}

//ParallelLoopCall is a call of a parallel loop
type ParallelLoopCall struct {
	// an array or object
	Range *Value    `json:"range,omitempty"`
	Run   Call      `json:"run,omitempty"`
	Vars  *LoopVars `json:"vars,omitempty"`
}

//Predicate is a predicate i.e. something that evaluates to true or false
type Predicate struct {
	Eq []*Value `json:"eq"`
	Ne []*Value `json:"ne"`
}

//SerialLoopCall is a call of a serial loop
type SerialLoopCall struct {
	// an array or object
	Range *Value    `json:"range,omitempty"`
	Run   Call      `json:"run,omitempty"`
	Until *bool     `json:"until,omitempty"`
	Vars  *LoopVars `json:"vars,omitempty"`
}
