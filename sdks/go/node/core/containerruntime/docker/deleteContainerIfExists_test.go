package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker/internal/fakes"
	"github.com/pkg/errors"
)

var _ = Context("DeleteContainerIfExists", func() {
	It("should call dockerClient.ContainerRemove w/ expected args", func() {
		/* arrange */
		fakeDockerClient := new(FakeCommonAPIClient)

		providedCtx := context.Background()
		providedContainerName := "dummyContainerName"
		expectedContainerName := "opctl_" + providedContainerName
		expectedContainerRemoveOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}

		objectUnderTest := _containerRuntime{
			dockerClient: fakeDockerClient,
		}

		/* act */
		objectUnderTest.DeleteContainerIfExists(
			providedCtx,
			providedContainerName,
		)

		/* assert */
		actualCtx,
			actualContainerName,
			actualContainerRemoveOptions := fakeDockerClient.ContainerRemoveArgsForCall(0)

		Expect(actualCtx).To(Equal(providedCtx))
		Expect(actualContainerName).To(Equal(expectedContainerName))
		Expect(actualContainerRemoveOptions).To(Equal(expectedContainerRemoveOptions))
	})
	Context("dockerClient.ContainerRemove errors", func() {
		It("should return", func() {
			/* arrange */
			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerRemoveReturns(errors.New("dummyError"))

			objectUnderTest := _containerRuntime{
				dockerClient: fakeDockerClient,
			}

			/* act */
			actualError := objectUnderTest.DeleteContainerIfExists(
				context.Background(),
				"containerID",
			)

			/* assert */
			Expect(actualError).To(MatchError("unable to delete container: dummyError"))
		})
	})
	Context("dockerClient.ContainerRemove doesn't error", func() {
		It("shouldn't error", func() {
			/* arrange */
			objectUnderTest := _containerRuntime{
				dockerClient: new(FakeCommonAPIClient),
			}

			/* act */
			actualError := objectUnderTest.DeleteContainerIfExists(
				context.Background(),
				"containerID",
			)

			/* assert */
			Expect(actualError).To(BeNil())
		})
	})
})
