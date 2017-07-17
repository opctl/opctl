package pkg

// PullCreds contains optional authentication attributes
type PullCreds struct {
	Username,
	Password string
}

type ResolveOpts struct {
	PullCreds *PullCreds
	// BasePath provides the base path for relative reference resolution
	BasePath string
}
