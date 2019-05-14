package core

import (
	"context"
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
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
			providedCallID := "providedCallID"
			providedInboundScope := map[string]*model.Value{}
			providedRootOpID := "providedRootOpID"
			providedOpHandle := new(data.FakeHandle)
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

			fakePubSub := new(pubsub.Fake)
			subscribeCallIndex := 0
			fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
				defer func() {
					subscribeCallIndex++
				}()
				eventChannel := make(chan model.Event, 100)
				eventChannel <- model.Event{OpEnded: &model.OpEndedEvent{OpID: fmt.Sprintf("%v", subscribeCallIndex)}}
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
				providedCallID,
				providedInboundScope,
				providedRootOpID,
				providedOpHandle,
				providedSCGSerialCalls,
			)

			/* assert */
			for expectedSCGIndex, expectedSCG := range providedSCGSerialCalls {
				actualNodeId,
					actualChildOutboundScope,
					actualSCG,
					actualOpHandle,
					actualParentCallID,
					actualRootOpID := fakeCaller.CallArgsForCall(expectedSCGIndex)
				Expect(actualNodeId).To(Equal(fmt.Sprintf("%v", expectedSCGIndex)))
				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualSCG).To(Equal(expectedSCG))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualParentCallID).To(Equal(&providedCallID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			}
		})
		Context("caller errors", func() {
			It("should return the expected error", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedInboundScope := map[string]*model.Value{}
				providedRootOpID := "dummyRootOpID"
				providedOpHandle := new(data.FakeHandle)
				providedSCGSerialCalls := []*model.SCG{
					{
						Container: &model.SCGContainerCall{},
					},
				}

				callErr := errors.New("dummyErr")

				fakeCaller := new(fakeCaller)
				fakeCaller.CallReturns(callErr)

				objectUnderTest := _serialCaller{
					caller:              fakeCaller,
					pubSub:              new(pubsub.Fake),
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				actualErr := objectUnderTest.Call(
					providedCallID,
					providedInboundScope,
					providedRootOpID,
					providedOpHandle,
					providedSCGSerialCalls,
				)

				/* assert */
				Expect(actualErr).To(Equal(callErr))
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
					providedOpHandle := new(data.FakeHandle)
					providedSCGSerialCalls := []*model.SCG{
						{
							Container: &model.SCGContainerCall{},
						},
						{
							Container: &model.SCGContainerCall{},
						},
					}

					fakePubSub := new(pubsub.Fake)
					subscribeCallIndex := 0
					fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
						defer func() {
							subscribeCallIndex++
						}()
						eventChannel := make(chan model.Event, 100)
						eventChannel <- model.Event{OpEnded: &model.OpEndedEvent{OpID: fmt.Sprintf("%v", subscribeCallIndex)}}
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
						providedCallID,
						providedInboundScope,
						providedRootOpID,
						providedOpHandle,
						providedSCGSerialCalls,
					)

					/* assert */
					_, actualInboundScopeToSecondChild, _, _, _, _ := fakeCaller.CallArgsForCall(1)
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
					providedOpHandle := new(data.FakeHandle)
					providedSCGSerialCalls := []*model.SCG{
						{
							Container: &model.SCGContainerCall{},
						},
						{
							Container: &model.SCGContainerCall{},
						},
					}

					fakePubSub := new(pubsub.Fake)
					subscribeCallIndex := 0
					fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
						defer func() {
							subscribeCallIndex++
						}()
						eventChannel := make(chan model.Event, 100)
						eventChannel <- model.Event{
							ContainerExited: &model.ContainerExitedEvent{
								RootOpID:    providedRootOpID,
								ContainerID: fmt.Sprintf("%v", subscribeCallIndex),
								Outputs:     firstChildOutputs,
							},
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
						providedCallID,
						providedInboundScope,
						providedRootOpID,
						providedOpHandle,
						providedSCGSerialCalls,
					)

					/* assert */
					_, actualInboundScopeToSecondChild, _, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualInboundScopeToSecondChild).To(Equal(expectedInboundScopeToSecondChild))
				})
			})
		})
	})
})
