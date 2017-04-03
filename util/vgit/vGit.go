package vgit

//go:generate counterfeiter -o fake.go --fake-name Fake ./ VGit

import "gopkg.in/src-d/go-git.v4"

type VGit interface {
	// PlainOpen opens a git repository from the given path. It detects if the
	// repository is bare or a normal one. If the path doesn't contain a valid
	// repository ErrRepositoryNotExists is returned
	PlainOpen(path string) (*git.Repository, error)

	// PlainClone a repository into the path with the given options, isBare defines
	// if the new repository will be bare or normal. If the path is not empty
	// ErrRepositoryAlreadyExists is returned
	PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
}

func New() VGit {
	return _VGit{}
}

type _VGit struct{}

func (this _VGit) PlainOpen(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func (this _VGit) PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	return git.PlainClone(path, isBare, o)
}
