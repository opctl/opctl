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
			)).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every serialCall w/ expected args", func() {
			/* arrange */
			providedCallId := "dummyCallId"
			providedInboundScope := map[string]*model.Data{}
			providedRootOpId := "dummyRootOpId"
			providedPkgRef := "dummyPkgRef"
			providedScgSerialCalls := []*model.Scg{
				{
					Container: &model.ScgContainerCall{},
				},
				{
					Op: &model.ScgOpCall{},
				},
				{
					Parallel: []*model.Scg{},
				},
				{
					Serial: []*model.Scg{},
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
				providedScgSerialCalls,
			)

			/* assert */
			for expectedScgIndex, expectedScg := range providedScgSerialCalls {
				actualNodeId,
					actualChildOutboundScope,
					actualScg,
					actualPkgRef,
					actualRootOpId := fakeCaller.CallArgsForCall(expectedScgIndex)
				Expect(actualNodeId).To(Equal(fmt.Sprintf("%v", expectedScgIndex)))
				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualScg).To(Equal(expectedScg))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			}
		})
		Context("caller errors", func() {
			It("should return the expected error", func() {
				/* arrange */
				providedCallId := "dummyCallId"
				providedInboundScope := map[string]*model.Data{}
				providedRootOpId := "dummyRootOpId"
				providedPkgRef := "dummyPkgRef"
				providedScgSerialCalls := []*model.Scg{
					{
						Container: &model.ScgContainerCall{},
					},
				}

				expectedError := errors.New("Error encountered during serial call")
				fakeCaller := new(fakeCaller)
				fakeCaller.CallReturns(errors.New("dummyError"))

				objectUnderTest := newSerialCaller(fakeCaller, new(pubsub.Fake), new(uniquestring.Fake))

				/* act */
				actualErr := objectUnderTest.Call(
					providedCallId,
					providedInboundScope,
					providedRootOpId,
					providedPkgRef,
					providedScgSerialCalls,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedError))
			})
		})
		Context("caller doesn't error", func() {
			Context("childOutboundScope empty", func() {
				It("should call secondChild w/ inboundScope", func() {
					/* arrange */
					providedCallId := "dummyCallId"
					providedInboundScope := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyParentVar1Data"},
						"dummyVar2Name": {Dir: "dummyParentVar2Data"},
					}
					expectedInboundScopeToSecondChild := providedInboundScope
					providedRootOpId := "dummyRootOpId"
					providedPkgRef := "dummyPkgRef"
					providedScgSerialCalls := []*model.Scg{
						{
							Container: &model.ScgContainerCall{},
						},
						{
							Container: &model.ScgContainerCall{},
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
						providedScgSerialCalls,
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
					providedInboundScope := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyParentVar1Data"},
						"dummyVar2Name": {Dir: "dummyParentVar2Data"},
						"dummyVar3Name": {File: "dummyParentVar3Data"},
					}
					firstChildOutputs := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyFirstChildVar1Data"},
						"dummyVar2Name": {Dir: "dummyFirstChildVar2Data"},
					}
					expectedInboundScopeToSecondChild := map[string]*model.Data{
						"dummyVar1Name": firstChildOutputs["dummyVar1Name"],
						"dummyVar2Name": firstChildOutputs["dummyVar2Name"],
						"dummyVar3Name": providedInboundScope["dummyVar3Name"],
					}
					providedRootOpId := "dummyRootOpId"
					providedPkgRef := "dummyPkgRef"
					providedScgSerialCalls := []*model.Scg{
						{
							Container: &model.ScgContainerCall{},
						},
						{
							Container: &model.ScgContainerCall{},
						},
					}

					fakePubSub := new(pubsub.Fake)
					subscribeCallIndex := 0
					fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
						defer func() {
							subscribeCallIndex++
						}()
						for outputName, outputValue := range firstChildOutputs {
							eventChannel <- &model.Event{
								OutputInitialized: &model.OutputInitializedEvent{
									Name:     outputName,
									Value:    outputValue,
									RootOpId: providedRootOpId,
									CallId:   fmt.Sprintf("%v", subscribeCallIndex),
								},
							}
						}
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
						providedScgSerialCalls,
					)

					/* assert */
					_, actualInboundScopeToSecondChild, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualInboundScopeToSecondChild).To(Equal(expectedInboundScopeToSecondChild))
				})
			})
		})
	})
})
