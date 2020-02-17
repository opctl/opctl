package docker

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/internal/iruntime"
)

var _ = Context("hostConfigFactory", func() {
	Context("Construct", func() {
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
				Binds: []string{
					"/unixSocket1HostAddress:/unixSocket1ContainerAddress",
					"/unixSocket2HostAddress:/unixSocket2ContainerAddress",
					"dir1HostPath:dir1ContainerPath:cached",
					"dir2HostPath:dir2ContainerPath:cached",
					"file1HostPath:file1ContainerPath:cached",
					"file2HostPath:file2ContainerPath:cached",
				},
				PortBindings: providedPortBindings,
				Privileged:   true,
			}

			objectUnderTest := _hostConfigFactory{
				fsPathConverter: _fsPathConverter{runtime: iruntime.New()},
			}

			/* act */
			actualHostConfig := objectUnderTest.Construct(
				providedContainerDirs,
				providedContainerFiles,
				providedContainerSockets,
				providedPortBindings,
			)

			/* assert */
			Expect(actualHostConfig).To(Equal(expectedHostConfig))
		})
	})
})
