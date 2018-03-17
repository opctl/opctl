package model

// ErrDataProviderAuthentication conveys pkg pull failed due to authentication
type ErrDataProviderAuthentication struct{}

func (ear ErrDataProviderAuthentication) Error() string {
	return "Pkg pull failed due to invalid/lack of authentication"
}

// ErrDataProviderAuthorization conveys pkg pull failed due to authorization
type ErrDataProviderAuthorization struct{}

func (ear ErrDataProviderAuthorization) Error() string {
	return "Pkg pull failed due to insufficient/lack of authorization"
}

// ErrDataRefResolution conveys no such package could be found
type ErrDataRefResolution struct{}

func (ear ErrDataRefResolution) Error() string {
	return "Provider failed to resolve the requested pkg"
}
