package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Context("pullImage", func() {
	dbDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}

	db, err := badger.Open(
		badger.DefaultOptions(dbDir).WithLogger(nil),
	)
	if err != nil {
		panic(err)
	}

	pubSub := pubsub.New(db)

	Context("imageRef valid", func() {
		It("should call dockerClient.ImagePull w/ expected args", func() {
			/* arrange */
			providedImageRef := "imageRef"
			expectedImagePullOptions := types.ImagePullOptions{Platform: "linux"}
			providedCtx := context.Background()

			imagePullResponse := io.NopCloser(bytes.NewBufferString(""))

			_fakeDockerClient := new(FakeCommonAPIClient)
			_fakeDockerClient.ImagePullReturns(imagePullResponse, nil)

			/* act */
			err := pullImage(
				providedCtx,
				&model.ContainerCall{
					ContainerID: "",
					Image: &model.ContainerCallImage{
						PullCreds: &model.Creds{},
						Ref:       &providedImageRef,
					},
				},
				_fakeDockerClient,
				"",
				pubSub,
			)
			if err != nil {
				panic(err)
			}

			/* assert */
			actualCtx, actualImageRef, actualImagePullOptions := _fakeDockerClient.ImagePullArgsForCall(0)
			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualImageRef).To(Equal(providedImageRef))
			Expect(actualImagePullOptions).To(Equal(expectedImagePullOptions))
		})
		It("should skip pulling when image is present and ref is tagged non-latest", func() {
			/* arrange */
			providedImageRef := "imageRef:myversion"
			providedCtx := context.Background()
			_fakeDockerClient := new(FakeCommonAPIClient)

			/* act */
			err := pullImage(
				providedCtx,
				&model.ContainerCall{
					ContainerID: "",
					Image: &model.ContainerCallImage{
						PullCreds: &model.Creds{},
						Ref:       &providedImageRef,
					},
				},
				_fakeDockerClient,
				"",
				pubSub,
			)
			if err != nil {
				panic(err)
			}

			/* assert */
			// Checked if image exists
			ctx, inspectedImageRef := _fakeDockerClient.ImageInspectWithRawArgsForCall(0)
			Expect(ctx).To(Equal(providedCtx))
			Expect(inspectedImageRef).To(Equal(providedImageRef))
			// Didn't pull
			Expect(_fakeDockerClient.ImagePullCallCount()).To(Equal(0))
		})
		It("should pul when image is not present and ref is tagged non-latest", func() {
			/* arrange */
			providedImageRef := "imageRef:myversion"
			expectedImagePullOptions := types.ImagePullOptions{Platform: "linux"}
			providedCtx := context.Background()

			imagePullResponse := io.NopCloser(bytes.NewBufferString(""))

			_fakeDockerClient := new(FakeCommonAPIClient)
			_fakeDockerClient.ImagePullReturns(imagePullResponse, nil)
			_fakeDockerClient.ImageInspectWithRawReturns(types.ImageInspect{}, nil, errors.New("not found"))

			/* act */
			err := pullImage(
				providedCtx,
				&model.ContainerCall{
					ContainerID: "",
					Image: &model.ContainerCallImage{
						PullCreds: &model.Creds{},
						Ref:       &providedImageRef,
					},
				},
				_fakeDockerClient,
				"",
				pubSub,
			)
			if err != nil {
				panic(err)
			}

			/* assert */
			// Checked if image exists
			ctx, inspectedImageRef := _fakeDockerClient.ImageInspectWithRawArgsForCall(0)
			Expect(ctx).To(Equal(providedCtx))
			Expect(inspectedImageRef).To(Equal(providedImageRef))
			// Pulled
			actualCtx, actualImageRef, actualImagePullOptions := _fakeDockerClient.ImagePullArgsForCall(0)
			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualImageRef).To(Equal(providedImageRef))
			Expect(actualImagePullOptions).To(Equal(expectedImagePullOptions))
		})
		Context("dockerClient.ImagePull errors", func() {
			It("should return expected error", func() {
				/* arrange */
				imagePullError := errors.New("dummyerror")
				expectedError := imagePullError
				imagePullResponse := io.NopCloser(bytes.NewBufferString(""))

				_fakeDockerClient := new(FakeCommonAPIClient)
				_fakeDockerClient.ImagePullReturns(imagePullResponse, imagePullError)

				imageRef := "dummyImageRef"

				/* act */
				actualError := pullImage(
					context.Background(),
					&model.ContainerCall{
						ContainerID: "",
						Image: &model.ContainerCallImage{
							PullCreds: &model.Creds{},
							Ref:       &imageRef,
						},
					},
					_fakeDockerClient,
					"",
					pubSub,
				)

				/* assert */
				Expect(actualError).To(MatchError(expectedError))
			})
		})
	})
})
