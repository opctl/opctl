package core

import (
	"context"
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uniquestringFakes "github.com/opctl/opctl/sdks/go/internal/uniquestring/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("parallelCaller", func() {
	Context("newParallelCaller", func() {
		It("should return parallelCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelCaller(
				new(FakeCaller),
				new(FakePubSub),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every parallelCall w/ expected args", func() {
			/* arrange */
			providedCallID := "dummyCallID"
			providedInboundScope := map[string]*model.Value{}
			providedRootOpID := "dummyRootOpID"
			providedOpHandle := new(modelFakes.FakeDataHandle)
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

			mtx := sync.Mutex{}

			fakeCaller := new(FakeCaller)
			eventChannel := make(chan model.Event, 100)
			callerCallIndex := 0
			fakeCaller.CallStub = func(context.Context, string, map[string]*model.Value, *model.SCG, model.DataHandle, *string, string) {
				mtx.Lock()
				eventChannel <- model.Event{
					CallEnded: &model.CallEndedEvent{
						CallID: fmt.Sprintf("%v", callerCallIndex),
					},
				}

				callerCallIndex++

				mtx.Unlock()
			}

			fakePubSub := new(FakePubSub)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
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
			It("should publish expected CallEndedEvent", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedInboundScope := map[string]*model.Value{}
				providedRootOpID := "dummyRootOpID"
				providedOpHandle := new(modelFakes.FakeDataHandle)
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

				mtx := sync.Mutex{}

				fakeCaller := new(FakeCaller)
				eventChannel := make(chan model.Event, 100)
				callerCallIndex := 0
				fakeCaller.CallStub = func(context.Context, string, map[string]*model.Value, *model.SCG, model.DataHandle, *string, string) {
					mtx.Lock()

					eventChannel <- model.Event{
						CallEnded: &model.CallEndedEvent{
							CallID: fmt.Sprintf("%v", callerCallIndex),
							Error: &model.CallEndedEventError{
								Message: errorMessage,
							},
						},
					}

					callerCallIndex++

					mtx.Unlock()
				}

				fakePubSub := new(FakePubSub)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
				uniqueStringCallIndex := 0
				expectedChildCallIds := []string{}
				fakeUniqueStringFactory.ConstructStub = func() (string, error) {
					defer func() {
						uniqueStringCallIndex++
					}()
					childCallID := fmt.Sprintf("%v", uniqueStringCallIndex)
					expectedChildCallIds = append(expectedChildCallIds, childCallID)
					return childCallID, nil
				}

				var formattedChildErrorMessages string
				for _, childErrorMessage := range childErrorMessages {
					formattedChildErrorMessages = fmt.Sprintf("\t-%v\n", childErrorMessage)
				}
				expectedErrorMessage := fmt.Sprintf(
					"-\nError(s) during parallel call. Error(s) were:\n%v\n-",
					formattedChildErrorMessages,
				)

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
				providedOpHandle := new(modelFakes.FakeDataHandle)
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

				mtx := sync.Mutex{}

				fakeCaller := new(FakeCaller)
				eventChannel := make(chan model.Event, 100)
				callerCallIndex := 0
				fakeCaller.CallStub = func(context.Context, string, map[string]*model.Value, *model.SCG, model.DataHandle, *string, string) {
					mtx.Lock()

					eventChannel <- model.Event{
						CallEnded: &model.CallEndedEvent{
							CallID: fmt.Sprintf("%v", callerCallIndex),
						},
					}

					callerCallIndex++
					mtx.Unlock()
				}

				fakePubSub := new(FakePubSub)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
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
