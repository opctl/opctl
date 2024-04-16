package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/docker/internal/fakes"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

var _ = Context("RunContainer", func() {
	closedContainerWaitOkBodyChan := make(chan container.WaitResponse)
	close(closedContainerWaitOkBodyChan)

	dbDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}

	db, err := badger.Open(
		badger.DefaultOptions(dbDir).WithLogger(nil),
	)
	if err != nil {
		panic(err)
	}

	pubSub := pubsub.New(db)

	It("should call dockerClient.ContainerRemove w/ expected args", func() {
		/* arrange */
		providedReq := &model.ContainerCall{
			BaseCall:    model.BaseCall{},
			ContainerID: "containerID",
			Image:       &model.ContainerCallImage{Ref: new(string)},
			// invalid to trigger early return
			Ports: map[string]string{"*": "&"},
		}

		expectedContainerRemoveOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}

		fakeDockerClient := new(FakeCommonAPIClient)
		fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

		objectUnderTest := _runContainer{
			containerStdErrStreamer: new(FakeContainerLogStreamer),
			containerStdOutStreamer: new(FakeContainerLogStreamer),
			dockerClient:            fakeDockerClient,
		}

		/* act */
		objectUnderTest.RunContainer(
			context.Background(),
			providedReq,
			"rootCallID",
			pubSub,
			nopWriteCloser{io.Discard},
			nopWriteCloser{io.Discard},
		)

		/* assert */
		_, actualContainerName, _ := fakeDockerClient.ContainerStopArgsForCall(0)
		Expect(actualContainerName).To(Equal(fmt.Sprintf("opctl_%s", providedReq.ContainerID)))

		_, actualContainerName, actualContainerRemoveOptions := fakeDockerClient.ContainerRemoveArgsForCall(0)
		Expect(actualContainerName).To(Equal(fmt.Sprintf("opctl_%s", providedReq.ContainerID)))
		Expect(actualContainerRemoveOptions).To(Equal(expectedContainerRemoveOptions))

	})
	Context("portBindingsFactory.Construct errs", func() {
		It("should return expected result", func() {
			/* arrange */

			objectUnderTest := _runContainer{
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            new(FakeCommonAPIClient),
			}

			/* act */
			_, actualErr := objectUnderTest.RunContainer(
				context.Background(),
				&model.ContainerCall{
					Image: &model.ContainerCallImage{Ref: new(string)},
					Ports: map[string]string{
						"*": "&",
					},
				},
				"rootCallID",
				pubSub,
				nopWriteCloser{io.Discard},
				nopWriteCloser{io.Discard},
			)

			/* assert */
			Expect(actualErr).To(MatchError("Invalid containerPort: *"))
		})
	})
	Context("constructPortBindings doesn't err", func() {

		It("should call dockerClient.ImagePull w/ expected args", func() {

			/* arrange */
			providedImageRef := "test"
			providedCtx := context.Background()
			providedReq := &model.ContainerCall{
				BaseCall:    model.BaseCall{},
				ContainerID: "dummyContainerID",
				Image:       &model.ContainerCallImage{Ref: &providedImageRef},
			}
			providedRootCallID := "providedRootCallID"
			expectedImagePullOptions := types.ImagePullOptions{Platform: "linux"}

			_fakeDockerClient := new(FakeCommonAPIClient)
			_fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)
			_fakeDockerClient.ImagePullReturns(io.NopCloser(bytes.NewBufferString("")), nil)

			objectUnderTest := _runContainer{
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            _fakeDockerClient,
			}

			/* act */
			objectUnderTest.RunContainer(
				providedCtx,
				providedReq,
				providedRootCallID,
				pubSub,
				nopWriteCloser{io.Discard},
				nopWriteCloser{io.Discard},
			)

			/* assert */
			actualCtx, actualImageRef, actualImagePullOptions := _fakeDockerClient.ImagePullArgsForCall(0)
			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualImageRef).To(Equal(providedImageRef))
			Expect(actualImagePullOptions).To(Equal(expectedImagePullOptions))
		})

		It("should call dockerClient.ContainerCreate w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedReq := &model.ContainerCall{
				BaseCall:    model.BaseCall{},
				ContainerID: "dummyContainerID",
				Dirs: map[string]string{
					"dir1ContainerPath": "dir1HostPath",
				},
				Files: map[string]string{
					"file1ContainerPath": "file1HostPath",
				},
				Image: &model.ContainerCallImage{Ref: new(string)},
				Name:  new(string),
				Sockets: map[string]string{
					"/unixSocket1ContainerAddress": "/unixSocket1HostAddress",
				},
				Ports: map[string]string{
					"80": "80",
				},
			}

			expectedPortBindings, err := constructPortBindings(providedReq.Ports)
			if err != nil {
				panic(err)
			}

			expectedContainerConfig := constructContainerConfig(
				providedReq.Cmd,
				providedReq.EnvVars,
				*providedReq.Image.Ref,
				expectedPortBindings,
				providedReq.WorkDir,
			)

			expectedHostConfig := &container.HostConfig{
				AutoRemove: true,
				Mounts: []mount.Mount{
					{
						Type:          "bind",
						Target:        "file1ContainerPath",
						Source:        "file1HostPath",
						Consistency:   "cached",
						ReadOnly:      false,
						BindOptions:   nil,
						VolumeOptions: nil,
						TmpfsOptions:  nil,
					},
					{
						Type:          "bind",
						Source:        "dir1HostPath",
						Target:        "dir1ContainerPath",
						Consistency:   "cached",
						ReadOnly:      false,
						BindOptions:   nil,
						TmpfsOptions:  nil,
						VolumeOptions: nil,
					},
					{
						Type:          "bind",
						Source:        "/unixSocket1HostAddress",
						Target:        "/unixSocket1ContainerAddress",
						ReadOnly:      false,
						Consistency:   "",
						BindOptions:   nil,
						VolumeOptions: nil,
						TmpfsOptions:  nil,
					},
				},
				PortBindings: nat.PortMap{
					"80/tcp": []nat.PortBinding{
						{
							HostPort: "80",
						},
					},
				},
				Privileged: true,
				Resources: container.Resources{
					DeviceRequests: []container.DeviceRequest{
						{
							Capabilities: [][]string{{"gpu"}},
							Count:        -1,
						},
					},
				},
			}

			expectedNetworkingConfig := &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					networkName: {
						Aliases: []string{
							*providedReq.Name,
						},
					},
				},
			}

			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			objectUnderTest := _runContainer{
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
			}

			/* act */
			objectUnderTest.RunContainer(
				providedCtx,
				providedReq,
				"rootCallID",
				pubSub,
				nopWriteCloser{io.Discard},
				nopWriteCloser{io.Discard},
			)

			/* assert */
			actualCtx,
				actualContainerConfig,
				actualHostConfig,
				actualNetworkingConfig,
				actualPlatformConfig,
				actualContainerName := fakeDockerClient.ContainerCreateArgsForCall(1)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualContainerConfig).To(Equal(expectedContainerConfig))
			Expect(*actualHostConfig).To(Equal(*expectedHostConfig))
			Expect(actualNetworkingConfig).To(Equal(expectedNetworkingConfig))
			Expect(actualPlatformConfig).To(BeNil())
			Expect(actualContainerName).To(Equal(fmt.Sprintf("opctl_%s", providedReq.ContainerID)))
		})
	})
})
