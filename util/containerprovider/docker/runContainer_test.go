package docker

import (
	"bytes"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/vruntime"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
)

var _ = Context("RunContainer", func() {
	It("should call dockerClient.ContainerCreate w/ expected args", func() {
		/* arrange */
		providedReq := &model.DCGContainerCall{
			DCGBaseCall: &model.DCGBaseCall{
				RootOpId: "dummyRootOpId",
				PkgRef:   "dummyPkgRef",
			},
			ContainerId: "dummyContainerId",
			Dirs: map[string]string{
				"dir1ContainerPath": "dir1HostPath",
				"dir2ContainerPath": "dir2HostPath",
			},
			EnvVars: map[string]string{
				"envVar1Name": "envVar1Value",
				"envVar2Name": "envVar2Value",
				"envVar3Name": "envVar3Value",
			},
			Files: map[string]string{
				"file1ContainerPath": "file1HostPath",
				"file2ContainerPath": "file2HostPath",
			},
			Image: &model.DCGContainerCallImage{Ref: "dummyImage"},
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
			Image:      providedReq.Image.Ref,
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
			runtime:      new(vruntime.Fake),
		}

		/* act */
		objectUnderTest.RunContainer(providedReq, new(pubsub.FakeEventPublisher))

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
