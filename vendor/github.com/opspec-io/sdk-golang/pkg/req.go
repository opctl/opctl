package pkg

type CreateReq struct {
	Path        string
	Name        string
	Description string
}

type SetDescriptionReq struct {
	Path        string
	Description string
}
