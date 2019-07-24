package core

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/types"
	"github.com/opctl/opctl/sdks/go/util/pubsub"
	"github.com/opctl/opctl/sdks/go/util/uniquestring"
)

var _ = Context("serialCaller", func() {
	Context("newSerialCaller", func() {
		It("should return serialCaller", func() {
			/* arrange/act/assert */
			Expect(newSerialCaller(
				new(fakeCaller),
				new(pubsub.Fake),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every serialCall w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedCallID := "providedCallID"
			providedInboundScope := map[string]*types.Value{}
			providedRootOpID := "providedRootOpID"
			providedOpHandle := new(data.FakeHandle)
			providedSCGSerialCalls := []*types.SCG{
				{
					Container: &types.SCGContainerCall{},
				},
				{
					Op: &types.SCGOpCall{},
				},
				{
					Parallel: []*types.SCG{},
				},
				{
					Serial: []*types.SCG{},
				},
			}

			fakePubSub := new(pubsub.Fake)
			eventChannel := make(chan types.Event, 100)
			fakePubSub.SubscribeStub = func(ctx context.Context, filter types.EventFilter) (<-chan types.Event, <-chan error) {
				for index := range providedSCGSerialCalls {
					eventChannel <- types.Event{
						CallEnded: &types.CallEndedEvent{
							CallID: fmt.Sprintf("%v", index),
						},
					}
				}
				return eventChannel, make(chan error)
			}

			fakeCaller := new(fakeCaller)

			fakeUniqueStringFactory := new(uniquestring.Fake)
			uniqueStringCallIndex := 0
			fakeUniqueStringFactory.ConstructStub = func() (string, error) {
				defer func() {
					uniqueStringCallIndex++
				}()
				return fmt.Sprintf("%v", uniqueStringCallIndex), nil
			}

			objectUnderTest := _serialCaller{
				caller:              fakeCaller,
				pubSub:              fakePubSub,
				uniqueStringFactory: fakeUniqueStringFactory,
			}

			/* act */
			objectUnderTest.Call(
				providedCtx,
				providedCallID,
				providedInboundScope,
				providedRootOpID,
				providedOpHandle,
				providedSCGSerialCalls,
			)

			/* assert */
			for expectedSCGIndex, expectedSCG := range providedSCGSerialCalls {
				actualCtx,
					actualNodeID,
					actualChildOutboundScope,
					actualSCG,
					actualOpHandle,
					actualParentCallID,
					actualRootOpID := fakeCaller.CallArgsForCall(expectedSCGIndex)

				Expect(actualCtx).To(Equal(actualCtx))
				Expect(actualNodeID).To(Equal(fmt.Sprintf("%v", expectedSCGIndex)))
				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualSCG).To(Equal(expectedSCG))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualParentCallID).To(Equal(&providedCallID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			}
		})
		Context("caller errors", func() {
			It("should publish expected CallEndedEvent", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedInboundScope := map[string]*types.Value{}
				providedRootOpID := "dummyRootOpID"
				providedOpHandle := new(data.FakeHandle)
				providedSCGSerialCalls := []*types.SCG{
					{
						Container: &types.SCGContainerCall{},
					},
				}

				callID := "callID"

				expectedErrorMessage := "expectedErrorMessage"
				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan types.Event, 100)
				fakePubSub.SubscribeStub = func(ctx context.Context, filter types.EventFilter) (<-chan types.Event, <-chan error) {
					for range providedSCGSerialCalls {
						eventChannel <- types.Event{
							CallEnded: &types.CallEndedEvent{
								CallID: callID,
								Error: &types.CallEndedEventError{
									Message: expectedErrorMessage,
								},
							},
						}
					}

					return eventChannel, make(chan error)
				}

				fakeUniqueStringFactory := new(uniquestring.Fake)
				fakeUniqueStringFactory.ConstructReturns(callID, nil)

				objectUnderTest := _serialCaller{
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
					providedSCGSerialCalls,
				)

				/* assert */
				actualEvent := fakePubSub.PublishArgsForCall(0)

				Expect(actualEvent.SerialCallEnded.Error.Message).To(Equal(expectedErrorMessage))
			})
		})
		Context("caller doesn't error", func() {
			Context("childOutboundScope empty", func() {
				It("should call secondChild w/ inboundScope", func() {
					/* arrange */
					providedCallID := "dummyCallID"
					providedScopeName1String := "dummyParentVar1Data"
					providedScopeName2Dir := "dummyParentVar2Data"
					providedInboundScope := map[string]*types.Value{
						"dummyVar1Name": {String: &providedScopeName1String},
						"dummyVar2Name": {Dir: &providedScopeName2Dir},
					}
					expectedInboundScopeToSecondChild := providedInboundScope
					providedRootOpID := "dummyRootOpID"
					providedOpHandle := new(data.FakeHandle)
					providedSCGSerialCalls := []*types.SCG{
						{
							Container: &types.SCGContainerCall{},
						},
						{
							Container: &types.SCGContainerCall{},
						},
					}

					fakePubSub := new(pubsub.Fake)
					eventChannel := make(chan types.Event, 100)
					fakePubSub.SubscribeStub = func(ctx context.Context, filter types.EventFilter) (<-chan types.Event, <-chan error) {
						for index := range providedSCGSerialCalls {
							eventChannel <- types.Event{
								CallEnded: &types.CallEndedEvent{
									CallID: fmt.Sprintf("%v", index),
								},
							}
						}

						return eventChannel, make(chan error)
					}

					fakeCaller := new(fakeCaller)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					uniqueStringCallIndex := 0
					fakeUniqueStringFactory.ConstructStub = func() (string, error) {
						defer func() {
							uniqueStringCallIndex++
						}()
						return fmt.Sprintf("%v", uniqueStringCallIndex), nil
					}

					objectUnderTest := _serialCaller{
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
						providedSCGSerialCalls,
					)

					/* assert */
					_, _, actualInboundScopeToSecondChild, _, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualInboundScopeToSecondChild).To(Equal(expectedInboundScopeToSecondChild))
				})
			})
			Context("childOutboundScope not empty", func() {
				It("should call secondChild w/ firstChildOutputs overlaying inboundScope", func() {
					/* arrange */
					providedCallID := "dummyCallID"

					providedInboundVar1String := "dummyParentVar1Data"
					providedInboundVar2Dir := "dummyParentVar2Data"
					providedInboundVar3File := "dummyParentVar3Data"
					providedInboundScope := map[string]*types.Value{
						"dummyVar1Name": {String: &providedInboundVar1String},
						"dummyVar2Name": {Dir: &providedInboundVar2Dir},
						"dummyVar3Name": {File: &providedInboundVar3File},
					}

					firstChildOutput1String := "dummyFirstChildVar1Data"
					firstChildOutput2String := "dummyFirstChildVar2Data"
					firstChildOutputs := map[string]*types.Value{
						"dummyVar1Name": {String: &firstChildOutput1String},
						"dummyVar2Name": {Dir: &firstChildOutput2String},
					}

					expectedInboundScopeToSecondChild := map[string]*types.Value{
						"dummyVar1Name": firstChildOutputs["dummyVar1Name"],
						"dummyVar2Name": firstChildOutputs["dummyVar2Name"],
						"dummyVar3Name": providedInboundScope["dummyVar3Name"],
					}
					providedRootOpID := "dummyRootOpID"
					providedOpHandle := new(data.FakeHandle)
					providedSCGSerialCalls := []*types.SCG{
						{
							Container: &types.SCGContainerCall{},
						},
						{
							Container: &types.SCGContainerCall{},
						},
					}

					fakePubSub := new(pubsub.Fake)
					eventChannel := make(chan types.Event, 100)
					fakePubSub.SubscribeStub = func(ctx context.Context, filter types.EventFilter) (<-chan types.Event, <-chan error) {
						for index := range providedSCGSerialCalls {
							eventChannel <- types.Event{
								CallEnded: &types.CallEndedEvent{
									RootCallID: providedRootOpID,
									CallID:     fmt.Sprintf("%v", index),
									Outputs:    firstChildOutputs,
								},
							}
						}

						return eventChannel, make(chan error)
					}

					fakeCaller := new(fakeCaller)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					uniqueStringCallIndex := 0
					fakeUniqueStringFactory.ConstructStub = func() (string, error) {
						defer func() {
							uniqueStringCallIndex++
						}()
						return fmt.Sprintf("%v", uniqueStringCallIndex), nil
					}

					objectUnderTest := _serialCaller{
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
						providedSCGSerialCalls,
					)

					/* assert */
					_, _, actualInboundScopeToSecondChild, _, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualInboundScopeToSecondChild).To(Equal(expectedInboundScopeToSecondChild))
				})
			})
		})
	})
})
