package pkg

type CreateReq struct {
	Path        string
	Name        string
	Description string
}

type GetReq struct {
	Path     string
	PkgRef   string
	Username string
	Password string
}

type SetDescriptionReq struct {
	Path        string
	Description string
}
