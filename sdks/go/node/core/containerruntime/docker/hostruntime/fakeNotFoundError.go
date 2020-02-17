package hostruntime

type FakeNotFoundError struct{}

func (err FakeNotFoundError) NotFound() bool {
	return true
}

func (err FakeNotFoundError) Error() string {
	return "not found"
}
