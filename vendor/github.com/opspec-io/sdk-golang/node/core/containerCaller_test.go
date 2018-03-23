package core

import (
	"errors"
	"github.com/golang-interfaces/iio"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/core/containerruntime"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall"
	"github.com/opspec-io/sdk-golang/util/pubsub"
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
				new(containerruntime.Fake),
				new(containercall.FakeInterpreter),
				new(pubsub.Fake),
				new(fakeDCGNodeRepo),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call dcgNodeRepo.Add w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Value{}
			providedContainerID := "dummyContainerID"
			providedSCGContainerCall := &model.SCGContainerCall{}
			providedOpHandle := new(data.FakeHandle)
			providedRootOpID := "dummyRootOpID"

			fakePubSub := new(pubsub.Fake)

			fakeContainerCallInterpreter := new(containercall.FakeInterpreter)
			// error to trigger immediate return
			fakeContainerCallInterpreter.InterpretReturns(nil, errors.New("dummyError"))

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:        providedContainerID,
				OpRef:     providedOpHandle.Ref(),
				RootOpID:  providedRootOpID,
				Container: &dcgContainerDescriptor{},
			}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerRuntime: new(containerruntime.Fake),
				containerCall:    fakeContainerCallInterpreter,
				pubSub:           fakePubSub,
				dcgNodeRepo:      fakeDCGNodeRepo,
				io:               fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerID,
				providedSCGContainerCall,
				providedOpHandle,
				providedRootOpID,
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
		It("should call pubSub.Publish w/ expected ContainerStartedEvent", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Value{}
			providedContainerID := "dummyContainerID"
			providedSCGContainerCall := &model.SCGContainerCall{}
			providedOpHandle := new(data.FakeHandle)
			providedRootOpID := "dummyRootOpID"

			expectedEvent := model.Event{
				Timestamp: time.Now().UTC(),
				ContainerStarted: &model.ContainerStartedEvent{
					ContainerID: providedContainerID,
					OpRef:       providedOpHandle.Ref(),
					RootOpID:    providedRootOpID,
				},
			}

			fakeContainerCallInterpreter := new(containercall.FakeInterpreter)
			fakeContainerCallInterpreter.InterpretReturns(&model.DCGContainerCall{Image: &model.DCGContainerCallImage{}}, nil)

			fakePubSub := new(pubsub.Fake)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerRuntime: new(containerruntime.Fake),
				containerCall:    fakeContainerCallInterpreter,
				pubSub:           fakePubSub,
				dcgNodeRepo:      new(fakeDCGNodeRepo),
				io:               fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerID,
				providedSCGContainerCall,
				providedOpHandle,
				providedRootOpID,
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
			expectedDCGContainerCall := &model.DCGContainerCall{}

			fakeContainerCallInterpreter := new(containercall.FakeInterpreter)
			fakeContainerCallInterpreter.InterpretReturns(expectedDCGContainerCall, nil)

			fakeContainerRuntime := new(containerruntime.Fake)

			fakePubSub := new(pubsub.Fake)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerRuntime: fakeContainerRuntime,
				containerCall:    fakeContainerCallInterpreter,
				pubSub:           fakePubSub,
				dcgNodeRepo:      new(fakeDCGNodeRepo),
				io:               fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				map[string]*model.Value{},
				"dummyContainerID",
				&model.SCGContainerCall{},
				new(data.FakeHandle),
				"dummyRootOpID",
			)

			/* assert */
			_, actualDCGContainerCall, actualEventPublisher, _, _ := fakeContainerRuntime.RunContainerArgsForCall(0)
			Expect(actualDCGContainerCall).To(Equal(expectedDCGContainerCall))
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
					containerCall:    new(containercall.FakeInterpreter),
					pubSub:           new(pubsub.Fake),
					dcgNodeRepo:      new(fakeDCGNodeRepo),
					io:               fakeIIO,
				}

				/* act */
				actualError := objectUnderTest.Call(
					map[string]*model.Value{},
					"dummyContainerID",
					&model.SCGContainerCall{},
					new(data.FakeHandle),
					"dummyRootOpID",
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
	It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
		/* arrange */
		providedContainerID := "dummyContainerID"

		fakeDCGNodeRepo := new(fakeDCGNodeRepo)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(containerruntime.Fake),
			containerCall:    new(containercall.FakeInterpreter),
			pubSub:           new(pubsub.Fake),
			dcgNodeRepo:      fakeDCGNodeRepo,
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			map[string]*model.Value{},
			providedContainerID,
			&model.SCGContainerCall{},
			new(data.FakeHandle),
			"dummyRootOpID",
		)

		/* assert */
		Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedContainerID))
	})

	It("should call pubSub.Publish w/ expected ContainerExitedEvent", func() {
		/* arrange */
		providedInboundScope := map[string]*model.Value{}
		providedContainerID := "dummyContainerID"
		providedSCGContainerCall := &model.SCGContainerCall{}
		providedOpHandle := new(data.FakeHandle)
		providedRootOpID := "dummyRootOpID"

		expectedEvent := model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerID: providedContainerID,
				OpRef:       providedOpHandle.Ref(),
				RootOpID:    providedRootOpID,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(containerruntime.Fake),
			containerCall:    new(containercall.FakeInterpreter),
			pubSub:           fakePubSub,
			dcgNodeRepo:      new(fakeDCGNodeRepo),
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			providedInboundScope,
			providedContainerID,
			providedSCGContainerCall,
			providedOpHandle,
			providedRootOpID,
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
		providedInboundScope := map[string]*model.Value{}
		providedContainerID := "dummyContainerID"
		providedSCGContainerCall := &model.SCGContainerCall{}
		providedOpHandle := new(data.FakeHandle)
		providedRootOpID := "dummyRootOpID"

		expectedEvent := model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerID: providedContainerID,
				OpRef:       providedOpHandle.Ref(),
				RootOpID:    providedRootOpID,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(containerruntime.Fake),
			containerCall:    new(containercall.FakeInterpreter),
			pubSub:           fakePubSub,
			dcgNodeRepo:      new(fakeDCGNodeRepo),
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			providedInboundScope,
			providedContainerID,
			providedSCGContainerCall,
			providedOpHandle,
			providedRootOpID,
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
