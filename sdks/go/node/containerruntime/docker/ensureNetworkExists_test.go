package docker

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
)

var _ = Context("ensureNetworkExists", func() {
	It("should call dockerClient.NetworkInspect w/ expected args", func() {
		/* arrange */
		fakeDockerClient := new(FakeCommonAPIClient)

		providedContainerID := "dummyContainerID"

		/* act */
		ensureNetworkExists(
			context.Background(),
			fakeDockerClient,
			providedContainerID,
		)

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

			/* act */
			actualErr := ensureNetworkExists(
				context.Background(),
				fakeDockerClient,
				providedContainerID,
			)

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

				/* act */
				ensureNetworkExists(
					context.Background(),
					fakeDockerClient,
					providedContainerID,
				)

				/* assert */
				_, actualContainerID, actualNetworkCreations := fakeDockerClient.NetworkCreateArgsForCall(0)
				Expect(actualContainerID).To(Equal(expectedContainerID))
				Expect(actualNetworkCreations).To(Equal(expectedNetworkCreations))
			})
			Context("dockerClient.NetworkCreate errors", func() {
				It("should return expected error", func() {
					/* arrange */
					fakeDockerClient := new(FakeCommonAPIClient)

					fakeDockerClient.NetworkInspectReturns(
						types.NetworkResource{},
						dockerNotFoundError{
							errors.New("dummyError"),
						},
					)

					fakeDockerClient.NetworkCreateReturns(types.NetworkCreateResponse{}, errors.New("dummyError"))

					/* act */
					actualError := ensureNetworkExists(
						context.Background(),
						fakeDockerClient,
						"",
					)

					/* assert */
					Expect(actualError).To(MatchError("unable to create network: dummyError"))
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

					/* act */
					actualError := ensureNetworkExists(
						context.Background(),
						fakeDockerClient,
						"",
					)

					/* assert */
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("isn't NotFoundErr", func() {
			It("should return expected error", func() {
				/* arrange */
				fakeDockerClient := new(FakeCommonAPIClient)
				fakeDockerClient.NetworkInspectReturns(types.NetworkResource{}, errors.New("dummyError"))

				/* act */
				actualError := ensureNetworkExists(
					context.Background(),
					fakeDockerClient,
					"",
				)

				/* assert */
				Expect(actualError).To(MatchError("unable to inspect network: dummyError"))
			})
		})
	})
})

type dockerNotFoundError struct {
	error
}

func (dockerNotFoundError) NotFound() bool {
	return true
}
