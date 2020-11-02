package core

import (
	"context"
	"errors"
	"io"

	"github.com/golang-interfaces/iio"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
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
				new(FakeStateStore),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call containerRuntime.RunContainer w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedContainerCall := &model.ContainerCall{
				BaseCall: model.BaseCall{},
				Image:    &model.ContainerCallImage{},
			}
			fakeContainerRuntime := new(FakeContainerRuntime)

			fakePubSub := new(FakePubSub)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerRuntime: fakeContainerRuntime,
				pubSub:           fakePubSub,
				stateStore:       new(FakeStateStore),
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
					stateStore:       new(FakeStateStore),
					io:               fakeIIO,
				}

				/* act */
				actualOutputs, actualErr := objectUnderTest.Call(
					context.Background(),
					&model.ContainerCall{
						BaseCall: model.BaseCall{},
						Image:    &model.ContainerCallImage{},
					},
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
				)

				/* assert */
				Expect(actualOutputs).To(Equal(map[string]*model.Value{}))
				Expect(actualErr).To(Equal(errors.New(expectedErrorMessage)))
			})
		})
	})

	It("should return expected results", func() {
		/* arrange */
		providedOpPath := "providedOpPath"
		providedContainerCall := &model.ContainerCall{
			BaseCall: model.BaseCall{
				OpPath:     providedOpPath,
				RootCallID: "providedRootID",
			},
			ContainerID: "providedContainerID",
			Image:       &model.ContainerCallImage{},
		}
		providedInboundScope := map[string]*model.Value{}
		providedContainerCallSpec := &model.ContainerCallSpec{}

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerRuntime: new(FakeContainerRuntime),
			pubSub:           new(FakePubSub),
			stateStore:       new(FakeStateStore),
			io:               fakeIIO,
		}

		/* act */
		actualOutputs, actualErr := objectUnderTest.Call(
			context.Background(),
			providedContainerCall,
			providedInboundScope,
			providedContainerCallSpec,
		)

		/* assert */
		Expect(actualOutputs).To(Equal(map[string]*model.Value{}))
		Expect(actualErr).To(Equal(errors.New("io: read/write on closed pipe")))
	})
})
