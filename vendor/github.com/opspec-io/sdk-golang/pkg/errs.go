package pkg

type PkgNotFoundError string

func (this PkgNotFoundError) Error() string {
	return string(this)
}
