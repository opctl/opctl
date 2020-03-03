package docker

import (
	"bytes"
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker/internal/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
	"io/ioutil"
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
				"",
				&model.PullCreds{},
				providedImageRef,
				"",
				new(FakeEventPublisher),
			)
			if nil != err {
				panic(err)
			}

			/* assert */
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
				imagePullResponse := ioutil.NopCloser(bytes.NewBufferString(""))

				_fakeDockerClient := new(FakeCommonAPIClient)
				_fakeDockerClient.ImagePullReturns(imagePullResponse, imagePullError)

				objectUnderTest := _imagePuller{
					dockerClient: _fakeDockerClient,
				}

				/* act */
				actualError := objectUnderTest.Pull(
					context.Background(),
					"",
					nil,
					"dummyImageRef",
					"",
					new(FakeEventPublisher),
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})
