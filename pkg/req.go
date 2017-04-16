package pkg

type CreateReq struct {
	Path        string
	Name        string
	Description string
}

type GetReq struct {
	PkgRef   string
	Username string
	Password string
}

type SetDescriptionReq struct {
	Path        string
	Description string
}
