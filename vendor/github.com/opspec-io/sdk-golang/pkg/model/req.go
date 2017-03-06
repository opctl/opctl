package model

type CreateCollectionReq struct {
	Path        string
	Name        string
	Description string
}

type CreateOpReq struct {
	Path        string
	Name        string
	Description string
}

type EventFilter struct {
	RootOpIds []string
}

type GetEventStreamReq struct {
	Filter *EventFilter
}

type KillOpReq struct {
	OpId string
}

type SetCollectionDescriptionReq struct {
	PathToCollection string
	Description      string
}

type SetOpDescriptionReq struct {
	PathToOp    string
	Description string
}

type StartOpReq struct {
	// map of args keyed by param name
	Args     map[string]*Data `json:"args"`
	OpPkgRef string           `json:"opPkgRef"`
}
