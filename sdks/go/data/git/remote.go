package git

import (
	"context"
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/opctl/opctl/sdks/go/model"
)

// resolveRemoteHash returns the commit SHA that the given tag points to on the
// remote.  For annotated tags it returns the peeled (dereferenced) commit SHA;
// for lightweight tags it returns the tag SHA directly.  The operation is
// equivalent to `git ls-remote` — no clone is performed.
func resolveRemoteHash(ctx context.Context, repoRef *ref, creds *model.Creds) (string, error) {
	remote := gogit.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{fmt.Sprintf("https://%v", repoRef.Name)},
	})

	listOpts := &gogit.ListOptions{
		PeelingOption: gogit.AppendPeeled,
	}
	if creds != nil {
		listOpts.Auth = &http.BasicAuth{
			Username: creds.Username,
			Password: creds.Password,
		}
	}

	refs, err := remote.ListContext(ctx, listOpts)
	if err != nil {
		return "", err
	}

	tagRef := plumbing.NewTagReferenceName(repoRef.Version)
	peeledRef := plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s^{}", repoRef.Version))

	var tagHash string
	for _, r := range refs {
		switch r.Name() {
		case peeledRef:
			// annotated tag dereferenced to commit — most precise, prefer this
			return r.Hash().String(), nil
		case tagRef:
			tagHash = r.Hash().String()
		}
	}

	if tagHash != "" {
		return tagHash, nil
	}

	return "", fmt.Errorf("tag %q not found on remote %s", repoRef.Version, repoRef.Name)
}
