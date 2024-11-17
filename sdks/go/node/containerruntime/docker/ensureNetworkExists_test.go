package docker

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
)

var _ = Context("ensureNetworkExists", func() {
	Context("is NotFoundErr", func() {
		It("should call dockerClient.NetworkCreate w/ expected args", func() {
			/* arrange */
			fakeDockerClient := new(FakeCommonAPIClient)

			providedContainerID := "dummyContainerID"
			expectedContainerID := providedContainerID
			expectedNetworkCreations := network.CreateOptions{
				Attachable: true,
			}

			/* act */
			ensureNetworkExists(
				context.Background(),
				fakeDockerClient,
				&model.Creds{},
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
					network.Inspect{},
					dockerNotFoundError{
						errors.New("dummyError"),
					},
				)

				fakeDockerClient.NetworkCreateReturns(network.CreateResponse{}, errors.New("dummyError"))

				/* act */
				actualError := ensureNetworkExists(
					context.Background(),
					fakeDockerClient,
					&model.Creds{},
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

				/* act */
				actualError := ensureNetworkExists(
					context.Background(),
					fakeDockerClient,
					&model.Creds{},
					"",
				)

				/* assert */
				Expect(actualError).To(BeNil())
			})
		})
	})
})

type dockerNotFoundError struct {
	error
}

func (dockerNotFoundError) NotFound() {}
