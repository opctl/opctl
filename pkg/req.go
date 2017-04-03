package pkg

type CreateReq struct {
	Path        string
	Name        string
	Description string
}

type GetReq struct {
	path     string
	pkgRef   string
	username string
	password string
}

type SetDescriptionReq struct {
	Path        string
	Description string
}
