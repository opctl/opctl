package core

import (
	"errors"
	"github.com/golang-interfaces/iio"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core/containerruntime"
	"github.com/opctl/sdk-golang/util/pubsub"
	"io"
	"time"
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
				new(fakeDCGNodeRepo),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call dcgNodeRepo.Add w/ expected args", func() {
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

			fakePubSub := new(pubsub.Fake)

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:        providedDCGContainerCall.ContainerID,
				OpRef:     providedDCGContainerCall.OpHandle.Ref(),
				RootOpID:  providedDCGContainerCall.RootOpID,
				Container: &dcgContainerDescriptor{},
			}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			fakeContainerRuntime := new(containerruntime.Fake)
			// error to trigger immediate return
			fakeContainerRuntime.RunContainerReturns(nil, errors.New("dummyError"))

			objectUnderTest := _containerCaller{
				containerRuntime: fakeContainerRuntime,
				pubSub:           fakePubSub,
				dcgNodeRepo:      fakeDCGNodeRepo,
				io:               fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				providedDCGContainerCall,
				providedInboundScope,
				providedSCGContainerCall,
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
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
				dcgNodeRepo:      new(fakeDCGNodeRepo),
				io:               fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
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
				dcgNodeRepo:      new(fakeDCGNodeRepo),
				io:               fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
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
					dcgNodeRepo:      new(fakeDCGNodeRepo),
					io:               fakeIIO,
				}

				/* act */
				actualError := objectUnderTest.Call(
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
	It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
		/* arrange */
		providedDCGContainerCall := &model.DCGContainerCall{
			DCGBaseCall: model.DCGBaseCall{
				OpHandle: fakeOpHandle,
				RootOpID: "providedRootID",
			},
			ContainerID: "providedContainerID",
		}

		fakeDCGNodeRepo := new(fakeDCGNodeRepo)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(containerruntime.Fake),
			pubSub:           new(pubsub.Fake),
			dcgNodeRepo:      fakeDCGNodeRepo,
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			providedDCGContainerCall,
			map[string]*model.Value{},
			&model.SCGContainerCall{},
		)

		/* assert */
		Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedDCGContainerCall.ContainerID))
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
				RootOpID:    providedDCGContainerCall.RootOpID,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(containerruntime.Fake),
			pubSub:           fakePubSub,
			dcgNodeRepo:      new(fakeDCGNodeRepo),
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
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
				RootOpID:    providedDCGContainerCall.RootOpID,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(containerruntime.Fake),
			pubSub:           fakePubSub,
			dcgNodeRepo:      new(fakeDCGNodeRepo),
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
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
