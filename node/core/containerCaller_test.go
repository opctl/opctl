package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var _ = Context("containerCaller", func() {
	Context("newContainerCaller", func() {
		It("should return containerCaller", func() {
			/* arrange/act/assert */
			Expect(newContainerCaller(
				new(containerprovider.Fake),
				new(pubsub.Fake),
				new(fakeDCGNodeRepo),
			)).Should(Not(BeNil()))
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
			fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
				close(eventChannel)
			}

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:        providedContainerId,
				PkgRef:    providedPkgRef,
				RootOpId:  providedRootOpId,
				Container: &dcgContainerDescriptor{},
			}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			objectUnderTest := newContainerCaller(
				new(containerprovider.Fake),
				fakePubSub,
				fakeDCGNodeRepo,
			)

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
			fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
				close(eventChannel)
			}

			objectUnderTest := newContainerCaller(
				new(containerprovider.Fake),
				fakePubSub,
				new(fakeDCGNodeRepo),
			)

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

			// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
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

			expectedReq, _ := constructDCGContainerCall(
				providedInboundScope,
				providedSCGContainerCall,
				providedContainerId,
				providedRootOpId,
				providedPkgRef,
			)

			fakeContainerProvider := new(containerprovider.Fake)

			fakePubSub := new(pubsub.Fake)
			fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
				close(eventChannel)
			}

			objectUnderTest := newContainerCaller(
				fakeContainerProvider,
				fakePubSub,
				new(fakeDCGNodeRepo),
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerId,
				providedSCGContainerCall,
				providedPkgRef,
				providedRootOpId,
			)

			/* assert */
			actualReq, actualEventPublisher := fakeContainerProvider.RunContainerArgsForCall(0)
			Expect(actualReq).To(Equal(expectedReq))
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
				fakeContainerProvider.RunContainerReturns(expectedError)

				objectUnderTest := newContainerCaller(
					fakeContainerProvider,
					new(pubsub.Fake),
					new(fakeDCGNodeRepo),
				)

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

		fakePubSub := new(pubsub.Fake)
		fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
			close(eventChannel)
		}

		fakeDCGNodeRepo := new(fakeDCGNodeRepo)

		objectUnderTest := newContainerCaller(
			new(containerprovider.Fake),
			fakePubSub,
			fakeDCGNodeRepo,
		)

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
		fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
			close(eventChannel)
		}

		objectUnderTest := newContainerCaller(
			new(containerprovider.Fake),
			fakePubSub,
			new(fakeDCGNodeRepo),
		)

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

		// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
		Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
		// set temporal fields to expected vals since they're already asserted
		actualEvent.Timestamp = expectedEvent.Timestamp

		Expect(actualEvent).To(Equal(expectedEvent))
	})
	XIt("should call pubSub.Publish w/ expected OutputInitializedEvents", func() {
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
					Name:     providedSCGContainerCall.Sockets["0.0.0.0"],
					Value:    &model.Data{Socket: providedContainerId},
					RootOpId: providedRootOpId,
					CallId:   providedContainerId,
				},
			},
		}

		fakePubSub := new(pubsub.Fake)

		// record actual published events
		var actualOutputInitEventMutex sync.Mutex
		actualOutputInitEvents := []*model.Event{}
		fakePubSub.PublishStub = func(event *model.Event) {
			event.Timestamp = expectedEventTimestamp
			actualOutputInitEventMutex.Lock()
			actualOutputInitEvents = append(actualOutputInitEvents, event)
			actualOutputInitEventMutex.Unlock()
		}
		fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
			close(eventChannel)
		}

		objectUnderTest := newContainerCaller(
			new(containerprovider.Fake),
			fakePubSub,
			new(fakeDCGNodeRepo),
		)

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
			Eventually(func() []*model.Event {
				return actualOutputInitEvents
			}, 5*time.Second, 1*time.Second).Should(ContainElement(expectedOutputInitEvent))
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
		fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
			close(eventChannel)
		}

		objectUnderTest := newContainerCaller(
			new(containerprovider.Fake),
			fakePubSub,
			new(fakeDCGNodeRepo),
		)

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

		// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
		Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
		// set temporal fields to expected vals since they're already asserted
		actualEvent.Timestamp = expectedEvent.Timestamp

		Expect(actualEvent).To(Equal(expectedEvent))
	})
})
