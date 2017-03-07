package model

type CreateCollectionReq struct {
	Path        string
	Name        string
	Description string
}

type CreatePackageReq struct {
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

type SetPackageDescriptionReq struct {
	PathToOp    string
	Description string
}

type StartOpReq struct {
	// map of args keyed by param name
	Args     map[string]*Data `json:"args"`
	PkgRef string           `json:"pkgRef"`
}
