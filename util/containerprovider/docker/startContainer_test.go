package docker

import (
	"bytes"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/pubsub"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
)

var _ = Context("StartContainer", func() {
	It("should call dockerClient.ContainerCreate w/ expected args", func() {
		/* arrange */
		providedReq := &containerprovider.StartContainerReq{
			ContainerId: "dummyContainerId",
			Dirs: map[string]string{
				"dir1ContainerPath": "dir1HostPath",
				"dir2ContainerPath": "dir2HostPath",
			},
			Env: map[string]string{
				"envVar1Name": "envVar1Value",
				"envVar2Name": "envVar2Value",
				"envVar3Name": "envVar3Value",
			},
			Files: map[string]string{
				"file1ContainerPath": "file1HostPath",
				"file2ContainerPath": "file2HostPath",
			},
			Image: "dummyImage",
			Sockets: map[string]string{
				"/unixSocket1ContainerAddress": "/unixSocket1HostAddress",
				"/unixSocket2ContainerAddress": "/unixSocket2HostAddress",
			},
			WorkDir: "dummyWorkDir",
		}

		expectedConfig := &container.Config{
			Env: []string{
				"envVar1Name=envVar1Value",
				"envVar2Name=envVar2Value",
				"envVar3Name=envVar3Value",
			},
			Image:      providedReq.Image,
			Tty:        true,
			WorkingDir: providedReq.WorkDir,
		}
		expectedHostConfig := &container.HostConfig{
			Binds: []string{
				"/unixSocket1HostAddress:/unixSocket1ContainerAddress",
				"/unixSocket2HostAddress:/unixSocket2ContainerAddress",
				"dir1HostPath:dir1ContainerPath",
				"dir2HostPath:dir2ContainerPath",
				"file1HostPath:file1ContainerPath",
				"file2HostPath:file2ContainerPath",
			},
			Privileged: true,
		}
		expectedNetworkingConfig := &network.NetworkingConfig{}

		_fakeDockerClient := new(fakeDockerClient)

		_fakeDockerClient.ContainerLogsStub = func(
			ctx context.Context,
			container string,
			options types.ContainerLogsOptions,
		) (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewBufferString("")), nil
		}

		objectUnderTest := _containerProvider{
			dockerClient: _fakeDockerClient,
		}

		/* act */
		objectUnderTest.StartContainer(providedReq, new(pubsub.FakeEventPublisher))

		/* assert */
		_,
			actualConfig,
			actualHostConfig,
			actualNetworkingConfig,
			actualContainerName := _fakeDockerClient.ContainerCreateArgsForCall(0)
		Expect(actualConfig).To(Equal(expectedConfig))
		Expect(actualHostConfig).To(Equal(expectedHostConfig))
		Expect(actualNetworkingConfig).To(Equal(expectedNetworkingConfig))
		Expect(actualContainerName).To(Equal(providedReq.ContainerId))
	})
})
