package docker

import (
	"bytes"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/iruntime"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
)

var _ = Context("RunContainer", func() {
	closedContainerWaitOkBodyChan := make(chan container.ContainerWaitOKBody)
	close(closedContainerWaitOkBodyChan)

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
			Name:    "dummyName",
			Ports: map[string]string{
				"6060/udp":  "6060",
				"8080-8081": "9090-9091",
			},
		}

		expectedConfig := &container.Config{
			Env: []string{
				"envVar1Name=envVar1Value",
				"envVar2Name=envVar2Value",
				"envVar3Name=envVar3Value",
			},
			ExposedPorts: nat.PortSet{
				"6060/udp": struct{}{},
				"8080/tcp": struct{}{},
				"8081/tcp": struct{}{},
			},
			Image:      providedReq.Image.Ref,
			WorkingDir: providedReq.WorkDir,
			Tty:        true,
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
			PortBindings: nat.PortMap{
				"6060/udp": []nat.PortBinding{{HostPort: "6060"}},
				"8080/tcp": []nat.PortBinding{{HostPort: "9090"}},
				"8081/tcp": []nat.PortBinding{{HostPort: "9091"}},
			},

			Privileged: true,
		}
		expectedNetworkingConfig := &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"opctl": {
					Aliases: []string{
						providedReq.Name,
					},
				},
			},
		}

		_fakeDockerClient := new(fakeDockerClient)

		_fakeDockerClient.ContainerLogsStub = func(
			ctx context.Context,
			container string,
			options types.ContainerLogsOptions,
		) (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewBufferString("")), nil
		}

		_fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

		objectUnderTest := _containerProvider{
			dockerClient: _fakeDockerClient,
			runtime:      new(iruntime.Fake),
		}

		/* act */
		objectUnderTest.RunContainer(
			providedReq,
			new(pubsub.FakeEventPublisher),
			nopWriteCloser{ioutil.Discard},
			nopWriteCloser{ioutil.Discard},
		)

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

type nopWriteCloser struct {
	io.Writer
}

func (w nopWriteCloser) Close() error { return nil }
