package docker

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Context("DeleteNetworkIfExists", func() {
	It("should call dockerClient.NetworkRemove w/ expected args", func() {
		/* arrange */
		_fakeDockerClient := new(fakeDockerClient)

		providedContainerId := "dummyContainerId"
		expectedContainerId := providedContainerId

		objectUnderTest := _containerProvider{
			dockerClient: _fakeDockerClient,
		}

		/* act */
		objectUnderTest.DeleteNetworkIfExists(providedContainerId)

		/* assert */
		_, actualContainerId := _fakeDockerClient.NetworkRemoveArgsForCall(0)
		Expect(actualContainerId).To(Equal(expectedContainerId))
	})
	Context("dockerClient.NetworkRemove errors", func() {
		It("should return expected error", func() {
			/* arrange */
			errorReturnedFromNetworkRemove := errors.New("dummyError")

			fakeDockerClient := new(fakeDockerClient)
			fakeDockerClient.NetworkRemoveReturns(errorReturnedFromNetworkRemove)

			expectedError := fmt.Errorf(
				"Unable to delete network. Response from docker was:\n %v",
				errorReturnedFromNetworkRemove.Error(),
			)

			objectUnderTest := _containerProvider{
				dockerClient: fakeDockerClient,
			}

			/* act */
			actualError := objectUnderTest.DeleteNetworkIfExists("")

			/* assert */
			Expect(actualError).To(Equal(expectedError))
		})
	})
	Context("dockerClient.NetworkRemove doesn't error", func() {
		It("shouldn't error", func() {
			/* arrange */
			objectUnderTest := _containerProvider{
				dockerClient: new(fakeDockerClient),
			}

			/* act */
			actualError := objectUnderTest.DeleteNetworkIfExists("")

			/* assert */
			Expect(actualError).To(BeNil())
		})
	})
})
