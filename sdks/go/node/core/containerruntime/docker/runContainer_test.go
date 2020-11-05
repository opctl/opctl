package docker

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
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
			imagePuller:             new(FakeImagePuller),
			ensureNetworkExistser:   new(FakeEnsureNetworkExistser),
		}

		/* act */
		objectUnderTest.RunContainer(
			context.Background(),
			providedReq,
			"rootCallID",
			new(FakeEventPublisher),
			nopWriteCloser{ioutil.Discard},
			nopWriteCloser{ioutil.Discard},
		)

		/* assert */
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
				ensureNetworkExistser:   new(FakeEnsureNetworkExistser),
				hostConfigFactory:       new(FakeHostConfigFactory),
				imagePuller:             new(FakeImagePuller),
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
				new(FakeEventPublisher),
				nopWriteCloser{ioutil.Discard},
				nopWriteCloser{ioutil.Discard},
			)

			/* assert */
			Expect(actualErr).To(Equal(fmt.Errorf("Invalid containerPort: *")))
		})
	})
	Context("portBindingsFactory.Construct doesn't err", func() {

		It("should call hostConfigFactory.Construct w expected args", func() {
			/* arrange */
			providedReq := &model.ContainerCall{
				BaseCall: model.BaseCall{},
				Dirs: map[string]string{
					"dir1ContainerPath": "dir1HostPath",
					"dir2ContainerPath": "dir2HostPath",
				},
				Files: map[string]string{
					"file1ContainerPath": "file1HostPath",
					"file2ContainerPath": "file2HostPath",
				},
				Image: &model.ContainerCallImage{Ref: new(string)},
				Sockets: map[string]string{
					"/unixSocket1ContainerAddress": "/unixSocket1HostAddress",
					"/unixSocket2ContainerAddress": "/unixSocket2HostAddress",
				},
				Ports: map[string]string{
					"80": "80",
				},
			}

			portBindings, err := constructPortBindings(providedReq.Ports)
			if nil != err {
				panic(err)
			}

			fakeHostConfigFactory := new(FakeHostConfigFactory)

			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			objectUnderTest := _runContainer{
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				ensureNetworkExistser:   new(FakeEnsureNetworkExistser),
				hostConfigFactory:       fakeHostConfigFactory,
				imagePuller:             new(FakeImagePuller),
			}

			/* act */
			objectUnderTest.RunContainer(
				context.Background(),
				providedReq,
				"rootCallID",
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
			providedReq := &model.ContainerCall{
				BaseCall:    model.BaseCall{},
				ContainerID: "dummyContainerID",
				Image:       &model.ContainerCallImage{Ref: new(string)},
			}
			providedRootCallID := "providedRootCallID"

			providedEventPublisher := new(FakeEventPublisher)

			fakeImagePuller := new(FakeImagePuller)

			fakeDockerClient := new(FakeCommonAPIClient)
			fakeDockerClient.ContainerWaitReturns(closedContainerWaitOkBodyChan, nil)

			objectUnderTest := _runContainer{
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				ensureNetworkExistser:   new(FakeEnsureNetworkExistser),
				hostConfigFactory:       new(FakeHostConfigFactory),
				imagePuller:             fakeImagePuller,
			}

			/* act */
			objectUnderTest.RunContainer(
				providedCtx,
				providedReq,
				providedRootCallID,
				providedEventPublisher,
				nopWriteCloser{ioutil.Discard},
				nopWriteCloser{ioutil.Discard},
			)

			/* assert */
			actualCtx,
				actualContainerID,
				actualImagePullCreds,
				actualImageRef,
				actualRootCallID,
				actualEventPublisher := fakeImagePuller.PullArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualContainerID).To(Equal(providedReq.ContainerID))
			Expect(actualImagePullCreds).To(Equal(providedReq.Image.PullCreds))
			Expect(actualImageRef).To(Equal(*providedReq.Image.Ref))
			Expect(actualRootCallID).To(Equal(providedRootCallID))
			Expect(actualEventPublisher).To(Equal(providedEventPublisher))
		})

		It("should call dockerClient.ContainerCreate w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedReq := &model.ContainerCall{
				BaseCall:    model.BaseCall{},
				ContainerID: "dummyContainerID",
				Image:       &model.ContainerCallImage{Ref: new(string)},
				Name:        new(string),
			}

			expectedPortBindings, err := constructPortBindings(providedReq.Ports)
			if nil != err {
				panic(err)
			}

			expectedContainerConfig := constructContainerConfig(
				providedReq.Cmd,
				providedReq.EnvVars,
				*providedReq.Image.Ref,
				expectedPortBindings,
				providedReq.WorkDir,
			)

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
				containerStdErrStreamer: new(FakeContainerLogStreamer),
				containerStdOutStreamer: new(FakeContainerLogStreamer),
				dockerClient:            fakeDockerClient,
				ensureNetworkExistser:   new(FakeEnsureNetworkExistser),
				hostConfigFactory:       fakeHostConfigFactory,
				imagePuller:             new(FakeImagePuller),
			}

			/* act */
			objectUnderTest.RunContainer(
				providedCtx,
				providedReq,
				"rootCallID",
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
			Expect(actualContainerName).To(Equal(fmt.Sprintf("opctl_%s", providedReq.ContainerID)))
		})
	})
})
