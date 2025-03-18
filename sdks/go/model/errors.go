package model

import (
	"errors"
)

// ErrDataProviderAuthentication conveys data pull failed due to authentication
type ErrDataProviderAuthentication struct{}

func (ErrDataProviderAuthentication) Error() string {
	return "unauthenticated"
}

// ErrDataProviderAuthorization conveys data pull failed due to authorization
type ErrDataProviderAuthorization struct{}

func (ErrDataProviderAuthorization) Error() string {
	return "unauthorized"
}

// ErrDataNotFoundResolution conveys no such data could be found
type ErrDataNotFoundResolution struct{}

func (e ErrDataNotFoundResolution) Error() string {
	return "not found"
}

// ErrDataUnableToResolve conveys unable to resolve data
type ErrDataUnableToResolve struct{}

func (e ErrDataUnableToResolve) Error() string {
	return "unable to resolve"
}

// ErrDataGitInvalidRef conveys invalid git ref specified
type ErrDataGitInvalidRef struct{}

func (e ErrDataGitInvalidRef) Error() string {
	return "invalid git ref"
}

// ErrDataMissingVersion conveys missing version
type ErrDataMissingVersion struct{}

func (e ErrDataMissingVersion) Error() string {
	return "missing version"
}

// ErrDataSkipped conveys data was skipped
type ErrDataSkipped struct{}

func (e ErrDataSkipped) Error() string {
	return "skipped"
}

// IsAuthError returns true if this is an authorization or authentication error
func IsAuthError(err error) bool {
	return errors.Is(err, ErrDataProviderAuthorization{}) ||
		errors.Is(err, ErrDataProviderAuthentication{})
}
