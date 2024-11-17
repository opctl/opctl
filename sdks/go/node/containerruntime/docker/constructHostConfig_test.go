package docker

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("constructHostConfig", func() {
	It("should return expected result", func() {
		/* arrange */
		providedContainerDirs := map[string]string{
			"dir1ContainerPath": "dir1HostPath",
			"dir2ContainerPath": "dir2HostPath",
		}

		providedContainerFiles := map[string]string{
			"file1ContainerPath": "file1HostPath",
			"file2ContainerPath": "file2HostPath",
		}

		providedContainerSockets := map[string]string{
			"/unixSocket1ContainerAddress": "/unixSocket1HostAddress",
			"/unixSocket2ContainerAddress": "/unixSocket2HostAddress",
		}

		providedPortBindings := nat.PortMap{
			"6060/udp": []nat.PortBinding{{HostPort: "6060"}},
			"8080/tcp": []nat.PortBinding{{HostPort: "9090"}},
			"8081/tcp": []nat.PortBinding{{HostPort: "9091"}},
		}

		expectedHostConfig := &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:        mount.TypeBind,
					Source:      "file1HostPath",
					Target:      "file1ContainerPath",
					Consistency: mount.ConsistencyCached,
				},
				{
					Type:        mount.TypeBind,
					Source:      "file2HostPath",
					Target:      "file2ContainerPath",
					Consistency: mount.ConsistencyCached,
				},
				{
					Type:        mount.TypeBind,
					Source:      "dir1HostPath",
					Target:      "dir1ContainerPath",
					Consistency: mount.ConsistencyCached,
				},
				{
					Type:        mount.TypeBind,
					Source:      "dir2HostPath",
					Target:      "dir2ContainerPath",
					Consistency: mount.ConsistencyCached,
				},
				{
					Type:   mount.TypeBind,
					Source: "/unixSocket1HostAddress",
					Target: "/unixSocket1ContainerAddress",
				},
				{
					Type:   mount.TypeBind,
					Source: "/unixSocket2HostAddress",
					Target: "/unixSocket2ContainerAddress",
				},
			},
			PortBindings: providedPortBindings,
			Privileged:   true,
			Resources: container.Resources{
				DeviceRequests: []container.DeviceRequest{
					{
						Capabilities: [][]string{{"gpu"}},
						Count:        -1,
					},
				},
			},
		}

		/* act */
		actualHostConfig := constructHostConfig(
			providedContainerDirs,
			providedContainerFiles,
			providedContainerSockets,
			providedPortBindings,
			true,
		)

		/* assert */
		// assert mounts separately since ranged inputs won't produce deterministic output order
		Expect(actualHostConfig.Mounts).To(ConsistOf(expectedHostConfig.Mounts))
		actualHostConfig.Mounts = nil
		expectedHostConfig.Mounts = nil

		Expect(actualHostConfig).To(Equal(expectedHostConfig))
	})
})
