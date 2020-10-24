package model

import (
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
	PkgRef    string     `json:"pkgRef"`
}

type KillOpReq struct {
	OpID     string `json:"opId"`
	RootOpID string `json:"rootOpId"`
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
