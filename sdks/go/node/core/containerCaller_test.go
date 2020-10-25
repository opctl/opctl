package core

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/golang-interfaces/iio"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("containerCaller", func() {
	closedPipeReader, closedPipeWriter := io.Pipe()
	closedPipeReader.Close()
	closedPipeWriter.Close()

	Context("newContainerCaller", func() {
		It("should return containerCaller", func() {
			/* arrange/act/assert */
			Expect(newContainerCaller(
				new(FakeContainerRuntime),
				new(FakePubSub),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call pubSub.Publish w/ expected ContainerStarted", func() {
			/* arrange */
			providedOpPath := "providedOpPath"
			providedContainerCall := &model.ContainerCall{
				BaseCall: model.BaseCall{
					OpPath:   providedOpPath,
					RootOpID: "providedRootID",
				},
				ContainerID: "providedContainerID",
			}
			providedInboundScope := map[string]*model.Value{}
			providedContainerCallSpec := &model.ContainerCallSpec{}

			expectedEvent := model.Event{
				Timestamp: time.Now().UTC(),
				ContainerStarted: &model.ContainerStarted{
					ContainerID: providedContainerCall.ContainerID,
					OpRef:       providedOpPath,
					RootOpID:    providedContainerCall.RootOpID,
				},
			}

			fakePubSub := new(FakePubSub)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerRuntime: new(FakeContainerRuntime),
				pubSub:           fakePubSub,
				io:               fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				context.Background(),
				providedContainerCall,
				providedInboundScope,
				providedContainerCallSpec,
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
			providedContainerCall := &model.ContainerCall{
				BaseCall: model.BaseCall{},
			}
			fakeContainerRuntime := new(FakeContainerRuntime)

			fakePubSub := new(FakePubSub)

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
				providedContainerCall,
				map[string]*model.Value{},
				&model.ContainerCallSpec{},
			)

			/* assert */
			_,
				actualContainerCall,
				actualEventPublisher,
				_,
				_ := fakeContainerRuntime.RunContainerArgsForCall(0)
			Expect(actualContainerCall).To(Equal(providedContainerCall))
			Expect(actualEventPublisher).To(Equal(fakePubSub))
		})
		Context("containerRuntime.RunContainer errors", func() {
			It("should publish expected ContainerExited", func() {
				/* arrange */
				expectedErrorMessage := "expectedErrorMessage"
				fakePubSub := new(FakePubSub)

				fakeContainerRuntime := new(FakeContainerRuntime)
				fakeContainerRuntime.RunContainerReturns(nil, errors.New(expectedErrorMessage))

				fakeIIO := new(iio.Fake)
				fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

				objectUnderTest := _containerCaller{
					containerRuntime: fakeContainerRuntime,
					pubSub:           fakePubSub,
					io:               fakeIIO,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					&model.ContainerCall{
						BaseCall: model.BaseCall{},
					},
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
				)

				/* assert */
				actualEvent := fakePubSub.PublishArgsForCall(1)

				Expect(actualEvent.ContainerExited.Error.Message).To(Equal(expectedErrorMessage))
			})
		})
	})

	It("should call pubSub.Publish w/ expected ContainerExited", func() {
		/* arrange */
		providedOpPath := "providedOpPath"
		providedContainerCall := &model.ContainerCall{
			BaseCall: model.BaseCall{
				OpPath:   providedOpPath,
				RootOpID: "providedRootID",
			},
			ContainerID: "providedContainerID",
		}
		providedInboundScope := map[string]*model.Value{}
		providedContainerCallSpec := &model.ContainerCallSpec{}

		expectedEvent := model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExited{
				ContainerID: providedContainerCall.ContainerID,
				Error: &model.CallEndedError{
					Message: "io: read/write on closed pipe",
				},
				OpRef:    providedOpPath,
				Outputs:  map[string]*model.Value{},
				RootOpID: providedContainerCall.RootOpID,
			},
		}

		fakePubSub := new(FakePubSub)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(FakeContainerRuntime),
			pubSub:           fakePubSub,
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			context.Background(),
			providedContainerCall,
			providedInboundScope,
			providedContainerCallSpec,
		)

		/* assert */
		actualEvent := fakePubSub.PublishArgsForCall(1)

		// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
		Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
		// set temporal fields to expected vals since they're already asserted
		actualEvent.Timestamp = expectedEvent.Timestamp

		Expect(actualEvent).To(Equal(expectedEvent))
	})
	It("should call pubSub.Publish w/ expected ContainerExited", func() {
		/* arrange */
		providedOpPath := "providedOpPath"
		providedContainerCall := &model.ContainerCall{
			BaseCall: model.BaseCall{
				OpPath:   providedOpPath,
				RootOpID: "providedRootID",
			},
			ContainerID: "providedContainerID",
		}
		providedInboundScope := map[string]*model.Value{}
		providedContainerCallSpec := &model.ContainerCallSpec{}

		expectedEvent := model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExited{
				ContainerID: providedContainerCall.ContainerID,
				Error: &model.CallEndedError{
					Message: "io: read/write on closed pipe",
				},
				OpRef:    providedOpPath,
				Outputs:  map[string]*model.Value{},
				RootOpID: providedContainerCall.RootOpID,
			},
		}

		fakePubSub := new(FakePubSub)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(FakeContainerRuntime),
			pubSub:           fakePubSub,
			io:               fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			context.Background(),
			providedContainerCall,
			providedInboundScope,
			providedContainerCallSpec,
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
