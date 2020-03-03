package hostruntime

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker/hostruntime/internal/fakes"
	uuid "github.com/satori/go.uuid"
)

var _ = Context("RuntimeInfo", func() {
	Context("when not running in a container", func() {
		cli := FakeContainerInspector{}
		cu := containerUtils{
			inAContainer: func() bool { return false },
		}

		It("should return empty info", func() {
			/* act */
			actual, _ := newContainerRuntimeInfo(context.TODO(), &cli, cu)

			/* assert */
			Expect(actual.InAContainer).To(BeFalse())
			Expect(actual.HostPathMap).To(BeEmpty())
			Expect(actual.DockerID).To(BeEmpty())
		})
	})

	Context("when running in a container", func() {
		dockerID, _ := uuid.NewV4()
		cu := containerUtils{
			inAContainer: func() bool { return true },
			getDockerID:  func() (string, error) { return dockerID.String(), nil },
		}

		Context("when given docker engine doesn't host opctl container", func() {
			cli := FakeContainerInspector{}
			cli.ContainerInspectReturns(types.ContainerJSON{}, FakeNotFoundError{})

			It("should return empty info", func() {
				/* act */
				actual, _ := newContainerRuntimeInfo(context.TODO(), &cli, cu)

				/* assert */
				Expect(actual.InAContainer).To(BeFalse())
				Expect(actual.HostPathMap).To(BeEmpty())
			})
		})

		Context("when given docker engine hosts opctl container", func() {
			cli := FakeContainerInspector{}
			containerJSON := types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					HostConfig: &container.HostConfig{
						Binds: []string{"/host/app:/app"},
					},
				},
			}
			cli.ContainerInspectReturns(containerJSON, nil)

			It("should create HostPathMap", func() {
				/* act */
				actual, _ := newContainerRuntimeInfo(context.TODO(), &cli, cu)

				/* assert */
				Expect(actual.InAContainer).To(BeTrue())
				Expect(actual.HostPathMap).To(Not(BeEmpty()))
				Expect(actual.DockerID).To(Equal(dockerID.String()))
			})
		})
	})
})
