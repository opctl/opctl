package core

import (
	"errors"
	"github.com/golang-interfaces/iio"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/containercall"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/util/containerprovider"
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
				new(containerprovider.Fake),
				new(containercall.Fake),
				new(pubsub.Fake),
				new(fakeDCGNodeRepo),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call dcgNodeRepo.Add w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Value{}
			providedContainerId := "dummyContainerId"
			providedSCGContainerCall := &model.SCGContainerCall{}
			providedPkgHandle := new(pkg.FakeHandle)
			providedRootOpId := "dummyRootOpId"

			fakePubSub := new(pubsub.Fake)

			fakeContainerCall := new(containercall.Fake)
			// error to trigger immediate return
			fakeContainerCall.InterpretReturns(nil, errors.New("dummyError"))

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:        providedContainerId,
				PkgRef:    providedPkgHandle.Ref(),
				RootOpId:  providedRootOpId,
				Container: &dcgContainerDescriptor{},
			}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerProvider: new(containerprovider.Fake),
				containerCall:     fakeContainerCall,
				pubSub:            fakePubSub,
				dcgNodeRepo:       fakeDCGNodeRepo,
				io:                fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerId,
				providedSCGContainerCall,
				providedPkgHandle,
				providedRootOpId,
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
		It("should call pubSub.Publish w/ expected ContainerStartedEvent", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Value{}
			providedContainerId := "dummyContainerId"
			providedSCGContainerCall := &model.SCGContainerCall{}
			providedPkgHandle := new(pkg.FakeHandle)
			providedRootOpId := "dummyRootOpId"

			expectedEvent := model.Event{
				Timestamp: time.Now().UTC(),
				ContainerStarted: &model.ContainerStartedEvent{
					ContainerId: providedContainerId,
					PkgRef:      providedPkgHandle.Ref(),
					RootOpId:    providedRootOpId,
				},
			}

			fakeContainerCall := new(containercall.Fake)
			fakeContainerCall.InterpretReturns(&model.DCGContainerCall{Image: &model.DCGContainerCallImage{}}, nil)

			fakePubSub := new(pubsub.Fake)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerProvider: new(containerprovider.Fake),
				containerCall:     fakeContainerCall,
				pubSub:            fakePubSub,
				dcgNodeRepo:       new(fakeDCGNodeRepo),
				io:                fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerId,
				providedSCGContainerCall,
				providedPkgHandle,
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
			expectedDCGContainerCall := &model.DCGContainerCall{}

			fakeContainerCall := new(containercall.Fake)
			fakeContainerCall.InterpretReturns(expectedDCGContainerCall, nil)

			fakeContainerProvider := new(containerprovider.Fake)

			fakePubSub := new(pubsub.Fake)

			fakeIIO := new(iio.Fake)
			fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

			objectUnderTest := _containerCaller{
				containerProvider: fakeContainerProvider,
				containerCall:     fakeContainerCall,
				pubSub:            fakePubSub,
				dcgNodeRepo:       new(fakeDCGNodeRepo),
				io:                fakeIIO,
			}

			/* act */
			objectUnderTest.Call(
				map[string]*model.Value{},
				"dummyContainerId",
				&model.SCGContainerCall{},
				new(pkg.FakeHandle),
				"dummyRootOpId",
			)

			/* assert */
			_, actualDCGContainerCall, actualEventPublisher, _, _ := fakeContainerProvider.RunContainerArgsForCall(0)
			Expect(actualDCGContainerCall).To(Equal(expectedDCGContainerCall))
			Expect(actualEventPublisher).To(Equal(fakePubSub))
		})
		Context("containerProvider.RunContainer errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")

				fakeContainerProvider := new(containerprovider.Fake)
				fakeContainerProvider.RunContainerReturns(nil, expectedError)

				fakeIIO := new(iio.Fake)
				fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

				objectUnderTest := _containerCaller{
					containerProvider: fakeContainerProvider,
					containerCall:     new(containercall.Fake),
					pubSub:            new(pubsub.Fake),
					dcgNodeRepo:       new(fakeDCGNodeRepo),
					io:                fakeIIO,
				}

				/* act */
				actualError := objectUnderTest.Call(
					map[string]*model.Value{},
					"dummyContainerId",
					&model.SCGContainerCall{},
					new(pkg.FakeHandle),
					"dummyRootOpId",
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
	It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
		/* arrange */
		providedContainerId := "dummyContainerId"

		fakeDCGNodeRepo := new(fakeDCGNodeRepo)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerProvider: new(containerprovider.Fake),
			containerCall:     new(containercall.Fake),
			pubSub:            new(pubsub.Fake),
			dcgNodeRepo:       fakeDCGNodeRepo,
			io:                fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			map[string]*model.Value{},
			providedContainerId,
			&model.SCGContainerCall{},
			new(pkg.FakeHandle),
			"dummyRootOpId",
		)

		/* assert */
		Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedContainerId))
	})

	It("should call pubSub.Publish w/ expected ContainerExitedEvent", func() {
		/* arrange */
		providedInboundScope := map[string]*model.Value{}
		providedContainerId := "dummyContainerId"
		providedSCGContainerCall := &model.SCGContainerCall{}
		providedPkgHandle := new(pkg.FakeHandle)
		providedRootOpId := "dummyRootOpId"

		expectedEvent := model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerId: providedContainerId,
				PkgRef:      providedPkgHandle.Ref(),
				RootOpId:    providedRootOpId,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerProvider: new(containerprovider.Fake),
			containerCall:     new(containercall.Fake),
			pubSub:            fakePubSub,
			dcgNodeRepo:       new(fakeDCGNodeRepo),
			io:                fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			providedInboundScope,
			providedContainerId,
			providedSCGContainerCall,
			providedPkgHandle,
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
	It("should call pubSub.Publish w/ expected ContainerExitedEvent", func() {
		/* arrange */
		providedInboundScope := map[string]*model.Value{}
		providedContainerId := "dummyContainerId"
		providedSCGContainerCall := &model.SCGContainerCall{}
		providedPkgHandle := new(pkg.FakeHandle)
		providedRootOpId := "dummyRootOpId"

		expectedEvent := model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerId: providedContainerId,
				PkgRef:      providedPkgHandle.Ref(),
				RootOpId:    providedRootOpId,
			},
		}

		fakePubSub := new(pubsub.Fake)

		fakeIIO := new(iio.Fake)
		fakeIIO.PipeReturns(closedPipeReader, closedPipeWriter)

		objectUnderTest := _containerCaller{
			containerProvider: new(containerprovider.Fake),
			containerCall:     new(containercall.Fake),
			pubSub:            fakePubSub,
			dcgNodeRepo:       new(fakeDCGNodeRepo),
			io:                fakeIIO,
		}

		/* act */
		objectUnderTest.Call(
			providedInboundScope,
			providedContainerId,
			providedSCGContainerCall,
			providedPkgHandle,
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
