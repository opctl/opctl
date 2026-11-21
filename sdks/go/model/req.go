package model

import (
	"time"
)

// AddAuthReq holds data for adding source (git or OCI Distribution API) credentials
type AddAuthReq struct {
	// Resources designates which resources this auth is for in the form of a reference (or prefix of).
	Resources string
	Creds
}

type EventFilter struct {
	// filter to events from these root op id's
	Roots []string
	// filter to events occurring after & including this time
	Since *time.Time
	// Ordered delivers live events in publication order after replay.
	Ordered bool
}

type GetEventStreamReq struct {
	Filter EventFilter
}

type GetDataReq struct {
	PullCreds *Creds `json:"pullCreds,omitempty"`
	DataRef   string
}

type ListDescendantsReq struct {
	PullCreds *Creds `json:"pullCreds,omitempty"`
	DataRef   string `json:"dataRef"`
}

type KillOpReq struct {
	OpID       string `json:"opId"`
	RootCallID string `json:"rootCallId"`
}

type StartOpReq struct {
	// map of args keyed by input name
	Args map[string]*Value `json:"args,omitempty"`
	// Op details the op to start
	Op          StartOpReqOp `json:"op,omitempty"`
	TermSize    *TermSize    `json:"termSize,omitempty"`
	Interactive bool         `json:"interactive,omitempty"`
}

type TermSize struct {
	Rows uint `json:"rows"`
	Cols uint `json:"cols"`
}

type StartOpReqOp struct {
	Ref       string
	PullCreds *Creds `json:"pullCreds,omitempty"`
}
