package docker

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"io"
	"io/ioutil"
)

var _ = Context("RunContainer", func() {
	closedContainerWaitOkBodyChan := make(chan container.ContainerWaitOKBody)
	close(closedContainerWaitOkBodyChan)

	It("should call portBindingsFactory.Construct w expected args", func() {
		/* arrange */
		providedReq := &model.DCGContainerCall{
			DCGBaseCall: model.DCGBaseCall{},
			Image:       &model.DCGContainerCallImage{},
			Ports: map[string]string{
				"6060/udp":  "6060",
				"8080-8081": "9090-9091",
			},
		}

		fakeDockerClient := new(fakeDockerClient)
		fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

		fakePortBindingsFactory := new(fakePortBindingsFactory)

		objectUnderTest := _runContainer{
			containerConfigFactory:  new(fakeContainerConfigFactory),
			containerStdErrStreamer: new(fakeContainerLogStreamer),
			containerStdOutStreamer: new(fakeContainerLogStreamer),
			dockerClient:            fakeDockerClient,
			hostConfigFactory:       new(fakeHostConfigFactory),
			imagePuller:             new(fakeImagePuller),
			portBindingsFactory:     fakePortBindingsFactory,
		}

		/* act */
		objectUnderTest.RunContainer(
			context.Background(),
			providedReq,
			new(pubsub.FakeEventPublisher),
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
			fakePortBindingsFactory := new(fakePortBindingsFactory)
			expectedErr := errors.New("dummyErr")
			fakePortBindingsFactory.ConstructReturnsOnCall(0, nil, expectedErr)

			objectUnderTest := _runContainer{
				containerConfigFactory:  new(fakeContainerConfigFactory),
				containerStdErrStreamer: new(fakeContainerLogStreamer),
				containerStdOutStreamer: new(fakeContainerLogStreamer),
				dockerClient:            new(fakeDockerClient),
				hostConfigFactory:       new(fakeHostConfigFactory),
				imagePuller:             new(fakeImagePuller),
				portBindingsFactory:     fakePortBindingsFactory,
			}

			/* act */
			_, actualErr := objectUnderTest.RunContainer(
				context.Background(),
				&model.DCGContainerCall{},
				new(pubsub.FakeEventPublisher),
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
				Image:   &model.DCGContainerCallImage{Ref: "dummyImage"},
				WorkDir: "dummyWorkDir",
			}

			fakeContainerConfigFactory := new(fakeContainerConfigFactory)

			fakeDockerClient := new(fakeDockerClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			portBindings := nat.PortMap{"80/tcp": []nat.PortBinding{}}
			fakePortBindingsFactory := new(fakePortBindingsFactory)
			fakePortBindingsFactory.ConstructReturns(portBindings, nil)

			objectUnderTest := _runContainer{
				containerConfigFactory:  fakeContainerConfigFactory,
				containerStdErrStreamer: new(fakeContainerLogStreamer),
				containerStdOutStreamer: new(fakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				hostConfigFactory:       new(fakeHostConfigFactory),
				imagePuller:             new(fakeImagePuller),
				portBindingsFactory:     fakePortBindingsFactory,
			}

			/* act */
			objectUnderTest.RunContainer(
				context.Background(),
				providedReq,
				new(pubsub.FakeEventPublisher),
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
			Expect(actualImage).To(Equal(providedReq.Image.Ref))
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
				Image: &model.DCGContainerCallImage{},
				Sockets: map[string]string{
					"/unixSocket1ContainerAddress": "/unixSocket1HostAddress",
					"/unixSocket2ContainerAddress": "/unixSocket2HostAddress",
				},
			}

			portBindings := nat.PortMap{"80/tcp": []nat.PortBinding{}}
			fakePortBindingsFactory := new(fakePortBindingsFactory)
			fakePortBindingsFactory.ConstructReturns(portBindings, nil)

			fakeHostConfigFactory := new(fakeHostConfigFactory)

			fakeDockerClient := new(fakeDockerClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			objectUnderTest := _runContainer{
				containerConfigFactory:  new(fakeContainerConfigFactory),
				containerStdErrStreamer: new(fakeContainerLogStreamer),
				containerStdOutStreamer: new(fakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				hostConfigFactory:       fakeHostConfigFactory,
				imagePuller:             new(fakeImagePuller),
				portBindingsFactory:     fakePortBindingsFactory,
			}

			/* act */
			objectUnderTest.RunContainer(
				context.Background(),
				providedReq,
				new(pubsub.FakeEventPublisher),
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
				Image:       &model.DCGContainerCallImage{Ref: "dummyImage"},
			}

			providedEventPublisher := new(pubsub.FakeEventPublisher)

			fakeImagePuller := new(fakeImagePuller)

			fakeDockerClient := new(fakeDockerClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			objectUnderTest := _runContainer{
				containerConfigFactory:  new(fakeContainerConfigFactory),
				containerStdErrStreamer: new(fakeContainerLogStreamer),
				containerStdOutStreamer: new(fakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				hostConfigFactory:       new(fakeHostConfigFactory),
				imagePuller:             fakeImagePuller,
				portBindingsFactory:     new(fakePortBindingsFactory),
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
				actualImage,
				actualContainerID,
				actualRootOpID,
				actualEventPublisher := fakeImagePuller.PullArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualImage).To(Equal(providedReq.Image))
			Expect(actualContainerID).To(Equal(providedReq.ContainerID))
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
				Image:       &model.DCGContainerCallImage{},
				Name:        new(string),
			}

			fakeContainerConfigFactory := new(fakeContainerConfigFactory)
			expectedContainerConfig := &container.Config{}
			fakeContainerConfigFactory.ConstructReturns(expectedContainerConfig)

			fakeHostConfigFactory := new(fakeHostConfigFactory)
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

			fakeDockerClient := new(fakeDockerClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			objectUnderTest := _runContainer{
				containerConfigFactory:  fakeContainerConfigFactory,
				containerStdErrStreamer: new(fakeContainerLogStreamer),
				containerStdOutStreamer: new(fakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				hostConfigFactory:       fakeHostConfigFactory,
				imagePuller:             new(fakeImagePuller),
				portBindingsFactory:     new(fakePortBindingsFactory),
			}

			/* act */
			objectUnderTest.RunContainer(
				providedCtx,
				providedReq,
				new(pubsub.FakeEventPublisher),
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

type nopWriteCloser struct {
	io.Writer
}

func (w nopWriteCloser) Close() error { return nil }
