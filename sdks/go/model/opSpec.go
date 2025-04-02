package model

// OpSpec is a spec for an op
type OpSpec struct {
	Description string                `json:"description"`
	Inputs      map[string]*ParamSpec `json:"inputs,omitempty"`
	Name        string                `json:"name"`
	Outputs     map[string]*ParamSpec `json:"outputs,omitempty"`
	Run         *CallSpec             `json:"run,omitempty"`
	Version     string                `json:"version,omitempty"`
}

// CallSpec is a spec for a node of a call graph; see https://en.wikipedia.org/wiki/Call_graph
type CallSpec struct {
	Container    *ContainerCallSpec    `json:"container,omitempty"`
	Description  string                `json:"description,omitempty"`
	If           *[]*PredicateSpec     `json:"if,omitempty"`
	Name         *string               `json:"name,omitempty"`
	Needs        []string              `json:"needs,omitempty"`
	Op           *OpCallSpec           `json:"op,omitempty"`
	Parallel     *[]*CallSpec          `json:"parallel,omitempty"`
	ParallelLoop *ParallelLoopCallSpec `json:"parallelLoop,omitempty"`
	Serial       *[]*CallSpec          `json:"serial,omitempty"`
	SerialLoop   *SerialLoopCallSpec   `json:"serialLoop,omitempty"`
}

// ContainerCallSpec is a spec for calling a container
type ContainerCallSpec struct {
	// Cmd entries will be interpreted to strings
	Cmd interface{} `json:"cmd,omitempty"`
	// Dirs entries will be interpreted to dirs
	Dirs map[string]interface{} `json:"dirs,omitempty"`

	// EnvVars entries will be interpreted to strings
	EnvVars interface{} `json:"envVars,omitempty"`
	// Dirs entries will be interpreted to files
	Files   map[string]interface{}  `json:"files,omitempty"`
	Image   *ContainerCallImageSpec `json:"image"`
	Sockets map[string]string       `json:"sockets,omitempty"`
	WorkDir string                  `json:"workDir,omitempty"`
	Name    *string                 `json:"name,omitempty"`
	Ports   map[string]string       `json:"ports,omitempty"`
}

// ContainerCallImageSpec is a spec for the image when calling a container
type ContainerCallImageSpec struct {
	Platform  *OCIImagePlatformSpec `json:"platform,omitempty"`
	PullCreds *CredsSpec            `json:"pullCreds,omitempty"`
	Ref       string                `json:"ref"`
}

// LoopVarsSpec is a spec for a loops vars
type LoopVarsSpec struct {
	Index *string `json:"index,omitempty"`
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

// OpCallSpec is a spec for calling an op
type OpCallSpec struct {
	// Ref represents a references to an op; will be interpolated
	Ref string `json:"ref"`
	// PullCreds represent creds for pulling the op from a provider
	PullCreds *CredsSpec `json:"pullCreds,omitempty"`
	// binds scope to inputs of referenced op
	Inputs map[string]interface{} `json:"inputs,omitempty"`
	// binds scope to outputs of referenced op
	Outputs map[string]string `json:"outputs,omitempty"`
}

type OCIImagePlatformSpec struct {
	Arch string `json:"arch,omitempty"`
}

// ParallelLoopCallSpec is a spec for calling a parallel loop
type ParallelLoopCallSpec struct {
	Range interface{}   `json:"range,omitempty"`
	Run   CallSpec      `json:"run,omitempty"`
	Vars  *LoopVarsSpec `json:"vars,omitempty"`
}

// PredicateSpec is a spec for a predicate
type PredicateSpec struct {
	Eq        *[]interface{} `json:"eq,omitempty"`
	Exists    *string        `json:"exists,omitempty"`
	LT        *[]interface{} `json:"lt,omitempty"`
	LTE       *[]interface{} `json:"lte,omitempty"`
	GT        *[]interface{} `json:"gt,omitempty"`
	GTE       *[]interface{} `json:"gte,omitempty"`
	NE        *[]interface{} `json:"ne,omitempty"`
	NotExists *string        `json:"notExists,omitempty"`
}

// CredsSpec is a spec for authentication credentials
type CredsSpec struct {
	// will be interpolated
	Username string `json:"username"`
	// will be interpolated
	Password string `json:"password"`
}

// SerialLoopCallSpec is a spec for calling a serial loop
type SerialLoopCallSpec struct {
	Range interface{}      `json:"range,omitempty"`
	Run   CallSpec         `json:"run,omitempty"`
	Until []*PredicateSpec `json:"until,omitempty"`
	Vars  *LoopVarsSpec    `json:"vars,omitempty"`
}

type ReferenceOpts struct {
	Type string
	// for creating dirs/files
	ScratchDir string
}
