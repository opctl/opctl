package model

import (
	"sort"

	"github.com/ipld/go-ipld-prime"
)

// Auth holds auth data
type Auth struct {
	// Resources designates which resources this auth applies to in the form of a reference (or prefix thereof)
	Resources string
	Creds     Creds
}

// Call is a node of a call graph; see https://en.wikipedia.org/wiki/Call_graph
type Call struct {
	Container *ContainerCall `json:"container,omitempty"`
	// id of call
	ID           string            `json:"id"`
	If           *bool             `json:"if,omitempty"`
	IsKilled     bool              `json:"isKilled"`
	Name         *string           `json:"name,omitempty"`
	Needs        []string          `json:"needs,omitempty"`
	Op           *OpCall           `json:"op,omitempty"`
	Parallel     []*CallSpec       `json:"parallel,omitempty"`
	ParallelLoop *ParallelLoopCall `json:"parallelLoop,omitempty"`
	// id of parent call
	ParentID *string `json:"parentId,omitempty"`
	// id of root call
	RootID     string          `json:"rootId"`
	Serial     []*CallSpec     `json:"serial,omitempty"`
	SerialLoop *SerialLoopCall `json:"serialLoop,omitempty"`
}

type BaseCall struct {
	OpPath string `json:"opPath"`
}

// ContainerCall is a call of a container
type ContainerCall struct {
	BaseCall
	ContainerID string   `json:"containerId"`
	Cmd         []string `json:"cmd"`
	// format: containerPath => hostPath
	Dirs StringMap `json:"dirs"`
	// format: name => value
	EnvVars StringMap `json:"envVars"`
	// format: containerPath => hostPath
	Files StringMap           `json:"files"`
	Image *ContainerCallImage `json:"image"`
	// format: containerSocket => hostSocket
	Sockets StringMap `json:"sockets"`
	WorkDir string    `json:"workDir"`
	Name    *string   `json:"name,omitempty"`
	Ports   StringMap `json:"ports,omitempty"`
}

func NewStringMap(input map[string]string) StringMap {
	var keys []string

	if input != nil {
		// Extract keys from the map and sort them
		keys = make([]string, 0, len(input))
		for key := range input {
			keys = append(keys, key)
		}
	}

	// Sort keys for deterministic ordering
	sort.Strings(keys)

	return StringMap{
		Keys:   keys,
		Values: input,
	}
}

type StringMap struct {
	Keys   []string
	Values map[string]string
}

// ContainerCallImage is the image used when calling a container
type ContainerCallImage struct {
	//	Src       *Value  `json:"src,omitempty"`
	Ref       *string `json:"ref"`
	PullCreds *Creds  `json:"pullCreds,omitempty"`
}

// Creds contains authentication credentials
type Creds struct {
	Username string
	Password string
}

// LoopVars is a loops vars
type LoopVars struct {
	Index *string `json:"index,omitempty"`
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

// OpCall is a call of an op
type OpCall struct {
	BaseCall
	OpID              string            `json:"opId"`
	Inputs            map[string]*ipld.Node `json:"inputs"`
	ChildCallCallSpec *CallSpec         `json:"childCallScg"`
	ChildCallID       string            `json:"childCallId"`
}

// ParallelLoopCall is a call of a parallel loop
type ParallelLoopCall struct {
	// an array or object
	Range *ipld.Node    `json:"range,omitempty"`
	Run   Call      `json:"run,omitempty"`
	Vars  *LoopVars `json:"vars,omitempty"`
}

// SerialLoopCall is a call of a serial loop
type SerialLoopCall struct {
	// an array or object
	Range *ipld.Node    `json:"range,omitempty"`
	Run   Call      `json:"run,omitempty"`
	Until *bool     `json:"until,omitempty"`
	Vars  *LoopVars `json:"vars,omitempty"`
}
