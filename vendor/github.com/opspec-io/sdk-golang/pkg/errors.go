package pkg

// ErrAuthenticationFailed conveys authentication failed while attempting to transport a package
type ErrAuthenticationFailed struct{}

func (ear ErrAuthenticationFailed) Error() string {
	return "Authentication failed while attempting to transport package"
}
