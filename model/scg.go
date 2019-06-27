package model

// static call graph; see https://en.wikipedia.org/wiki/Call_graph
type SCG struct {
	Container    *SCGContainerCall `json:"container,omitempty"`
	If           *[]*SCGPredicate  `json:"if,omitempty"`
	Op           *SCGOpCall        `json:"op,omitempty"`
	Parallel     []*SCG            `json:"parallel,omitempty"`
	ParallelLoop *SCGParallelLoop  `json:"parallelLoop,omitempty"`
	Serial       []*SCG            `json:"serial,omitempty"`
	SerialLoop   *SCGSerialLoop    `json:"serialLoop,omitempty"`
}

type SCGContainerCall struct {
	// Cmd entries will be evaluated to strings
	Cmd []interface{} `json:"cmd,omitempty"`

	// Dirs entries will be evaluated to files
	Dirs map[string]string `json:"dirs,omitempty"`

	// EnvVars entries will be evaluated to strings
	EnvVars interface{} `json:"envVars,omitempty"`

	// Files entries will be evaluated to files
	Files   map[string]interface{} `json:"files,omitempty"`
	Image   *SCGContainerCallImage `json:"image"`
	Sockets map[string]string      `json:"sockets,omitempty"`
	StdErr  map[string]string      `json:"stdErr,omitempty"`
	StdOut  map[string]string      `json:"stdOut,omitempty"`
	WorkDir string                 `json:"workDir,omitempty"`
	Name    *string                `json:"name,omitempty"`
	Ports   map[string]string      `json:"ports,omitempty"`
}

type SCGContainerCallImage struct {
	// will be interpolated
	Ref       string        `json:"ref"`
	PullCreds *SCGPullCreds `json:"pullCreds,omitempty"`
}

type SCGLoopVars struct {
	Index *string `json:"index,omitempty"`
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type SCGOpCall struct {
	// Ref represents a references to the op; will be interpolated
	Ref string `json:"ref"`
	// PullCreds represent creds for pulling the op from a provider
	PullCreds *SCGPullCreds `json:"pullCreds,omitempty"`
	// binds scope to inputs of referenced op
	Inputs map[string]interface{} `json:"inputs,omitempty"`
	// binds scope to outputs of referenced op
	Outputs map[string]string `json:"outputs,omitempty"`
}

type SCGParallelLoop struct {
	Range interface{}  `json:"range,omitempty"`
	Run   SCG          `json:"run,omitempty"`
	Vars  *SCGLoopVars `json:"vars,omitempty"`
}

type SCGPredicate struct {
	Eq        *[]interface{} `json:"eq,omitempty"`
	Exists    *string        `json:"exists,omitempty"`
	Ne        *[]interface{} `json:"ne,omitempty"`
	NotExists *string        `json:"notExists,omitempty"`
}

type SCGPullCreds struct {
	// will be interpolated
	Username string `json:"username"`
	// will be interpolated
	Password string `json:"password"`
}

type SCGSerialLoop struct {
	Range interface{}     `json:"range,omitempty"`
	Run   SCG             `json:"run,omitempty"`
	Until []*SCGPredicate `json:"until,omitempty"`
	Vars  *SCGLoopVars    `json:"vars,omitempty"`
}
