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

// GetPkgContentReq deprecated
type GetPkgContentReq struct {
	ContentPath string
	PullCreds   *PullCreds `json:",omitempty"`
	PkgRef      string
}

// ListPkgContentsReq deprecated
type ListPkgContentsReq struct {
	PullCreds *PullCreds `json:",omitempty"`
	PkgRef    string
}

type KillOpReq struct {
	OpID string `json:"opId"`
}

type StartOpReq struct {
	// map of args keyed by input name
	Args map[string]*Value
	// Op details the op to start
	Op StartOpReqOp
}

type StartOpReqOp struct {
	Ref       string
	PullCreds *PullCreds `json:",omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface to handle deprecated properties gracefully in one place
func (sor *StartOpReq) UnmarshalJSON(
	b []byte,
) error {
	if err := json.Unmarshal(b, sor); nil != err {
		return err
	}

	// handle deprecated property
	deprecated := struct {
		Pkg *DCGOpCallPkg `json:"pkg"`
	}{}
	if err := json.Unmarshal(b, &deprecated); nil != err {
		return err
	}

	if nil != deprecated.Pkg {
		sor.Op.Ref = deprecated.Pkg.Ref
		sor.Op.PullCreds = deprecated.Pkg.PullCreds
	}
	return nil

}
