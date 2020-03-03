package core

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uniquestringFakes "github.com/opctl/opctl/sdks/go/internal/uniquestring/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("serialCaller", func() {
	Context("newSerialCaller", func() {
		It("should return serialCaller", func() {
			/* arrange/act/assert */
			Expect(newSerialCaller(
				new(FakeCaller),
				new(FakePubSub),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every serialCall w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedCallID := "providedCallID"
			providedInboundScope := map[string]*model.Value{}
			providedRootOpID := "providedRootOpID"
			providedOpHandle := new(modelFakes.FakeDataHandle)
			providedSCGSerialCalls := []*model.SCG{
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

			fakePubSub := new(FakePubSub)
			eventChannel := make(chan model.Event, 100)
			fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
				for index := range providedSCGSerialCalls {
					eventChannel <- model.Event{
						CallEnded: &model.CallEndedEvent{
							CallID: fmt.Sprintf("%v", index),
						},
					}
				}
				return eventChannel, make(chan error)
			}

			fakeCaller := new(FakeCaller)

			fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
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
				providedInboundScope := map[string]*model.Value{}
				providedRootOpID := "dummyRootOpID"
				providedOpHandle := new(modelFakes.FakeDataHandle)
				providedSCGSerialCalls := []*model.SCG{
					{
						Container: &model.SCGContainerCall{},
					},
				}

				callID := "callID"

				expectedErrorMessage := "expectedErrorMessage"
				fakePubSub := new(FakePubSub)
				eventChannel := make(chan model.Event, 100)
				fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
					for range providedSCGSerialCalls {
						eventChannel <- model.Event{
							CallEnded: &model.CallEndedEvent{
								CallID: callID,
								Error: &model.CallEndedEventError{
									Message: expectedErrorMessage,
								},
							},
						}
					}

					return eventChannel, make(chan error)
				}

				fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
				fakeUniqueStringFactory.ConstructReturns(callID, nil)

				objectUnderTest := _serialCaller{
					caller:              new(FakeCaller),
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
					providedInboundScope := map[string]*model.Value{
						"dummyVar1Name": {String: &providedScopeName1String},
						"dummyVar2Name": {Dir: &providedScopeName2Dir},
					}
					expectedInboundScopeToSecondChild := providedInboundScope
					providedRootOpID := "dummyRootOpID"
					providedOpHandle := new(modelFakes.FakeDataHandle)
					providedSCGSerialCalls := []*model.SCG{
						{
							Container: &model.SCGContainerCall{},
						},
						{
							Container: &model.SCGContainerCall{},
						},
					}

					fakePubSub := new(FakePubSub)
					eventChannel := make(chan model.Event, 100)
					fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
						for index := range providedSCGSerialCalls {
							eventChannel <- model.Event{
								CallEnded: &model.CallEndedEvent{
									CallID: fmt.Sprintf("%v", index),
								},
							}
						}

						return eventChannel, make(chan error)
					}

					fakeCaller := new(FakeCaller)

					fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
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
					providedInboundScope := map[string]*model.Value{
						"dummyVar1Name": {String: &providedInboundVar1String},
						"dummyVar2Name": {Dir: &providedInboundVar2Dir},
						"dummyVar3Name": {File: &providedInboundVar3File},
					}

					firstChildOutput1String := "dummyFirstChildVar1Data"
					firstChildOutput2String := "dummyFirstChildVar2Data"
					firstChildOutputs := map[string]*model.Value{
						"dummyVar1Name": {String: &firstChildOutput1String},
						"dummyVar2Name": {Dir: &firstChildOutput2String},
					}

					expectedInboundScopeToSecondChild := map[string]*model.Value{
						"dummyVar1Name": firstChildOutputs["dummyVar1Name"],
						"dummyVar2Name": firstChildOutputs["dummyVar2Name"],
						"dummyVar3Name": providedInboundScope["dummyVar3Name"],
					}
					providedRootOpID := "dummyRootOpID"
					providedOpHandle := new(modelFakes.FakeDataHandle)
					providedSCGSerialCalls := []*model.SCG{
						{
							Container: &model.SCGContainerCall{},
						},
						{
							Container: &model.SCGContainerCall{},
						},
					}

					fakePubSub := new(FakePubSub)
					eventChannel := make(chan model.Event, 100)
					fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
						for index := range providedSCGSerialCalls {
							eventChannel <- model.Event{
								CallEnded: &model.CallEndedEvent{
									RootCallID: providedRootOpID,
									CallID:     fmt.Sprintf("%v", index),
									Outputs:    firstChildOutputs,
								},
							}
						}

						return eventChannel, make(chan error)
					}

					fakeCaller := new(FakeCaller)

					fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
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
