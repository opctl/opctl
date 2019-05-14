package core

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/outputs"
	dotyml "github.com/opctl/sdk-golang/opspec/opfile"
	"github.com/opctl/sdk-golang/util/pubsub"
)

var _ = Context("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(fakeCallStore),
				new(pubsub.Fake),
				new(fakeCaller),
				"",
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call pubSub.Publish w/ expected args", func() {
			/* arrange */
			providedOpHandleRef := "dummyOpRef"
			fakeOpHandle := new(data.FakeHandle)
			fakeOpHandle.RefReturns(providedOpHandleRef)

			providedDCGOpCall := &model.DCGOpCall{
				DCGBaseCall: model.DCGBaseCall{
					OpHandle: fakeOpHandle,
					RootOpID: "providedRootID",
				},
				OpID: "providedOpId",
			}

			providedSCGOpCall := &model.SCGOpCall{}

			expectedEvent := model.Event{
				Timestamp: time.Now().UTC(),
				OpStarted: &model.OpStartedEvent{
					OpID:     providedDCGOpCall.OpID,
					OpRef:    providedOpHandleRef,
					RootOpID: providedDCGOpCall.RootOpID,
				},
			}

			fakePubSub := new(pubsub.Fake)
			eventChannel := make(chan model.Event)
			// close eventChannel to trigger immediate return
			close(eventChannel)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeCaller := new(fakeCaller)
			// err to trigger immediate return
			fakeCaller.CallReturns(errors.New("dummyError"))

			objectUnderTest := _opCaller{
				caller:    fakeCaller,
				callStore: new(fakeCallStore),
				pubSub:    fakePubSub,
			}

			/* act */
			objectUnderTest.Call(
				providedDCGOpCall,
				map[string]*model.Value{},
				nil,
				providedSCGOpCall,
			)

			/* assert */
			actualEvent := fakePubSub.PublishArgsForCall(0)

			// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
			Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
			// set temporal fields to expected vals since they're already asserted
			actualEvent.Timestamp = expectedEvent.Timestamp

			Expect(actualEvent).To(Equal(expectedEvent))
		})
		It("should call caller.Call w/ expected args", func() {
			/* arrange */
			dummyString := "dummyString"
			providedDCGOpCall := &model.DCGOpCall{
				DCGBaseCall: model.DCGBaseCall{
					OpHandle: new(data.FakeHandle),
					RootOpID: "providedRootID",
				},
				ChildCallID: "dummyChildCallID",
				ChildCallSCG: &model.SCG{
					Parallel: []*model.SCG{
						{
							Container: &model.SCGContainerCall{},
						},
					},
				},
				Inputs: map[string]*model.Value{
					"dummyScopeName": {String: &dummyString},
				},
				OpID: "providedOpID",
			}

			providedParentCallIDValue := "providedParentCallID"
			providedParentCallID := &providedParentCallIDValue

			fakePubSub := new(pubsub.Fake)
			eventChannel := make(chan model.Event)
			// close eventChannel to trigger immediate return
			close(eventChannel)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeCaller := new(fakeCaller)
			// err to trigger immediate return
			fakeCaller.CallReturns(errors.New("dummyErr"))

			objectUnderTest := _opCaller{
				caller:    fakeCaller,
				callStore: new(fakeCallStore),
				pubSub:    fakePubSub,
			}

			/* act */
			objectUnderTest.Call(
				providedDCGOpCall,
				map[string]*model.Value{},
				providedParentCallID,
				&model.SCGOpCall{},
			)

			/* assert */
			actualChildCallID,
				actualChildCallScope,
				actualChildSCG,
				actualOpRef,
				actualParentCallID,
				actualRootOpID := fakeCaller.CallArgsForCall(0)

			Expect(actualChildCallID).To(Equal(providedDCGOpCall.ChildCallID))
			Expect(actualChildCallScope).To(Equal(providedDCGOpCall.Inputs))
			Expect(actualChildSCG).To(Equal(providedDCGOpCall.ChildCallSCG))
			Expect(actualOpRef).To(Equal(providedDCGOpCall.OpHandle))
			Expect(actualParentCallID).To(Equal(providedParentCallID))
			Expect(actualRootOpID).To(Equal(providedDCGOpCall.RootOpID))
		})
		Context("callStore.Get(callID).IsKilled returns true", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedOpHandleRef := "dummyOpRef"
				fakeOpHandle := new(data.FakeHandle)
				fakeOpHandle.RefReturns(providedOpHandleRef)
				fakeOpHandle.PathReturns(new(string))

				providedDCGOpCall := &model.DCGOpCall{
					DCGBaseCall: model.DCGBaseCall{
						OpHandle: fakeOpHandle,
						RootOpID: "providedRootID",
					},
					OpID: "providedOpID",
				}

				providedSCGOpCall := &model.SCGOpCall{}

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					OpEnded: &model.OpEndedEvent{
						OpID:     providedDCGOpCall.OpID,
						Outcome:  model.OpOutcomeKilled,
						RootOpID: providedDCGOpCall.RootOpID,
						OpRef:    providedOpHandleRef,
						Outputs:  map[string]*model.Value{},
					},
				}

				fakeCallStore := new(fakeCallStore)
				fakeCallStore.GetReturns(model.DCG{IsKilled: true})

				fakeDotYmlGetter := new(dotyml.FakeGetter)
				fakeDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan model.Event)
				// close eventChannel to trigger immediate return
				close(eventChannel)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				objectUnderTest := _opCaller{
					caller:             new(fakeCaller),
					callStore:          fakeCallStore,
					dotYmlGetter:       fakeDotYmlGetter,
					pubSub:             fakePubSub,
					outputsInterpreter: new(outputs.FakeInterpreter),
				}

				/* act */
				objectUnderTest.Call(
					providedDCGOpCall,
					map[string]*model.Value{},
					nil,
					providedSCGOpCall,
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
		Context("callStore.Get(callID).IsKilled returns false", func() {

			Context("caller.Call errs", func() {
				It("should call pubSub.Publish w/ expected args", func() {
					/* arrange */
					providedOpHandleRef := "dummyOpRef"
					fakeOpHandle := new(data.FakeHandle)
					fakeOpHandle.RefReturns(providedOpHandleRef)
					fakeOpHandle.PathReturns(new(string))

					providedDCGOpCall := &model.DCGOpCall{
						DCGBaseCall: model.DCGBaseCall{
							OpHandle: fakeOpHandle,
							RootOpID: "providedRootID",
						},
						OpID: "providedOpId",
					}

					providedSCGOpCall := &model.SCGOpCall{}

					expectedEvent := model.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &model.OpEndedEvent{
							OpID:     providedDCGOpCall.OpID,
							OpRef:    providedOpHandleRef,
							Outcome:  model.OpOutcomeFailed,
							RootOpID: providedDCGOpCall.RootOpID,
							Outputs:  map[string]*model.Value{},
						},
					}

					fakePubSub := new(pubsub.Fake)
					eventChannel := make(chan model.Event)
					// close eventChannel to trigger immediate return
					close(eventChannel)
					fakePubSub.SubscribeReturns(eventChannel, nil)

					fakeCaller := new(fakeCaller)
					fakeCaller.CallReturns(
						errors.New("dummyError"),
					)

					objectUnderTest := _opCaller{
						caller:    fakeCaller,
						callStore: new(fakeCallStore),
						pubSub:    fakePubSub,
					}

					/* act */
					objectUnderTest.Call(
						providedDCGOpCall,
						map[string]*model.Value{},
						nil,
						providedSCGOpCall,
					)

					/* assert */
					actualEvent := fakePubSub.PublishArgsForCall(2)

					// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
					Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
					// set temporal fields to expected vals since they're already asserted
					actualEvent.Timestamp = expectedEvent.Timestamp

					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
		})
		Context("caller.Call didn't error", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedOpHandleRef := "dummyOpRef"
				fakeOpHandle := new(data.FakeHandle)
				fakeOpHandle.RefReturns(providedOpHandleRef)
				fakeOpHandle.PathReturns(new(string))

				providedDCGOpCall := &model.DCGOpCall{
					DCGBaseCall: model.DCGBaseCall{
						OpHandle: fakeOpHandle,
						RootOpID: "providedRootID",
					},
					OpID: "providedOpId",
				}

				expectedOutputName := "expectedOutputName"

				providedSCGOpCall := &model.SCGOpCall{
					Outputs: map[string]string{
						expectedOutputName: "",
					},
				}

				fakeOutputsInterpreter := new(outputs.FakeInterpreter)
				interpretedOutputs := map[string]*model.Value{
					expectedOutputName: new(model.Value),
					// include unbound output to ensure it's not added to scope
					"unexpectedOutputName": new(model.Value),
				}
				fakeOutputsInterpreter.InterpretReturns(interpretedOutputs, nil)

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					OpEnded: &model.OpEndedEvent{
						OpID:     providedDCGOpCall.OpID,
						OpRef:    providedOpHandleRef,
						Outcome:  model.OpOutcomeSucceeded,
						RootOpID: providedDCGOpCall.RootOpID,
						Outputs: map[string]*model.Value{
							expectedOutputName: interpretedOutputs[expectedOutputName],
						},
					},
				}

				fakeDotYmlGetter := new(dotyml.FakeGetter)
				fakeDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan model.Event)
				// close eventChannel to trigger immediate return
				close(eventChannel)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				objectUnderTest := _opCaller{
					caller:             new(fakeCaller),
					callStore:          new(fakeCallStore),
					dotYmlGetter:       fakeDotYmlGetter,
					pubSub:             fakePubSub,
					outputsInterpreter: fakeOutputsInterpreter,
				}

				/* act */
				objectUnderTest.Call(
					providedDCGOpCall,
					map[string]*model.Value{},
					nil,
					providedSCGOpCall,
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
	})
})
