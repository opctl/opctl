package pkg

type ErrPkgNotFound struct{}

func (this ErrPkgNotFound) Error() string {
	return "Pkg not found"
}
