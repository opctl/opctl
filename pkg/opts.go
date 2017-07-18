package pkg

import "github.com/opspec-io/sdk-golang/model"

type ResolveOpts struct {
	PullCreds *model.PullCreds
	// BasePath provides the base path for relative reference resolution
	BasePath string
}
