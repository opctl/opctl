package pkg

type CreateReq struct {
	Path        string
	Name        string
	Description string
}

type GetReq struct {
	BasePath string
	PkgRef   string
	Username string
	Password string
}

type SetDescriptionReq struct {
	Path        string
	Description string
}
