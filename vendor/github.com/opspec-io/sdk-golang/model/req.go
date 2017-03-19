package model

type EventFilter struct {
	RootOpIds []string
}

type GetEventStreamReq struct {
	Filter *EventFilter
}

type KillOpReq struct {
	OpId string
}

type StartOpReq struct {
	// map of args keyed by param name
	Args   map[string]*Data
	PkgRef string
}
