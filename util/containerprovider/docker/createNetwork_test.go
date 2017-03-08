package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Context("CreateNetwork", func() {
	It("should call dockerClient.NetworkCreate w/ expected args", func() {
		/* arrange */
		fakeDockerClient := new(fakeDockerClient)

		providedContainerId := "dummyContainerId"
		expectedContainerId := providedContainerId
		expectedNetworkCreatePackagetions := types.NetworkCreate{
			Attachable: true,
		}

		objectUnderTest := _containerProvider{
			dockerClient: fakeDockerClient,
		}

		/* act */
		objectUnderTest.CreateNetwork(providedContainerId)

		/* assert */
		_, actualContainerId, actualNetworkCreatePackagetions := fakeDockerClient.NetworkCreateArgsForCall(0)
		Expect(actualContainerId).To(Equal(expectedContainerId))
		Expect(actualNetworkCreatePackagetions).Should(Equal(expectedNetworkCreatePackagetions))
	})
	Context("dockerClient.NetworkCreate errors", func() {
		It("should return expected error", func() {
			/* arrange */
			errorReturnedFromNetworkCreate := errors.New("dummyError")

			fakeDockerClient := new(fakeDockerClient)
			fakeDockerClient.NetworkCreateReturns(types.NetworkCreateResponse{}, errorReturnedFromNetworkCreate)

			expectedError := fmt.Errorf(
				"Unable to create network. Response from docker was:\n %v",
				errorReturnedFromNetworkCreate.Error(),
			)

			objectUnderTest := _containerProvider{
				dockerClient: fakeDockerClient,
			}

			/* act */
			actualError := objectUnderTest.CreateNetwork("")

			/* assert */
			Expect(actualError).To(Equal(expectedError))
		})
	})
	Context("dockerClient.NetworkCreate doesn't error", func() {
		It("shouldn't error", func() {
			/* arrange */
			objectUnderTest := _containerProvider{
				dockerClient: new(fakeDockerClient),
			}

			/* act */
			actualError := objectUnderTest.CreateNetwork("")

			/* assert */
			Expect(actualError).To(BeNil())
		})
	})
})
