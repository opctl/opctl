package core

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/outputs"
	dotyml "github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/opctl/opctl/sdks/go/types"
	"github.com/opctl/opctl/sdks/go/util/pubsub"
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

			providedDCGOpCall := &types.DCGOpCall{
				DCGBaseCall: types.DCGBaseCall{
					OpHandle: fakeOpHandle,
					RootOpID: "providedRootID",
				},
				OpID: "providedOpId",
			}

			providedSCGOpCall := &types.SCGOpCall{}

			expectedEvent := types.Event{
				Timestamp: time.Now().UTC(),
				OpStarted: &types.OpStartedEvent{
					OpID:     providedDCGOpCall.OpID,
					OpRef:    providedOpHandleRef,
					RootOpID: providedDCGOpCall.RootOpID,
				},
			}

			fakePubSub := new(pubsub.Fake)
			eventChannel := make(chan types.Event)
			// close eventChannel to trigger immediate return
			close(eventChannel)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeDotYmlGetter := new(dotyml.FakeGetter)
			// err to trigger immediate return
			fakeDotYmlGetter.GetReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _opCaller{
				caller:       new(fakeCaller),
				callStore:    new(fakeCallStore),
				dotYmlGetter: fakeDotYmlGetter,
				pubSub:       fakePubSub,
			}

			/* act */
			objectUnderTest.Call(
				context.Background(),
				providedDCGOpCall,
				map[string]*types.Value{},
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
			providedCtx := context.Background()
			providedDCGOpCall := &types.DCGOpCall{
				DCGBaseCall: types.DCGBaseCall{
					OpHandle: new(data.FakeHandle),
					RootOpID: "providedRootID",
				},
				ChildCallID: "dummyChildCallID",
				ChildCallSCG: &types.SCG{
					Parallel: []*types.SCG{
						{
							Container: &types.SCGContainerCall{},
						},
					},
				},
				Inputs: map[string]*types.Value{
					"dummyScopeName": {String: &dummyString},
				},
				OpID: "providedOpID",
			}

			fakePubSub := new(pubsub.Fake)
			eventChannel := make(chan types.Event)
			// close eventChannel to trigger immediate return
			close(eventChannel)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeCaller := new(fakeCaller)

			fakeDotYmlGetter := new(dotyml.FakeGetter)
			// err to trigger immediate return
			fakeDotYmlGetter.GetReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _opCaller{
				caller:       fakeCaller,
				callStore:    new(fakeCallStore),
				dotYmlGetter: fakeDotYmlGetter,
				pubSub:       fakePubSub,
			}

			/* act */
			objectUnderTest.Call(
				providedCtx,
				providedDCGOpCall,
				map[string]*types.Value{},
				nil,
				&types.SCGOpCall{},
			)

			/* assert */
			actualCtx,
				actualChildCallID,
				actualChildCallScope,
				actualChildSCG,
				actualOpRef,
				actualParentCallID,
				actualRootOpID := fakeCaller.CallArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualChildCallID).To(Equal(providedDCGOpCall.ChildCallID))
			Expect(actualChildCallScope).To(Equal(providedDCGOpCall.Inputs))
			Expect(actualChildSCG).To(Equal(providedDCGOpCall.ChildCallSCG))
			Expect(actualOpRef).To(Equal(providedDCGOpCall.OpHandle))
			Expect(actualParentCallID).To(Equal(&providedDCGOpCall.OpID))
			Expect(actualRootOpID).To(Equal(providedDCGOpCall.RootOpID))
		})
		Context("callStore.Get(callID).IsKilled returns true", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedOpHandleRef := "dummyOpRef"
				fakeOpHandle := new(data.FakeHandle)
				fakeOpHandle.RefReturns(providedOpHandleRef)
				fakeOpHandle.PathReturns(new(string))

				providedDCGOpCall := &types.DCGOpCall{
					DCGBaseCall: types.DCGBaseCall{
						OpHandle: fakeOpHandle,
						RootOpID: "providedRootID",
					},
					OpID: "providedOpID",
				}

				providedSCGOpCall := &types.SCGOpCall{}

				expectedEvent := types.Event{
					Timestamp: time.Now().UTC(),
					OpEnded: &types.OpEndedEvent{
						OpID:     providedDCGOpCall.OpID,
						Outcome:  types.OpOutcomeKilled,
						RootOpID: providedDCGOpCall.RootOpID,
						OpRef:    providedOpHandleRef,
						Outputs:  map[string]*types.Value{},
					},
				}

				fakeCallStore := new(fakeCallStore)
				fakeCallStore.GetReturns(types.DCG{IsKilled: true})

				fakeDotYmlGetter := new(dotyml.FakeGetter)
				fakeDotYmlGetter.GetReturns(&types.OpDotYml{}, nil)

				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan types.Event)
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
					context.Background(),
					providedDCGOpCall,
					map[string]*types.Value{},
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

					providedDCGOpCall := &types.DCGOpCall{
						DCGBaseCall: types.DCGBaseCall{
							OpHandle: fakeOpHandle,
							RootOpID: "providedRootID",
						},
						OpID: "providedOpId",
					}

					providedSCGOpCall := &types.SCGOpCall{}
					errMsg := "errMsg"

					expectedEvent := types.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &types.OpEndedEvent{
							Error: &types.CallEndedEventError{
								Message: errMsg,
							},
							OpID:     providedDCGOpCall.OpID,
							OpRef:    providedOpHandleRef,
							Outcome:  types.OpOutcomeFailed,
							RootOpID: providedDCGOpCall.RootOpID,
							Outputs:  map[string]*types.Value{},
						},
					}

					fakePubSub := new(pubsub.Fake)
					eventChannel := make(chan types.Event)
					// close eventChannel to trigger immediate return
					close(eventChannel)
					fakePubSub.SubscribeReturns(eventChannel, nil)

					fakeDotYmlGetter := new(dotyml.FakeGetter)
					// err to trigger immediate return
					fakeDotYmlGetter.GetReturns(nil, errors.New(errMsg))

					objectUnderTest := _opCaller{
						caller:       new(fakeCaller),
						callStore:    new(fakeCallStore),
						dotYmlGetter: fakeDotYmlGetter,
						pubSub:       fakePubSub,
					}

					/* act */
					objectUnderTest.Call(
						context.Background(),
						providedDCGOpCall,
						map[string]*types.Value{},
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

				providedDCGOpCall := &types.DCGOpCall{
					DCGBaseCall: types.DCGBaseCall{
						OpHandle: fakeOpHandle,
						RootOpID: "providedRootID",
					},
					OpID: "providedOpId",
				}

				expectedOutputName := "expectedOutputName"

				providedSCGOpCall := &types.SCGOpCall{
					Outputs: map[string]string{
						expectedOutputName: "",
					},
				}

				fakeOutputsInterpreter := new(outputs.FakeInterpreter)
				interpretedOutputs := map[string]*types.Value{
					expectedOutputName: new(types.Value),
					// include unbound output to ensure it's not added to scope
					"unexpectedOutputName": new(types.Value),
				}
				fakeOutputsInterpreter.InterpretReturns(interpretedOutputs, nil)

				expectedEvent := types.Event{
					Timestamp: time.Now().UTC(),
					OpEnded: &types.OpEndedEvent{
						OpID:     providedDCGOpCall.OpID,
						OpRef:    providedOpHandleRef,
						Outcome:  types.OpOutcomeSucceeded,
						RootOpID: providedDCGOpCall.RootOpID,
						Outputs: map[string]*types.Value{
							expectedOutputName: interpretedOutputs[expectedOutputName],
						},
					},
				}

				fakeDotYmlGetter := new(dotyml.FakeGetter)
				fakeDotYmlGetter.GetReturns(&types.OpDotYml{}, nil)

				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan types.Event)
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
					context.Background(),
					providedDCGOpCall,
					map[string]*types.Value{},
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
