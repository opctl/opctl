package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("serialCaller", func() {
	Context("newSerialCaller", func() {
		It("should return serialCaller", func() {
			/* arrange/act/assert */
			Expect(newSerialCaller(
				new(fakeCaller),
				new(pubsub.Fake),
				new(uniquestring.Fake),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every serialCall w/ expected args", func() {
			/* arrange */
			providedCallId := "dummyCallId"
			providedInboundScope := map[string]*model.Value{}
			providedRootOpId := "dummyRootOpId"
			providedPkgRef := "dummyPkgRef"
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
			fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
				defer func() {
					subscribeCallIndex++
				}()
				eventChannel <- &model.Event{OpEnded: &model.OpEndedEvent{OpId: fmt.Sprintf("%v", subscribeCallIndex)}}
			}

			fakeCaller := new(fakeCaller)

			fakeUniqueStringFactory := new(uniquestring.Fake)
			uniqueStringCallIndex := 0
			fakeUniqueStringFactory.ConstructStub = func() (uniqueString string) {
				defer func() {
					uniqueStringCallIndex++
				}()
				return fmt.Sprintf("%v", uniqueStringCallIndex)
			}

			objectUnderTest := newSerialCaller(fakeCaller, fakePubSub, fakeUniqueStringFactory)

			/* act */
			objectUnderTest.Call(
				providedCallId,
				providedInboundScope,
				providedRootOpId,
				providedPkgRef,
				providedSCGSerialCalls,
			)

			/* assert */
			for expectedSCGIndex, expectedSCG := range providedSCGSerialCalls {
				actualNodeId,
					actualChildOutboundScope,
					actualSCG,
					actualPkgRef,
					actualRootOpId := fakeCaller.CallArgsForCall(expectedSCGIndex)
				Expect(actualNodeId).To(Equal(fmt.Sprintf("%v", expectedSCGIndex)))
				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualSCG).To(Equal(expectedSCG))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			}
		})
		Context("caller errors", func() {
			It("should return the expected error", func() {
				/* arrange */
				providedCallId := "dummyCallId"
				providedInboundScope := map[string]*model.Value{}
				providedRootOpId := "dummyRootOpId"
				providedPkgRef := "dummyPkgRef"
				providedSCGSerialCalls := []*model.SCG{
					{
						Container: &model.SCGContainerCall{},
					},
				}

				callErr := errors.New("dummyErr")

				fakeCaller := new(fakeCaller)
				fakeCaller.CallReturns(callErr)

				objectUnderTest := newSerialCaller(fakeCaller, new(pubsub.Fake), new(uniquestring.Fake))

				/* act */
				actualErr := objectUnderTest.Call(
					providedCallId,
					providedInboundScope,
					providedRootOpId,
					providedPkgRef,
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
					providedCallId := "dummyCallId"
					providedScopeName1String := "dummyParentVar1Data"
					providedScopeName2Dir := "dummyParentVar2Data"
					providedInboundScope := map[string]*model.Value{
						"dummyVar1Name": {String: &providedScopeName1String},
						"dummyVar2Name": {Dir: &providedScopeName2Dir},
					}
					expectedInboundScopeToSecondChild := providedInboundScope
					providedRootOpId := "dummyRootOpId"
					providedPkgRef := "dummyPkgRef"
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
					fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
						defer func() {
							subscribeCallIndex++
						}()
						eventChannel <- &model.Event{OpEnded: &model.OpEndedEvent{OpId: fmt.Sprintf("%v", subscribeCallIndex)}}
					}

					fakeCaller := new(fakeCaller)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					uniqueStringCallIndex := 0
					fakeUniqueStringFactory.ConstructStub = func() (uniqueString string) {
						defer func() {
							uniqueStringCallIndex++
						}()
						return fmt.Sprintf("%v", uniqueStringCallIndex)
					}

					objectUnderTest := newSerialCaller(fakeCaller, fakePubSub, fakeUniqueStringFactory)

					/* act */
					objectUnderTest.Call(
						providedCallId,
						providedInboundScope,
						providedRootOpId,
						providedPkgRef,
						providedSCGSerialCalls,
					)

					/* assert */
					_, actualInboundScopeToSecondChild, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualInboundScopeToSecondChild).To(Equal(expectedInboundScopeToSecondChild))
				})
			})
			Context("childOutboundScope not empty", func() {
				It("should call secondChild w/ firstChildOutputs overlaying inboundScope", func() {
					/* arrange */
					providedCallId := "dummyCallId"

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
					providedRootOpId := "dummyRootOpId"
					providedPkgRef := "dummyPkgRef"
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
					fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
						defer func() {
							subscribeCallIndex++
						}()
						eventChannel <- &model.Event{
							ContainerExited: &model.ContainerExitedEvent{
								RootOpId:    providedRootOpId,
								ContainerId: fmt.Sprintf("%v", subscribeCallIndex),
								Outputs:     firstChildOutputs,
							},
						}
					}

					fakeCaller := new(fakeCaller)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					uniqueStringCallIndex := 0
					fakeUniqueStringFactory.ConstructStub = func() (uniqueString string) {
						defer func() {
							uniqueStringCallIndex++
						}()
						return fmt.Sprintf("%v", uniqueStringCallIndex)
					}

					objectUnderTest := newSerialCaller(fakeCaller, fakePubSub, fakeUniqueStringFactory)

					/* act */
					objectUnderTest.Call(
						providedCallId,
						providedInboundScope,
						providedRootOpId,
						providedPkgRef,
						providedSCGSerialCalls,
					)

					/* assert */
					_, actualInboundScopeToSecondChild, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualInboundScopeToSecondChild).To(Equal(expectedInboundScopeToSecondChild))
				})
			})
		})
	})
})
