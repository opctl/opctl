package model

import (
	"encoding/json"
	"time"
)

type EventFilter struct {
	// filter to events from these root op id's
	Roots []string
	// filter to events occurring after & including this time
	Since *time.Time
}

type GetEventStreamReq struct {
	Filter EventFilter
}

// GetDataReq deprecated
type GetDataReq struct {
	ContentPath string
	PullCreds   *PullCreds `json:"pullCreds,omitempty"`
	PkgRef      string
}

// ListDescendantsReq deprecated
type ListDescendantsReq struct {
	PullCreds *PullCreds `json:"pullCreds,omitempty"`
	PkgRef    string     `json:pkgRef`
}

type KillOpReq struct {
	OpID string `json:"opId"`
}

type StartOpReq struct {
	// map of args keyed by input name
	Args map[string]*Value `json:"args,omitempty"`
	// Op details the op to start
	Op StartOpReqOp `json:"op,omitempty"`
}

type StartOpReqOp struct {
	Ref       string
	PullCreds *PullCreds `json:"pullCreds,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface to handle deprecated properties gracefully in one place
func (sor *StartOpReq) UnmarshalJSON(
	b []byte,
) error {
	// handle deprecated property
	deprecated := struct {
		Args map[string]*Value `json:"args,omitempty"`
		Op   *StartOpReqOp     `json:"op,omitempty"`
		Pkg  *StartOpReqOp     `json:"pkg,omitempty"`
	}{}
	if err := json.Unmarshal(b, &deprecated); nil != err {
		return err
	}

	sor.Args = deprecated.Args

	if nil != deprecated.Op {
		sor.Op = *deprecated.Op
	}
	if nil != deprecated.Pkg {
		sor.Op = *deprecated.Pkg
	}
	return nil

}
