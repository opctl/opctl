package core

import (
	"github.com/golang-interfaces/iio"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
	"io"
	"time"
)

var _ = Context("containerCaller", func() {
	closedPipeReader, closedPipeWriter := io.Pipe()
	closedPipeReader.Close()
	closedPipeWriter.Close()
	Context("newContainerCaller", func() {
		It("should return containerCaller", func() {
			/* arrange/act/assert */
			Expect(newContainerCaller(
				new(containerprovider.Fake),
				new(fakeDCGFactory),
				new(pubsub.Fake),
				new(fakeDCGNodeRepo),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call dcgNodeRepo.Add w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedContainerId := "dummyContainerId"
			providedSCGContainerCall := &model.SCGContainerCall{}
			providedPkgRef := "dummyPkgRef"
			providedRootOpId := "dummyRootOpId"

			fakePubSub := new(pubsub.Fake)

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:        providedContainerId,
				PkgRef:    providedPkgRef,
				RootOpId:  providedRootOpId,
				Container: &dcgContainerDescriptor{},
			}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerProvider: new(containerprovider.Fake),
				dcgFactory:        new(fakeDCGFactory),
				pubSub:            fakePubSub,
				dcgNodeRepo:       fakeDCGNodeRepo,
				io:                fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerId,
				providedSCGContainerCall,
				providedPkgRef,
				providedRootOpId,
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
		It("should call pubSub.Publish w/ expected ContainerStartedEvent", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedContainerId := "dummyContainerId"
			providedSCGContainerCall := &model.SCGContainerCall{}
			providedPkgRef := "dummyPkgRef"
			providedRootOpId := "dummyRootOpId"

			expectedEvent := &model.Event{
				Timestamp: time.Now().UTC(),
				ContainerStarted: &model.ContainerStartedEvent{
					ContainerId: providedContainerId,
					PkgRef:      providedPkgRef,
					RootOpId:    providedRootOpId,
				},
			}

			fakePubSub := new(pubsub.Fake)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerProvider: new(containerprovider.Fake),
				dcgFactory:        new(fakeDCGFactory),
				pubSub:            fakePubSub,
				dcgNodeRepo:       new(fakeDCGNodeRepo),
				io:                fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerId,
				providedSCGContainerCall,
				providedPkgRef,
				providedRootOpId,
			)

			/* assert */
			actualEvent := fakePubSub.PublishArgsForCall(0)

			// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
			Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
			// set temporal fields to expected vals since they're already asserted
			actualEvent.Timestamp = expectedEvent.Timestamp

			Expect(actualEvent).To(Equal(expectedEvent))
		})
		It("should call containerProvider.RunContainer w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedContainerId := "dummyContainerId"
			providedSCGContainerCall := &model.SCGContainerCall{}
			providedPkgRef := "dummyPkgRef"
			providedRootOpId := "dummyRootOpId"

			expectedDCGContainerCall := &model.DCGContainerCall{}

			fakeDCGFactory := new(fakeDCGFactory)
			fakeDCGFactory.ConstructReturns(expectedDCGContainerCall, nil)

			fakeContainerProvider := new(containerprovider.Fake)

			fakePubSub := new(pubsub.Fake)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerProvider: fakeContainerProvider,
				dcgFactory:        fakeDCGFactory,
				pubSub:            fakePubSub,
				dcgNodeRepo:       new(fakeDCGNodeRepo),
				io:                fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerId,
				providedSCGContainerCall,
				providedPkgRef,
				providedRootOpId,
			)

			/* assert */
			actualDCGContainerCall, actualEventPublisher, _, _ := fakeContainerProvider.RunContainerArgsForCall(0)
			Expect(actualDCGContainerCall).To(Equal(expectedDCGContainerCall))
			Expect(actualEventPublisher).To(Equal(fakePubSub))
		})
		Context("containerProvider.RunContainer errors", func() {
			It("should return expected error", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedContainerId := "dummyContainerId"
				providedSCGContainerCall := &model.SCGContainerCall{}
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"

				expectedError := errors.New("dummyError")

				fakeContainerProvider := new(containerprovider.Fake)
				fakeContainerProvider.RunContainerReturns(nil, expectedError)

				fakeIIO := new(iio.Fake)
				fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

				objectUnderTest := _containerCaller{
					containerProvider: fakeContainerProvider,
					dcgFactory:        new(fakeDCGFactory),
					pubSub:            new(pubsub.Fake),
					dcgNodeRepo:       new(fakeDCGNodeRepo),
					io:                fakeIIO,
				}

				/* act */
				actualError := objectUnderTest.Call(
					providedInboundScope,
					providedContainerId,
					providedSCGContainerCall,
					providedPkgRef,
					providedRootOpId,
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
	It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
		/* arrange */
		providedInboundScope := map[string]*model.Data{}
		providedContainerId := "dummyContainerId"
		providedSCGContainerCall := &model.SCGContainerCall{}
		providedPkgRef := "dummyPkgRef"
		providedRootOpId := "dummyRootOpId"

		fakeDCGNodeRepo := new(fakeDCGNodeRepo)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerProvider: new(containerprovider.Fake),
			dcgFactory:        new(fakeDCGFactory),
			pubSub:            new(pubsub.Fake),
			dcgNodeRepo:       fakeDCGNodeRepo,
			io:                fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			providedInboundScope,
			providedContainerId,
			providedSCGContainerCall,
			providedPkgRef,
			providedRootOpId,
		)

		/* assert */
		Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedContainerId))
	})

	It("should call pubSub.Publish w/ expected ContainerExitedEvent", func() {
		/* arrange */
		providedInboundScope := map[string]*model.Data{}
		providedContainerId := "dummyContainerId"
		providedSCGContainerCall := &model.SCGContainerCall{}
		providedPkgRef := "dummyPkgRef"
		providedRootOpId := "dummyRootOpId"

		expectedEvent := &model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerId: providedContainerId,
				PkgRef:      providedPkgRef,
				RootOpId:    providedRootOpId,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerProvider: new(containerprovider.Fake),
			dcgFactory:        new(fakeDCGFactory),
			pubSub:            fakePubSub,
			dcgNodeRepo:       new(fakeDCGNodeRepo),
			io:                fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			providedInboundScope,
			providedContainerId,
			providedSCGContainerCall,
			providedPkgRef,
			providedRootOpId,
		)

		/* assert */
		actualEvent := fakePubSub.PublishArgsForCall(1)

		// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
		Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
		// set temporal fields to expected vals since they're already asserted
		actualEvent.Timestamp = expectedEvent.Timestamp

		Expect(actualEvent).To(Equal(expectedEvent))
	})
	It("should call pubSub.Publish w/ expected OutputInitializedEvents", func() {
		/* arrange */
		providedInboundScope := map[string]*model.Data{}
		providedContainerId := "dummyContainerId"
		providedSCGContainerCall := &model.SCGContainerCall{
			Sockets: map[string]string{
				"0.0.0.0": "socket0Name",
			},
			Files:  map[string]string{},
			Dirs:   map[string]string{},
			StdErr: map[string]string{},
			StdOut: map[string]string{},
		}
		providedPkgRef := "dummyPkgRef"
		providedRootOpId := "dummyRootOpId"

		expectedEventTimestamp := time.Now().UTC()

		expectedOutputInitEvents := []*model.Event{
			{
				Timestamp: expectedEventTimestamp,
				OutputInitialized: &model.OutputInitializedEvent{
					CallId:   providedContainerId,
					RootOpId: providedRootOpId,
					Name:     providedSCGContainerCall.Sockets["0.0.0.0"],
					Value:    &model.Data{Socket: &providedContainerId},
				},
			},
		}

		fakeDCGFactory := new(fakeDCGFactory)
		fakeDCGFactory.ConstructReturns(
			&model.DCGContainerCall{
				DCGBaseCall: &model.DCGBaseCall{
					RootOpId: providedRootOpId,
				},
				ContainerId: providedContainerId,
			},
			nil,
		)

		fakePubSub := new(pubsub.Fake)

		// record actual published events
		actualOutputInitEvents := make(chan *model.Event, 1)
		fakePubSub.PublishStub = func(event *model.Event) {
			event.Timestamp = expectedEventTimestamp
			if event.OutputInitialized != nil {
				actualOutputInitEvents <- event
			}
		}

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerProvider: new(containerprovider.Fake),
			dcgFactory:        fakeDCGFactory,
			pubSub:            fakePubSub,
			dcgNodeRepo:       new(fakeDCGNodeRepo),
			io:                fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			providedInboundScope,
			providedContainerId,
			providedSCGContainerCall,
			providedPkgRef,
			providedRootOpId,
		)

		/* assert */
		for _, expectedOutputInitEvent := range expectedOutputInitEvents {
			Eventually(func() *model.Event {
				return <-actualOutputInitEvents
			}, 5*time.Second, 1*time.Second).Should(Equal(expectedOutputInitEvent))
		}
	})
	It("should call pubSub.Publish w/ expected ContainerExitedEvent", func() {
		/* arrange */
		providedInboundScope := map[string]*model.Data{}
		providedContainerId := "dummyContainerId"
		providedSCGContainerCall := &model.SCGContainerCall{}
		providedPkgRef := "dummyPkgRef"
		providedRootOpId := "dummyRootOpId"

		expectedEvent := &model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerId: providedContainerId,
				PkgRef:      providedPkgRef,
				RootOpId:    providedRootOpId,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerProvider: new(containerprovider.Fake),
			dcgFactory:        new(fakeDCGFactory),
			pubSub:            fakePubSub,
			dcgNodeRepo:       new(fakeDCGNodeRepo),
			io:                fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			providedInboundScope,
			providedContainerId,
			providedSCGContainerCall,
			providedPkgRef,
			providedRootOpId,
		)

		/* assert */
		actualEvent := fakePubSub.PublishArgsForCall(1)

		// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
		Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
		// set temporal fields to expected vals since they're already asserted
		actualEvent.Timestamp = expectedEvent.Timestamp

		Expect(actualEvent).To(Equal(expectedEvent))
	})
})
