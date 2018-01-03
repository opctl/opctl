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
	Filter *EventFilter
}

type GetPkgContentReq struct {
	ContentPath string
	PullCreds   *PullCreds
	PkgRef      string
}

type ListPkgContentsReq struct {
	PullCreds *PullCreds
	PkgRef    string
}

type KillOpReq struct {
	OpId string
}

type StartOpReq struct {
	// map of args keyed by input name
	Args map[string]interface{}
	Pkg  *DCGOpCallPkg
}
