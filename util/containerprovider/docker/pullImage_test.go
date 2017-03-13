package docker

import (
	"bytes"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/reference"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/pkg/errors"
	"io/ioutil"
)

var _ = Context("pullImage", func() {
	Context("imageRef invalid", func() {
		It("should return expected error", func() {
			/* arrange */
			providedImage := &model.DCGContainerCallImage{Ref: "%$^"}
			_, _, expectedError := reference.Parse(providedImage.Ref)

			objectUnderTest := _containerProvider{
				dockerClient: new(fakeDockerClient),
			}

			/* act */
			actualError := objectUnderTest.pullImage(
				providedImage,
				"",
				"",
				new(pubsub.FakeEventPublisher),
			)

			/* assert */
			Expect(actualError).To(Equal(expectedError))
		})
	})
	Context("imageRef valid", func() {
		It("should call dockerClient.ImagePull w/ expected args", func() {
			/* arrange */
			providedImage := &model.DCGContainerCallImage{Ref: "dummy-ref"}
			expectedImageRef := fmt.Sprintf("%v:latest", providedImage.Ref)
			expectedImagePullOptions := types.ImagePullOptions{}

			imagePullResponse := ioutil.NopCloser(bytes.NewBufferString(""))

			_fakeDockerClient := new(fakeDockerClient)
			_fakeDockerClient.ImagePullReturns(imagePullResponse, nil)

			objectUnderTest := _containerProvider{
				dockerClient: _fakeDockerClient,
			}

			/* act */
			err := objectUnderTest.pullImage(
				providedImage,
				"",
				"",
				new(pubsub.FakeEventPublisher),
			)
			if nil != err {
				panic(err)
			}

			/* assert */
			_, actualImageRef, actualImagePullOptions := _fakeDockerClient.ImagePullArgsForCall(0)
			Expect(actualImageRef).To(Equal(expectedImageRef))
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

				objectUnderTest := _containerProvider{
					dockerClient: _fakeDockerClient,
				}

				/* act */
				actualError := objectUnderTest.pullImage(
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
