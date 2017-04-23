package pkg

type CreateReq struct {
	Path        string
	Name        string
	Description string
}

type PullOpts struct {
	Username string
	Password string
}

type SetDescriptionReq struct {
	Path        string
	Description string
}
