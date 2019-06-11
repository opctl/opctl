package core

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

var _ = Context("parallelCaller", func() {
	Context("newParallelCaller", func() {
		It("should return parallelCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelCaller(
				new(fakeCaller),
				new(pubsub.Fake),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every parallelCall w/ expected args", func() {
			/* arrange */
			providedCallID := "dummyCallID"
			providedInboundScope := map[string]*model.Value{}
			providedRootOpID := "dummyRootOpID"
			providedOpHandle := new(data.FakeHandle)
			providedSCGParallelCalls := []*model.SCG{
				{
					Container: &model.SCGContainerCall{},
				},
				{
					Op: &model.SCGOpCall{},
				},
				{
					Parallel: []*model.SCG{},
				},
				{
					Serial: []*model.SCG{},
				},
			}

			fakeCaller := new(fakeCaller)
			eventChannel := make(chan model.Event, 100)
			callerCallIndex := 0
			fakeCaller.CallStub = func(context.Context, string, map[string]*model.Value, *model.SCG, model.DataHandle, *string, string) error {
				eventChannel <- model.Event{
					CallEnded: &model.CallEndedEvent{
						CallID: fmt.Sprintf("%v", callerCallIndex),
					},
				}

				callerCallIndex++
				return nil
			}

			fakePubSub := new(pubsub.Fake)
			fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
				return eventChannel, make(chan error)
			}

			fakeUniqueStringFactory := new(uniquestring.Fake)
			uniqueStringCallIndex := 0
			expectedChildCallIds := []string{}
			fakeUniqueStringFactory.ConstructStub = func() (string, error) {
				defer func() {
					uniqueStringCallIndex++
				}()
				childCallID := fmt.Sprintf("%v", uniqueStringCallIndex)
				expectedChildCallIds = append(expectedChildCallIds, fmt.Sprintf("%v", uniqueStringCallIndex))
				return childCallID, nil
			}

			objectUnderTest := _parallelCaller{
				caller:              fakeCaller,
				pubSub:              fakePubSub,
				uniqueStringFactory: fakeUniqueStringFactory,
			}

			/* act */
			objectUnderTest.Call(
				context.Background(),
				providedCallID,
				providedInboundScope,
				providedRootOpID,
				providedOpHandle,
				providedSCGParallelCalls,
			)

			/* assert */
			for callIndex := range providedSCGParallelCalls {
				_,
					actualNodeID,
					actualChildOutboundScope,
					actualSCG,
					actualOpHandle,
					actualParentCallID,
					actualRootOpID := fakeCaller.CallArgsForCall(callIndex)

				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualParentCallID).To(Equal(&providedCallID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))

				// handle unordered asserts because call order can't be relied on within go statement
				Expect(expectedChildCallIds).To(ContainElement(actualNodeID))
				Expect(providedSCGParallelCalls).To(ContainElement(actualSCG))
			}
		})
		Context("CallEnded event received w/ Error", func() {
			It("should fail fast on childCall error & return expected error", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedInboundScope := map[string]*model.Value{}
				providedRootOpID := "dummyRootOpID"
				providedOpHandle := new(data.FakeHandle)
				providedSCGParallelCalls := []*model.SCG{
					{
						Container: &model.SCGContainerCall{},
					},
					{
						Op: &model.SCGOpCall{},
					},
					{
						Parallel: []*model.SCG{},
					},
					{
						Serial: []*model.SCG{},
					},
				}

				errorMessage := "errorMessage"
				childErrorMessages := []string{}
				for range providedSCGParallelCalls {
					childErrorMessages = append(childErrorMessages, errorMessage)
				}

				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan model.Event, 100)
				fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
					for index := range providedSCGParallelCalls {
						eventChannel <- model.Event{
							CallEnded: &model.CallEndedEvent{
								CallID: fmt.Sprintf("%v", index),
								Error: &model.CallEndedEventError{
									Message: errorMessage,
								},
							},
						}
					}

					return eventChannel, make(chan error)
				}

				expectedErrorMessage := fmt.Sprintf(
					"-\nError(s) during parallel call. Error(s) were:\n%v\n-",
					strings.Join(childErrorMessages, "\n"),
				)

				fakeUniqueStringFactory := new(uniquestring.Fake)
				uniqueStringCallIndex := 0
				fakeUniqueStringFactory.ConstructStub = func() (string, error) {
					defer func() {
						uniqueStringCallIndex++
					}()
					return fmt.Sprintf("%v", uniqueStringCallIndex), nil
				}

				objectUnderTest := _parallelCaller{
					caller:              new(fakeCaller),
					pubSub:              fakePubSub,
					uniqueStringFactory: fakeUniqueStringFactory,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedInboundScope,
					providedRootOpID,
					providedOpHandle,
					providedSCGParallelCalls,
				)

				/* assert */
				actualEvent := fakePubSub.PublishArgsForCall(0)

				Expect(actualEvent.ParallelCallEnded.Error.Message).To(Equal(expectedErrorMessage))
			})
		})
		Context("caller doesn't error", func() {
			It("shouldn't exit until all childCalls complete & not error", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedInboundScope := map[string]*model.Value{}
				providedRootOpID := "dummyRootOpID"
				providedOpHandle := new(data.FakeHandle)
				providedSCGParallelCalls := []*model.SCG{
					{
						Container: &model.SCGContainerCall{},
					},
					{
						Op: &model.SCGOpCall{},
					},
					{
						Parallel: []*model.SCG{},
					},
					{
						Serial: []*model.SCG{},
					},
				}

				fakeCaller := new(fakeCaller)
				eventChannel := make(chan model.Event, 100)
				callerCallIndex := 0
				fakeCaller.CallStub = func(context.Context, string, map[string]*model.Value, *model.SCG, model.DataHandle, *string, string) error {
					eventChannel <- model.Event{
						CallEnded: &model.CallEndedEvent{
							CallID: fmt.Sprintf("%v", callerCallIndex),
						},
					}

					callerCallIndex++
					return nil
				}

				fakePubSub := new(pubsub.Fake)
				fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
					return eventChannel, make(chan error)
				}

				fakeUniqueStringFactory := new(uniquestring.Fake)
				uniqueStringCallIndex := 0
				expectedChildCallIds := []string{}
				fakeUniqueStringFactory.ConstructStub = func() (string, error) {
					defer func() {
						uniqueStringCallIndex++
					}()
					childCallID := fmt.Sprintf("%v", uniqueStringCallIndex)
					expectedChildCallIds = append(expectedChildCallIds, fmt.Sprintf("%v", uniqueStringCallIndex))
					return childCallID, nil
				}

				objectUnderTest := _parallelCaller{
					caller:              fakeCaller,
					pubSub:              fakePubSub,
					uniqueStringFactory: fakeUniqueStringFactory,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedInboundScope,
					providedRootOpID,
					providedOpHandle,
					providedSCGParallelCalls,
				)

				/* assert */
				for callIndex := range providedSCGParallelCalls {
					_,
						actualNodeID,
						actualChildOutboundScope,
						actualSCG,
						actualOpHandle,
						actualParentCallID,
						actualRootOpID := fakeCaller.CallArgsForCall(callIndex)

					Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
					Expect(actualOpHandle).To(Equal(providedOpHandle))
					Expect(actualParentCallID).To(Equal(&providedCallID))
					Expect(actualRootOpID).To(Equal(providedRootOpID))

					// handle unordered asserts because call order can't be relied on within go statement
					Expect(expectedChildCallIds).To(ContainElement(actualNodeID))
					Expect(providedSCGParallelCalls).To(ContainElement(actualSCG))
				}
			})
		})
	})
})
