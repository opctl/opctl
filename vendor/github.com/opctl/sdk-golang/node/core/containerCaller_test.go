package core

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/golang-interfaces/iio"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core/containerruntime"
	"github.com/opctl/sdk-golang/util/pubsub"
)

var _ = Context("containerCaller", func() {
	closedPipeReader, closedPipeWriter := io.Pipe()
	closedPipeReader.Close()
	closedPipeWriter.Close()
	opHandleRef := "dummyOpRef"
	fakeOpHandle := new(data.FakeHandle)
	fakeOpHandle.RefReturns(opHandleRef)

	Context("newContainerCaller", func() {
		It("should return containerCaller", func() {
			/* arrange/act/assert */
			Expect(newContainerCaller(
				new(containerruntime.Fake),
				new(pubsub.Fake),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call pubSub.Publish w/ expected ContainerStartedEvent", func() {
			/* arrange */
			providedDCGContainerCall := &model.DCGContainerCall{
				DCGBaseCall: model.DCGBaseCall{
					OpHandle: fakeOpHandle,
					RootOpID: "providedRootID",
				},
				ContainerID: "providedContainerID",
			}
			providedInboundScope := map[string]*model.Value{}
			providedSCGContainerCall := &model.SCGContainerCall{}

			expectedEvent := model.Event{
				Timestamp: time.Now().UTC(),
				ContainerStarted: &model.ContainerStartedEvent{
					ContainerID: providedDCGContainerCall.ContainerID,
					OpRef:       providedDCGContainerCall.OpHandle.Ref(),
					RootOpID:    providedDCGContainerCall.RootOpID,
				},
			}

			fakePubSub := new(pubsub.Fake)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerRuntime: new(containerruntime.Fake),
				pubSub:           fakePubSub,
				io:               fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				context.Background(),
				providedDCGContainerCall,
				providedInboundScope,
				providedSCGContainerCall,
			)

			/* assert */
			actualEvent := fakePubSub.PublishArgsForCall(0)

			// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
			Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
			// set temporal fields to expected vals since they're already asserted
			actualEvent.Timestamp = expectedEvent.Timestamp

			Expect(actualEvent).To(Equal(expectedEvent))
		})
		It("should call containerRuntime.RunContainer w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedDCGContainerCall := &model.DCGContainerCall{
				DCGBaseCall: model.DCGBaseCall{
					OpHandle: fakeOpHandle,
				},
			}
			fakeContainerRuntime := new(containerruntime.Fake)

			fakePubSub := new(pubsub.Fake)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerRuntime: fakeContainerRuntime,
				pubSub:           fakePubSub,
				io:               fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				providedCtx,
				providedDCGContainerCall,
				map[string]*model.Value{},
				&model.SCGContainerCall{},
			)

			/* assert */
			_,
				actualDCGContainerCall,
				actualEventPublisher,
				_,
				_ := fakeContainerRuntime.RunContainerArgsForCall(0)
			Expect(actualDCGContainerCall).To(Equal(providedDCGContainerCall))
			Expect(actualEventPublisher).To(Equal(fakePubSub))
		})
		Context("containerRuntime.RunContainer errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")

				fakeContainerRuntime := new(containerruntime.Fake)
				fakeContainerRuntime.RunContainerReturns(nil, expectedError)

				fakeIIO := new(iio.Fake)
				fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

				objectUnderTest := _containerCaller{
					containerRuntime: fakeContainerRuntime,
					pubSub:           new(pubsub.Fake),
					io:               fakeIIO,
				}

				/* act */
				actualError := objectUnderTest.Call(
					context.Background(),
					&model.DCGContainerCall{
						DCGBaseCall: model.DCGBaseCall{
							OpHandle: fakeOpHandle,
						},
					},
					map[string]*model.Value{},
					&model.SCGContainerCall{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})

	It("should call pubSub.Publish w/ expected ContainerExitedEvent", func() {
		/* arrange */
		providedDCGContainerCall := &model.DCGContainerCall{
			DCGBaseCall: model.DCGBaseCall{
				OpHandle: fakeOpHandle,
				RootOpID: "providedRootID",
			},
			ContainerID: "providedContainerID",
		}
		providedInboundScope := map[string]*model.Value{}
		providedSCGContainerCall := &model.SCGContainerCall{}

		expectedEvent := model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerID: providedDCGContainerCall.ContainerID,
				OpRef:       providedDCGContainerCall.OpHandle.Ref(),
				Outputs:     map[string]*model.Value{},
				RootOpID:    providedDCGContainerCall.RootOpID,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(containerruntime.Fake),
			pubSub:           fakePubSub,
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			context.Background(),
			providedDCGContainerCall,
			providedInboundScope,
			providedSCGContainerCall,
		)

		/* assert */
		actualEvent := fakePubSub.PublishArgsForCall(1)

		// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
		Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
		// set temporal fields to expected vals since they're already asserted
		actualEvent.Timestamp = expectedEvent.Timestamp

		Expect(actualEvent).To(Equal(expectedEvent))
	})
	It("should call pubSub.Publish w/ expected ContainerExitedEvent", func() {
		/* arrange */
		providedDCGContainerCall := &model.DCGContainerCall{
			DCGBaseCall: model.DCGBaseCall{
				OpHandle: fakeOpHandle,
				RootOpID: "providedRootID",
			},
			ContainerID: "providedContainerID",
		}
		providedInboundScope := map[string]*model.Value{}
		providedSCGContainerCall := &model.SCGContainerCall{}

		expectedEvent := model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerID: providedDCGContainerCall.ContainerID,
				OpRef:       providedDCGContainerCall.OpHandle.Ref(),
				Outputs:     map[string]*model.Value{},
				RootOpID:    providedDCGContainerCall.RootOpID,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(containerruntime.Fake),
			pubSub:           fakePubSub,
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			context.Background(),
			providedDCGContainerCall,
			providedInboundScope,
			providedSCGContainerCall,
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
