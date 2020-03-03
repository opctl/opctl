package docker

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker/internal/fakes"
)

var _ = Context("EnsureNetworkExists", func() {
	It("should call dockerClient.NetworkInspect w/ expected args", func() {
		/* arrange */
		fakeDockerClient := new(FakeCommonAPIClient)

		providedContainerID := "dummyContainerID"

		objectUnderTest := _containerRuntime{
			dockerClient: fakeDockerClient,
		}

		/* act */
		objectUnderTest.EnsureNetworkExists(providedContainerID)

		/* assert */
		_, actualContainerID, actualVerbose := fakeDockerClient.NetworkInspectArgsForCall(0)
		Expect(actualContainerID).To(Equal(providedContainerID))
		Expect(actualVerbose).To(Equal(types.NetworkInspectOptions{}))
	})
	Context("dockerClient.NetworkInspect doesn't err", func() {
		It("should return nil and not call dockerClient.NetworkCreate", func() {
			/* arrange */
			fakeDockerClient := new(FakeCommonAPIClient)

			providedContainerID := "dummyContainerID"

			objectUnderTest := _containerRuntime{
				dockerClient: fakeDockerClient,
			}

			/* act */
			actualErr := objectUnderTest.EnsureNetworkExists(providedContainerID)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(fakeDockerClient.NetworkCreateCallCount()).To(BeZero())
		})
	})
	Context("dockerClient.NetworkInspect errs", func() {
		Context("is NotFoundErr", func() {
			It("should call dockerClient.NetworkCreate w/ expected args", func() {
				/* arrange */
				fakeDockerClient := new(FakeCommonAPIClient)

				fakeDockerClient.NetworkInspectReturns(
					types.NetworkResource{},
					dockerNotFoundError{
						errors.New("dummyError"),
					},
				)

				providedContainerID := "dummyContainerID"
				expectedContainerID := providedContainerID
				expectedNetworkCreations := types.NetworkCreate{
					Attachable:     true,
					CheckDuplicate: true,
				}

				objectUnderTest := _containerRuntime{
					dockerClient: fakeDockerClient,
				}

				/* act */
				objectUnderTest.EnsureNetworkExists(providedContainerID)

				/* assert */
				_, actualContainerID, actualNetworkCreations := fakeDockerClient.NetworkCreateArgsForCall(0)
				Expect(actualContainerID).To(Equal(expectedContainerID))
				Expect(actualNetworkCreations).To(Equal(expectedNetworkCreations))
			})
			Context("dockerClient.NetworkCreate errors", func() {
				It("should return expected error", func() {
					/* arrange */
					errorReturnedFromNetworkCreate := errors.New("dummyError")

					fakeDockerClient := new(FakeCommonAPIClient)

					fakeDockerClient.NetworkInspectReturns(
						types.NetworkResource{},
						dockerNotFoundError{
							errors.New("dummyError"),
						},
					)

					fakeDockerClient.NetworkCreateReturns(types.NetworkCreateResponse{}, errorReturnedFromNetworkCreate)

					expectedError := fmt.Errorf(
						"unable to create network. Response from docker was: %v",
						errorReturnedFromNetworkCreate.Error(),
					)

					objectUnderTest := _containerRuntime{
						dockerClient: fakeDockerClient,
					}

					/* act */
					actualError := objectUnderTest.EnsureNetworkExists("")

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("dockerClient.NetworkCreate doesn't error", func() {
				It("shouldn't error", func() {
					/* arrange */
					fakeDockerClient := new(FakeCommonAPIClient)

					fakeDockerClient.NetworkInspectReturns(
						types.NetworkResource{},
						dockerNotFoundError{
							errors.New("dummyError"),
						},
					)

					objectUnderTest := _containerRuntime{
						dockerClient: fakeDockerClient,
					}

					/* act */
					actualError := objectUnderTest.EnsureNetworkExists("")

					/* assert */
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("isn't NotFoundErr", func() {
			It("should return expected error", func() {
				/* arrange */
				errorReturnedFromNetworkInspect := errors.New("dummyError")

				fakeDockerClient := new(FakeCommonAPIClient)
				fakeDockerClient.NetworkInspectReturns(types.NetworkResource{}, errorReturnedFromNetworkInspect)

				expectedError := fmt.Errorf(
					"unable to inspect network. Response from docker was: %v",
					errorReturnedFromNetworkInspect.Error(),
				)

				objectUnderTest := _containerRuntime{
					dockerClient: fakeDockerClient,
				}

				/* act */
				actualError := objectUnderTest.EnsureNetworkExists("")

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})

type dockerNotFoundError struct {
	error
}

func (this dockerNotFoundError) NotFound() bool {
	return true
}
