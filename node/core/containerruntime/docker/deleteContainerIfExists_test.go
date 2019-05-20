package docker

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("DeleteContainerIfExists", func() {
	It("should call dockerClient.ContainerRemove w/ expected args", func() {
		/* arrange */
		fakeDockerClient := new(fakeDockerClient)

		providedContainerID := "dummyContainerID"
		expectedContainerID := providedContainerID
		expectedContainerRemoveOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}

		objectUnderTest := _containerRuntime{
			dockerClient: fakeDockerClient,
		}

		/* act */
		objectUnderTest.DeleteContainerIfExists(providedContainerID)

		/* assert */
		_, actualContainerID, actualContainerRemoveOptions := fakeDockerClient.ContainerRemoveArgsForCall(0)
		Expect(actualContainerID).To(Equal(expectedContainerID))
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

			objectUnderTest := _containerRuntime{
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
			objectUnderTest := _containerRuntime{
				dockerClient: new(fakeDockerClient),
			}

			/* act */
			actualError := objectUnderTest.DeleteContainerIfExists("")

			/* assert */
			Expect(actualError).To(BeNil())
		})
	})
})
