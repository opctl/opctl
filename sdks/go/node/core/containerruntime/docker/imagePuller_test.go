package docker

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker/internal/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("imagePuller", func() {
	Context("imageRef valid", func() {
		It("should call dockerClient.ImagePull w/ expected args", func() {
			/* arrange */
			providedImageRef := "imageRef"
			expectedImagePullOptions := types.ImagePullOptions{}
			providedCtx := context.Background()

			imagePullResponse := ioutil.NopCloser(bytes.NewBufferString(""))

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
				new(FakeEventPublisher),
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

			imagePullResponse := ioutil.NopCloser(bytes.NewBufferString(""))

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
				new(FakeEventPublisher),
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
		Context("dockerClient.ImagePull errors", func() {
			It("should return expected error", func() {
				/* arrange */
				imagePullError := errors.New("dummyerror")
				expectedError := imagePullError
				imagePullResponse := ioutil.NopCloser(bytes.NewBufferString(""))

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
					new(FakeEventPublisher),
				)

				/* assert */
				Expect(actualError).To(MatchError(expectedError))
			})
		})
	})
})
