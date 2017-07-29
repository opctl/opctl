package pkg

// ErrAuthenticationFailed conveys authentication failed while attempting to transport a package
type ErrAuthenticationFailed struct{}

func (ear ErrAuthenticationFailed) Error() string {
	return "Authentication failed while attempting to transport package"
}

// ErrPkgNotResolved conveys the provider failed to resolve the requested package
type ErrPkgNotFound struct{}

func (ear ErrPkgNotFound) Error() string {
	return "Provider failed to resolve the requested package"
}
