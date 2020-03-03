package docker

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker/internal/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("RunContainer", func() {
	closedContainerWaitOkBodyChan := make(chan container.ContainerWaitOKBody)
	close(closedContainerWaitOkBodyChan)

	It("should call dockerClient.ContainerRemove w/ expected args", func() {
		/* arrange */
		providedReq := &model.DCGContainerCall{
			ContainerID: "containerID",
			DCGBaseCall: model.DCGBaseCall{},
			Image:       &model.DCGContainerCallImage{Ref: new(string)},
		}

		expectedContainerRemoveOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}

		fakeDockerClient := new(FakeCommonAPIClient)
		fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

		fakePortBindingsFactory := new(FakePortBindingsFactory)
		// err to trigger immediate return
		fakePortBindingsFactory.ConstructReturns(nil, errors.New("dummyError"))

		objectUnderTest := _runContainer{
			containerStdErrStreamer: new(FakeContainerLogStreamer),
			containerStdOutStreamer: new(FakeContainerLogStreamer),
			dockerClient:            fakeDockerClient,
			imagePuller:             new(FakeImagePuller),
			portBindingsFactory:     fakePortBindingsFactory,
		}

		/* act */
		objectUnderTest.RunContainer(
			context.Background(),
			providedReq,
			new(FakeEventPublisher),
			nopWriteCloser{ioutil.Discard},
			nopWriteCloser{ioutil.Discard},
		)

		/* assert */
		_, actualContainerID, actualContainerRemoveOptions := fakeDockerClient.ContainerRemoveArgsForCall(0)
		Expect(actualContainerID).To(Equal(providedReq.ContainerID))
		Expect(actualContainerRemoveOptions).To(Equal(expectedContainerRemoveOptions))

	})
	It("should call portBindingsFactory.Construct w expected args", func() {
		/* arrange */
		providedReq := &model.DCGContainerCall{
			DCGBaseCall: model.DCGBaseCall{},
			Image:       &model.DCGContainerCallImage{Ref: new(string)},
			Ports: map[string]string{
				"6060/udp":  "6060",
				"8080-8081": "9090-9091",
			},
		}

		fakeDockerClient := new(FakeCommonAPIClient)
		fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

		fakePortBindingsFactory := new(FakePortBindingsFactory)

		objectUnderTest := _runContainer{
			containerConfigFactory:  new(FakeContainerConfigFactory),
			containerStdErrStreamer: new(FakeContainerLogStreamer),
			containerStdOutStreamer: new(FakeContainerLogStreamer),
			dockerClient:            fakeDockerClient,
			hostConfigFactory:       new(FakeHostConfigFactory),
			imagePuller:             new(FakeImagePuller),
			portBindingsFactory:     fakePortBindingsFactory,
		}

		/* act */
		objectUnderTest.RunContainer(
			context.Background(),
			providedReq,
			new(FakeEventPublisher),
			nopWriteCloser{ioutil.Discard},
			nopWriteCloser{ioutil.Discard},
		)

		/* assert */
		actualPorts := fakePortBindingsFactory.ConstructArgsForCall(0)
		Expect(actualPorts).To(Equal(providedReq.Ports))
	})
	Context("portBindingsFactory.Construct errs", func() {
		It("should return expected result", func() {
			/* arrange */
			fakePortBindingsFactory := new(FakePortBindingsFactory)
			expectedErr := errors.New("dummyErr")
			fakePortBindingsFactory.ConstructReturnsOnCall(0, nil, expectedErr)

			objectUnderTest := _runContainer{
				containerConfigFactory:  new(FakeContainerConfigFactory),
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            new(FakeCommonAPIClient),
				hostConfigFactory:       new(FakeHostConfigFactory),
				imagePuller:             new(FakeImagePuller),
				portBindingsFactory:     fakePortBindingsFactory,
			}

			/* act */
			_, actualErr := objectUnderTest.RunContainer(
				context.Background(),
				&model.DCGContainerCall{
					Image: &model.DCGContainerCallImage{Ref: new(string)},
				},
				new(FakeEventPublisher),
				nopWriteCloser{ioutil.Discard},
				nopWriteCloser{ioutil.Discard},
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
	})
	Context("portBindingsFactory.Construct doesn't err", func() {
		It("should call containerConfigFactory.Construct w expected args", func() {
			/* arrange */
			providedReq := &model.DCGContainerCall{
				DCGBaseCall: model.DCGBaseCall{},
				Cmd:         []string{"dummyCmd"},
				EnvVars: map[string]string{
					"envVar1Name": "envVar1Value",
					"envVar2Name": "envVar2Value",
					"envVar3Name": "envVar3Value",
				},
				Image:   &model.DCGContainerCallImage{Ref: new(string)},
				WorkDir: "dummyWorkDir",
			}

			fakeContainerConfigFactory := new(FakeContainerConfigFactory)

			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			portBindings := nat.PortMap{"80/tcp": []nat.PortBinding{}}
			fakePortBindingsFactory := new(FakePortBindingsFactory)
			fakePortBindingsFactory.ConstructReturns(portBindings, nil)

			objectUnderTest := _runContainer{
				containerConfigFactory:  fakeContainerConfigFactory,
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				hostConfigFactory:       new(FakeHostConfigFactory),
				imagePuller:             new(FakeImagePuller),
				portBindingsFactory:     fakePortBindingsFactory,
			}

			/* act */
			objectUnderTest.RunContainer(
				context.Background(),
				providedReq,
				new(FakeEventPublisher),
				nopWriteCloser{ioutil.Discard},
				nopWriteCloser{ioutil.Discard},
			)

			/* assert */
			actualCmd,
				actualEnvVars,
				actualImage,
				actualPortBindings,
				actualWorkDir := fakeContainerConfigFactory.ConstructArgsForCall(0)

			Expect(actualCmd).To(Equal(providedReq.Cmd))
			Expect(actualEnvVars).To(Equal(providedReq.EnvVars))
			Expect(actualImage).To(Equal(*providedReq.Image.Ref))
			Expect(actualPortBindings).To(Equal(portBindings))
			Expect(actualWorkDir).To(Equal(providedReq.WorkDir))
		})

		It("should call hostConfigFactory.Construct w expected args", func() {
			/* arrange */
			providedReq := &model.DCGContainerCall{
				DCGBaseCall: model.DCGBaseCall{},
				Dirs: map[string]string{
					"dir1ContainerPath": "dir1HostPath",
					"dir2ContainerPath": "dir2HostPath",
				},
				Files: map[string]string{
					"file1ContainerPath": "file1HostPath",
					"file2ContainerPath": "file2HostPath",
				},
				Image: &model.DCGContainerCallImage{Ref: new(string)},
				Sockets: map[string]string{
					"/unixSocket1ContainerAddress": "/unixSocket1HostAddress",
					"/unixSocket2ContainerAddress": "/unixSocket2HostAddress",
				},
			}

			portBindings := nat.PortMap{"80/tcp": []nat.PortBinding{}}
			fakePortBindingsFactory := new(FakePortBindingsFactory)
			fakePortBindingsFactory.ConstructReturns(portBindings, nil)

			fakeHostConfigFactory := new(FakeHostConfigFactory)

			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			objectUnderTest := _runContainer{
				containerConfigFactory:  new(FakeContainerConfigFactory),
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				hostConfigFactory:       fakeHostConfigFactory,
				imagePuller:             new(FakeImagePuller),
				portBindingsFactory:     fakePortBindingsFactory,
			}

			/* act */
			objectUnderTest.RunContainer(
				context.Background(),
				providedReq,
				new(FakeEventPublisher),
				nopWriteCloser{ioutil.Discard},
				nopWriteCloser{ioutil.Discard},
			)

			/* assert */
			actualDirs,
				actualFiles,
				actualSockets,
				actualPortBindings := fakeHostConfigFactory.ConstructArgsForCall(0)
			Expect(actualDirs).To(Equal(providedReq.Dirs))
			Expect(actualFiles).To(Equal(providedReq.Files))
			Expect(actualSockets).To(Equal(providedReq.Sockets))
			Expect(actualPortBindings).To(Equal(portBindings))
		})

		It("should call imagePuller.Pull w/ expected args", func() {

			/* arrange */
			providedCtx := context.Background()
			providedReq := &model.DCGContainerCall{
				DCGBaseCall: model.DCGBaseCall{
					RootOpID: "dummyRootOpID",
				},
				ContainerID: "dummyContainerID",
				Image:       &model.DCGContainerCallImage{Ref: new(string)},
			}

			providedEventPublisher := new(FakeEventPublisher)

			fakeImagePuller := new(FakeImagePuller)

			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			objectUnderTest := _runContainer{
				containerConfigFactory:  new(FakeContainerConfigFactory),
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				hostConfigFactory:       new(FakeHostConfigFactory),
				imagePuller:             fakeImagePuller,
				portBindingsFactory:     new(FakePortBindingsFactory),
			}

			/* act */
			objectUnderTest.RunContainer(
				providedCtx,
				providedReq,
				providedEventPublisher,
				nopWriteCloser{ioutil.Discard},
				nopWriteCloser{ioutil.Discard},
			)

			/* assert */
			actualCtx,
				actualContainerID,
				actualImagePullCreds,
				actualImageRef,
				actualRootOpID,
				actualEventPublisher := fakeImagePuller.PullArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualContainerID).To(Equal(providedReq.ContainerID))
			Expect(actualImagePullCreds).To(Equal(providedReq.Image.PullCreds))
			Expect(actualImageRef).To(Equal(*providedReq.Image.Ref))
			Expect(actualRootOpID).To(Equal(providedReq.RootOpID))
			Expect(actualEventPublisher).To(Equal(providedEventPublisher))
		})

		It("should call dockerClient.ContainerCreate w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedReq := &model.DCGContainerCall{
				DCGBaseCall: model.DCGBaseCall{
					RootOpID: "dummyRootOpID",
				},
				ContainerID: "dummyContainerID",
				Image:       &model.DCGContainerCallImage{Ref: new(string)},
				Name:        new(string),
			}

			fakeContainerConfigFactory := new(FakeContainerConfigFactory)
			expectedContainerConfig := &container.Config{}
			fakeContainerConfigFactory.ConstructReturns(expectedContainerConfig)

			fakeHostConfigFactory := new(FakeHostConfigFactory)
			expectedHostConfig := &container.HostConfig{}
			fakeHostConfigFactory.ConstructReturns(expectedHostConfig)

			expectedNetworkingConfig := &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					dockerNetworkName: {
						Aliases: []string{
							*providedReq.Name,
						},
					},
				},
			}

			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			objectUnderTest := _runContainer{
				containerConfigFactory:  fakeContainerConfigFactory,
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				hostConfigFactory:       fakeHostConfigFactory,
				imagePuller:             new(FakeImagePuller),
				portBindingsFactory:     new(FakePortBindingsFactory),
			}

			/* act */
			objectUnderTest.RunContainer(
				providedCtx,
				providedReq,
				new(FakeEventPublisher),
				nopWriteCloser{ioutil.Discard},
				nopWriteCloser{ioutil.Discard},
			)

			/* assert */
			actualCtx,
				actualContainerConfig,
				actualHostConfig,
				actualNetworkingConfig,
				actualContainerName := fakeDockerClient.ContainerCreateArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualContainerConfig).To(Equal(expectedContainerConfig))
			Expect(actualHostConfig).To(Equal(expectedHostConfig))
			Expect(actualNetworkingConfig).To(Equal(expectedNetworkingConfig))
			Expect(actualContainerName).To(Equal(providedReq.ContainerID))
		})
	})
})
