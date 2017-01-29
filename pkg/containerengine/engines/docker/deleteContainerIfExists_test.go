package docker

import (
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Context("DeleteContainerIfExists", func() {
	It("should call dockerClient.ContainerRemove w/ expected args", func() {
		/* arrange */
		_fakeDockerClient := new(fakeDockerClient)

		providedContainerId := "dummyContainerId"
		expectedContainerId := providedContainerId
		expectedContainerRemoveOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}

		objectUnderTest := _containerEngine{
			dockerClient: _fakeDockerClient,
		}

		/* act */
		objectUnderTest.DeleteContainerIfExists(providedContainerId)

		/* assert */
		_, actualContainerId, actualContainerRemoveOptions := _fakeDockerClient.ContainerRemoveArgsForCall(0)
		Expect(actualContainerId).To(Equal(expectedContainerId))
		Expect(actualContainerRemoveOptions).Should(Equal(expectedContainerRemoveOptions))
	})
	Context("dockerClient.ContainerRemove errors", func() {
		It("should return", func() {
			/* arrange */
			_fakeDockerClient := new(fakeDockerClient)
			_fakeDockerClient.ContainerRemoveReturns(errors.New("dummyError"))

			objectUnderTest := _containerEngine{
				dockerClient: _fakeDockerClient,
			}

			/* act/assert */
			objectUnderTest.DeleteContainerIfExists("dummyContainerId")
		})
	})
	Context("dockerClient.ContainerRemove doesn't error", func() {
		It("should return", func() {
			/* arrange */
			objectUnderTest := _containerEngine{
				dockerClient: new(fakeDockerClient),
			}

			/* act/assert */
			objectUnderTest.DeleteContainerIfExists("dummyContainerId")
		})
	})
})
