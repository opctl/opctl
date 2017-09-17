package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Context("DeleteContainerIfExists", func() {
	It("should call dockerClient.ContainerRemove w/ expected args", func() {
		/* arrange */
		fakeDockerClient := new(fakeDockerClient)

		providedContainerId := "dummyContainerId"
		expectedContainerId := providedContainerId
		expectedContainerRemoveOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}

		objectUnderTest := _containerProvider{
			dockerClient: fakeDockerClient,
		}

		/* act */
		objectUnderTest.DeleteContainerIfExists(providedContainerId)

		/* assert */
		_, actualContainerId, actualContainerRemoveOptions := fakeDockerClient.ContainerRemoveArgsForCall(0)
		Expect(actualContainerId).To(Equal(expectedContainerId))
		Expect(actualContainerRemoveOptions).To(Equal(expectedContainerRemoveOptions))
	})
	Context("dockerClient.ContainerRemove errors", func() {
		It("should return", func() {
			/* arrange */
			errorReturnedFromContainerRemove := errors.New("dummyError")

			fakeDockerClient := new(fakeDockerClient)
			fakeDockerClient.ContainerRemoveReturns(errorReturnedFromContainerRemove)

			expectedError := fmt.Errorf(
				"unable to delete container. Response from docker was:\n %v",
				errorReturnedFromContainerRemove.Error(),
			)

			objectUnderTest := _containerProvider{
				dockerClient: fakeDockerClient,
			}

			/* act */
			actualError := objectUnderTest.DeleteContainerIfExists("")

			/* assert */
			Expect(actualError).To(Equal(expectedError))
		})
	})
	Context("dockerClient.ContainerRemove doesn't error", func() {
		It("shouldn't error", func() {
			/* arrange */
			objectUnderTest := _containerProvider{
				dockerClient: new(fakeDockerClient),
			}

			/* act */
			actualError := objectUnderTest.DeleteContainerIfExists("")

			/* assert */
			Expect(actualError).To(BeNil())
		})
	})
})
