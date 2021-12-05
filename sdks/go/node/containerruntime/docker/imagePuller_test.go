package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/dgraph-io/badger/v3"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

var _ = Context("imagePuller", func() {
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

			objectUnderTest := _imagePuller{
				dockerClient: _fakeDockerClient,
			}

			/* act */
			err := objectUnderTest.Pull(
				providedCtx,
				&model.ContainerCall{
					ContainerID: "",
					Image: &model.ContainerCallImage{
						PullCreds: &model.Creds{},
						Ref:       &providedImageRef,
					},
				},
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
			dockerClient := new(FakeCommonAPIClient)

			objectUnderTest := _imagePuller{
				dockerClient: dockerClient,
			}

			/* act */
			err := objectUnderTest.Pull(
				providedCtx,
				&model.ContainerCall{
					ContainerID: "",
					Image: &model.ContainerCallImage{
						PullCreds: &model.Creds{},
						Ref:       &providedImageRef,
					},
				},
				"",
				pubSub,
			)
			if err != nil {
				panic(err)
			}

			/* assert */
			// Checked if image exists
			ctx, inspectedImageRef := dockerClient.ImageInspectWithRawArgsForCall(0)
			Expect(ctx).To(Equal(providedCtx))
			Expect(inspectedImageRef).To(Equal(providedImageRef))
			// Didn't pull
			Expect(dockerClient.ImagePullCallCount()).To(Equal(0))
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

			objectUnderTest := _imagePuller{
				dockerClient: _fakeDockerClient,
			}

			/* act */
			err := objectUnderTest.Pull(
				providedCtx,
				&model.ContainerCall{
					ContainerID: "",
					Image: &model.ContainerCallImage{
						PullCreds: &model.Creds{},
						Ref:       &providedImageRef,
					},
				},
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

				objectUnderTest := _imagePuller{
					dockerClient: _fakeDockerClient,
				}

				imageRef := "dummyImageRef"

				/* act */
				actualError := objectUnderTest.Pull(
					context.Background(),
					&model.ContainerCall{
						ContainerID: "",
						Image: &model.ContainerCallImage{
							PullCreds: &model.Creds{},
							Ref:       &imageRef,
						},
					},
					"",
					pubSub,
				)

				/* assert */
				Expect(actualError).To(MatchError(expectedError))
			})
		})
	})
})
