package oppath

import (
	"errors"
	"net/url"

	"github.com/go-git/go-git/v5"
)

func tryGetOriginURL(
	currentPath string,
) (*url.URL, error) {
	gitDirPath, err := tryFindGitDir(
		currentPath,
	)
	if err != nil {
		return nil, err
	}
	if gitDirPath == "" {
		return nil, nil
	}

	repo, err := git.PlainOpen(gitDirPath)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return nil, nil
		}

		return nil, err
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		if errors.Is(err, git.ErrRemoteNotFound) {
			return nil, nil
		}

		return nil, err
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return nil, nil
	}

	return url.Parse(urls[0])

}
