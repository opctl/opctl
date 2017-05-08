package igit

//go:generate counterfeiter -o fake.go --fake-name Fake ./ IGit

import "gopkg.in/src-d/go-git.v4"

type IGit interface {
	// PlainClone a repository into the path with the given options, isBare defines
	// if the new repository will be bare or normal. If the path is not empty
	// ErrRepositoryAlreadyExists is returned
	PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
}

func New() IGit {
	return _IGit{}
}

type _IGit struct{}

func (this _IGit) PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	return git.PlainClone(path, isBare, o)
}
