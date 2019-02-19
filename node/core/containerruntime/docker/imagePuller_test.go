package docker

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/pkg/errors"
	"io/ioutil"
)

var _ = Context("imagePuller", func() {
	Context("imageRef valid", func() {
		It("should call dockerClient.ImagePull w/ expected args", func() {
			/* arrange */
			providedImage := &model.DCGContainerCallImage{Ref: "dummy-ref"}
			expectedImagePullOptions := types.ImagePullOptions{}
			providedCtx := context.Background()

			imagePullResponse := ioutil.NopCloser(bytes.NewBufferString(""))

			_fakeDockerClient := new(fakeDockerClient)
			_fakeDockerClient.ImagePullReturns(imagePullResponse, nil)

			objectUnderTest := _imagePuller{
				dockerClient: _fakeDockerClient,
			}

			/* act */
			err := objectUnderTest.Pull(
				providedCtx,
				providedImage,
				"",
				"",
				new(pubsub.FakeEventPublisher),
			)
			if nil != err {
				panic(err)
			}

			/* assert */
			actualCtx, actualImageRef, actualImagePullOptions := _fakeDockerClient.ImagePullArgsForCall(0)
			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualImageRef).To(Equal(providedImage.Ref))
			Expect(actualImagePullOptions).To(Equal(expectedImagePullOptions))
		})
		Context("dockerClient.ImagePull errors", func() {
			It("should return expected error", func() {
				/* arrange */
				imagePullError := errors.New("dummyerror")
				expectedError := imagePullError
				imagePullResponse := ioutil.NopCloser(bytes.NewBufferString(""))

				_fakeDockerClient := new(fakeDockerClient)
				_fakeDockerClient.ImagePullReturns(imagePullResponse, imagePullError)

				objectUnderTest := _imagePuller{
					dockerClient: _fakeDockerClient,
				}

				/* act */
				actualError := objectUnderTest.Pull(
					context.Background(),
					&model.DCGContainerCallImage{Ref: "dummy-ref"},
					"",
					"",
					new(pubsub.FakeEventPublisher),
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})
