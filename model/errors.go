package model

// ErrPkgPullAuthentication conveys pkg pull failed due to authentication
type ErrPkgPullAuthentication struct{}

func (ear ErrPkgPullAuthentication) Error() string {
	return "Pkg pull failed due to invalid/lack of authentication"
}

// ErrPkgPullAuthorization conveys pkg pull failed due to authorization
type ErrPkgPullAuthorization struct{}

func (ear ErrPkgPullAuthorization) Error() string {
	return "Pkg pull failed due to insufficient/lack of authorization"
}

// ErrPkgNotFound conveys no such package could be found
type ErrPkgNotFound struct{}

func (ear ErrPkgNotFound) Error() string {
	return "Provider failed to resolve the requested pkg"
}
