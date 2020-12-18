package core

// RunError is an error type that can be returned to allow specifying a specific
// exit code
type RunError struct {
	ExitCode int
	message  string
}

func (e *RunError) Error() string {
	return e.message
}
