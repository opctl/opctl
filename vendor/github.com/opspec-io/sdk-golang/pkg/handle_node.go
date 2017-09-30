package pkg

import (
	"context"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api/client"
)

func newNodeHandle(
	client client.Client,
	pkgRef string,
	pullCreds *model.PullCreds,
) model.PkgHandle {
	return nodeHandle{
		client:    client,
		pkgRef:    pkgRef,
		pullCreds: pullCreds,
	}
}

func (nh nodeHandle) GetContent(
	ctx context.Context,
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {
	return nh.client.GetPkgContent(
		ctx,
		model.GetPkgContentReq{
			ContentPath: contentPath,
			PkgRef:      nh.pkgRef,
			PullCreds:   nh.pullCreds,
		},
	)
}

// nodeHandle allows interacting w/ a package sourced from a node
type nodeHandle struct {
	client    client.Client
	pkgRef    string
	pullCreds *model.PullCreds
}

func (nh nodeHandle) ListContents(
	ctx context.Context,
) (
	[]*model.PkgContent,
	error,
) {
	return nh.client.ListPkgContents(
		ctx,
		model.ListPkgContentsReq{
			PkgRef:    nh.pkgRef,
			PullCreds: nh.pullCreds,
		},
	)
}

func (hn nodeHandle) Path() *string {
	return nil
}

func (nh nodeHandle) Ref() string {
	return nh.pkgRef
}
