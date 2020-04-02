package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
)

// NamedSCG represents a named SCG (as opposed to an anonymous SCG)
type NamedSCG map[string]SCG

func (s *NamedSCG) UnmarshalJSON(b []byte) error {
	rawSCGMap := make(map[string]json.RawMessage)
	err := yaml.Unmarshal(b, &rawSCGMap)
	if nil != err {
		return err
	}

	var scgBytesParts []string
	for key, value := range rawSCGMap {
		// if known call type found, consider anonymous and convert it to named.
		switch key {
		case "container":
			scgBytesParts = append(scgBytesParts, fmt.Sprintf("\"%v\": %s", key, value))
		case "if":
			scgBytesParts = append(scgBytesParts, fmt.Sprintf("\"%v\": %s", key, value))
		case "needs":
			scgBytesParts = append(scgBytesParts, fmt.Sprintf("\"%v\": %s", key, value))
		case "op":
			scgBytesParts = append(scgBytesParts, fmt.Sprintf("\"%v\": %s", key, value))
		case "parallel":
			scgBytesParts = append(scgBytesParts, fmt.Sprintf("\"%v\": %s", key, value))
		case "parallelLoop":
			scgBytesParts = append(scgBytesParts, fmt.Sprintf("\"%v\": %s", key, value))
		case "serial":
			scgBytesParts = append(scgBytesParts, fmt.Sprintf("\"%v\": %s", key, value))
		case "serialLoop":
			scgBytesParts = append(scgBytesParts, fmt.Sprintf("\"%v\": %s", key, value))
		}
	}

	var scgBytes []byte
	scgName := ""
	if nil != scgBytesParts {
		scgBytes = []byte(fmt.Sprintf("{%v}", strings.Join(scgBytesParts, ",")))
	} else {
		for key, value := range rawSCGMap {
			scgName = key
			scgBytes = value
			// named calls only have one map entry i.e. their name.
			break
		}
	}

	scg := SCG{}
	err = yaml.Unmarshal(
		scgBytes,
		&scg,
	)

	*s = map[string]SCG{
		scgName: scg,
	}

	return err
}

//SCG models a static call graph; see https://en.wikipedia.org/wiki/Call_graph
type SCG struct {
	Container    *SCGContainerCall    `json:"container,omitempty"`
	If           *[]*SCGPredicate     `json:"if,omitempty"`
	Needs        []string             `json:"needs,omitempty"`
	Op           *SCGOpCall           `json:"op,omitempty"`
	Parallel     []NamedSCG           `json:"parallel,omitempty"`
	ParallelLoop *SCGParallelLoopCall `json:"parallelLoop,omitempty"`
	Serial       []NamedSCG           `json:"serial,omitempty"`
	SerialLoop   *SCGSerialLoopCall   `json:"serialLoop,omitempty"`
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
	WorkDir string                 `json:"workDir,omitempty"`
	Name    *string                `json:"name,omitempty"`
	Ports   map[string]string      `json:"ports,omitempty"`
}

type SCGContainerCallImage struct {
	// @TODO: remove after release 0.1.28 in favor of image.ref
	Src       *string       `json:"src,omitempty"`
	Ref       *string       `json:"ref,omitempty"`
	PullCreds *SCGPullCreds `json:"pullCreds,omitempty"`
}

type SCGLoopVars struct {
	Index *string `json:"index,omitempty"`
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type SCGOpCall struct {
	// Ref represents a references to an op; will be interpolated
	Ref string `json:"ref"`
	// PullCreds represent creds for pulling the op from a provider
	PullCreds *SCGPullCreds `json:"pullCreds,omitempty"`
	// binds scope to inputs of referenced op
	Inputs map[string]interface{} `json:"inputs,omitempty"`
	// binds scope to outputs of referenced op
	Outputs map[string]string `json:"outputs,omitempty"`
}

type SCGParallelLoopCall struct {
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

type SCGSerialLoopCall struct {
	Range interface{}     `json:"range,omitempty"`
	Run   SCG             `json:"run,omitempty"`
	Until []*SCGPredicate `json:"until,omitempty"`
	Vars  *SCGLoopVars    `json:"vars,omitempty"`
}
