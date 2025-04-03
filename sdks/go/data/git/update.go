package git

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/opctl/opctl/sdks/go/model"
)

func update(
	ctx context.Context,
	repoPath string,
	authOpts *model.Creds,
) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	branchName := headRef.Name()
	if !branchName.IsBranch() {
		return err
	}

	fetchOptions := &git.FetchOptions{
		Depth: 1,
		Force: true,
		Prune: true,
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
		},
	}

	if authOpts != nil {
		fetchOptions.Auth = &http.BasicAuth{
			Username: authOpts.Username,
			Password: authOpts.Password,
		}
	}

	err = repo.FetchContext(
		ctx,
		fetchOptions,
	)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	remoteRef, err := repo.Reference(
		plumbing.NewRemoteReferenceName("origin", branchName.Short()),
		true,
	)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Reset the working tree to match the fetched state
	return wt.Reset(
		&git.ResetOptions{
			Mode:   git.HardReset,
			Commit: remoteRef.Hash(),
		},
	)
}
