package docker

import (
	"fmt"
	"github.com/docker/docker/api/types/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Context("NetworkContainer", func() {
	It("should call dockerClient.NetworkConnect w/ expected args", func() {
		/* arrange */
		fakeDockerClient := new(fakeDockerClient)

		expectedNetworkId := "dummyNetworkId"
		expectedContainerId := "dummyContainerId"
		expectedContainerAlias := "dummyContainerAlias"
		expectedEndpointSettings := &network.EndpointSettings{
			Aliases: []string{expectedContainerAlias},
		}

		objectUnderTest := _containerProvider{
			dockerClient: fakeDockerClient,
		}

		/* act */
		objectUnderTest.NetworkContainer(expectedNetworkId, expectedContainerId, expectedContainerAlias)

		/* assert */
		_, actualNetworkId, actualContainerId, actualNetworkConnectOptions :=
			fakeDockerClient.NetworkConnectArgsForCall(0)

		Expect(actualNetworkId).To(Equal(expectedNetworkId))
		Expect(actualContainerId).To(Equal(expectedContainerId))
		Expect(actualNetworkConnectOptions).To(Equal(expectedEndpointSettings))
	})
	Context("dockerClient.NetworkConnect errors", func() {
		It("should return", func() {
			/* arrange */
			errorReturnedFromNetworkConnect := errors.New("dummyError")

			fakeDockerClient := new(fakeDockerClient)
			fakeDockerClient.NetworkConnectReturns(errorReturnedFromNetworkConnect)

			expectedError := fmt.Errorf(
				"Unable to network container. Response from docker was:\n %v",
				errorReturnedFromNetworkConnect.Error(),
			)

			objectUnderTest := _containerProvider{
				dockerClient: fakeDockerClient,
			}

			/* act */
			actualError := objectUnderTest.NetworkContainer("", "", "")

			/* assert */
			Expect(actualError).To(Equal(expectedError))
		})
	})
	Context("dockerClient.NetworkConnect doesn't error", func() {
		It("shouldn't error", func() {
			/* arrange */
			objectUnderTest := _containerProvider{
				dockerClient: new(fakeDockerClient),
			}

			/* act */
			actualError := objectUnderTest.NetworkContainer("", "", "")

			/* assert */
			Expect(actualError).To(BeNil())
		})
	})
})
