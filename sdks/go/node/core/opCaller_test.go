package core

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	outputsFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/outputs/fakes"
	. "github.com/opctl/opctl/sdks/go/opspec/opfile/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(FakeCallStore),
				new(FakePubSub),
				new(FakeCaller),
				"",
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call pubSub.Publish w/ expected args", func() {
			/* arrange */
			providedOpPath := "providedOpPath"

			providedOpCall := &model.OpCall{
				BaseCall: model.BaseCall{
					OpPath:   providedOpPath,
					RootOpID: "providedRootID",
				},
				OpID: "providedOpId",
			}

			providedOpCallSpec := &model.OpCallSpec{}

			expectedEvent := model.Event{
				Timestamp: time.Now().UTC(),
				OpStarted: &model.OpStarted{
					OpID:     providedOpCall.OpID,
					OpRef:    providedOpPath,
					RootOpID: providedOpCall.RootOpID,
				},
			}

			fakePubSub := new(FakePubSub)
			eventChannel := make(chan model.Event)
			// close eventChannel to trigger immediate return
			close(eventChannel)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeOpFileGetter := new(FakeGetter)
			// err to trigger immediate return
			fakeOpFileGetter.GetReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _opCaller{
				caller:       new(FakeCaller),
				callStore:    new(FakeCallStore),
				opFileGetter: fakeOpFileGetter,
				pubSub:       fakePubSub,
			}

			/* act */
			objectUnderTest.Call(
				context.Background(),
				providedOpCall,
				map[string]*model.Value{},
				nil,
				providedOpCallSpec,
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
			providedOpPath := "providedOpPath"

			dummyString := "dummyString"
			providedCtx := context.Background()
			providedOpCall := &model.OpCall{
				BaseCall: model.BaseCall{
					OpPath:   providedOpPath,
					RootOpID: "providedRootID",
				},
				ChildCallID: "dummyChildCallID",
				ChildCallCallSpec: &model.CallSpec{
					Parallel: &[]*model.CallSpec{
						{
							Container: &model.ContainerCallSpec{},
						},
					},
				},
				Inputs: map[string]*model.Value{
					"dummyScopeName": {String: &dummyString},
				},
				OpID: "providedOpID",
			}

			expectedChildCallScope := map[string]*model.Value{
				"dummyScopeName": providedOpCall.Inputs["dummyScopeName"],
				"/": &model.Value{
					Dir: &providedOpPath,
				},
			}

			fakePubSub := new(FakePubSub)
			eventChannel := make(chan model.Event)
			// close eventChannel to trigger immediate return
			close(eventChannel)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeCaller := new(FakeCaller)

			fakeOpFileGetter := new(FakeGetter)
			// err to trigger immediate return
			fakeOpFileGetter.GetReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _opCaller{
				caller:       fakeCaller,
				callStore:    new(FakeCallStore),
				opFileGetter: fakeOpFileGetter,
				pubSub:       fakePubSub,
			}

			/* act */
			objectUnderTest.Call(
				providedCtx,
				providedOpCall,
				map[string]*model.Value{},
				nil,
				&model.OpCallSpec{},
			)

			/* assert */
			actualCtx,
				actualChildCallID,
				actualChildCallScope,
				actualChildCallSpec,
				actualOpPath,
				actualParentCallID,
				actualRootOpID := fakeCaller.CallArgsForCall(0)

			Expect(actualCtx).To(Not(BeNil()))
			Expect(actualChildCallID).To(Equal(providedOpCall.ChildCallID))
			Expect(actualChildCallScope).To(Equal(expectedChildCallScope))
			Expect(actualChildCallSpec).To(Equal(providedOpCall.ChildCallCallSpec))
			Expect(actualOpPath).To(Equal(providedOpPath))
			Expect(actualParentCallID).To(Equal(&providedOpCall.OpID))
			Expect(actualRootOpID).To(Equal(providedOpCall.RootOpID))
		})
		Context("callStore.Get(callID).IsKilled returns true", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedOpPath := "providedOpPath"

				providedOpCall := &model.OpCall{
					BaseCall: model.BaseCall{
						OpPath:   providedOpPath,
						RootOpID: "providedRootID",
					},
					OpID: "providedOpID",
				}

				providedOpCallSpec := &model.OpCallSpec{}

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					OpEnded: &model.OpEnded{
						OpID:     providedOpCall.OpID,
						Outcome:  model.OpOutcomeKilled,
						RootOpID: providedOpCall.RootOpID,
						OpRef:    providedOpPath,
						Outputs:  map[string]*model.Value{},
					},
				}

				fakeCallStore := new(FakeCallStore)
				fakeCallStore.TryGetReturns(&model.Call{IsKilled: true})

				fakeOpFileGetter := new(FakeGetter)
				fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

				fakePubSub := new(FakePubSub)
				eventChannel := make(chan model.Event)
				// close eventChannel to trigger immediate return
				close(eventChannel)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				objectUnderTest := _opCaller{
					caller:             new(FakeCaller),
					callStore:          fakeCallStore,
					opFileGetter:       fakeOpFileGetter,
					pubSub:             fakePubSub,
					outputsInterpreter: new(outputsFakes.FakeInterpreter),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedOpCall,
					map[string]*model.Value{},
					nil,
					providedOpCallSpec,
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
					providedOpPath := "providedOpPath"

					providedOpCall := &model.OpCall{
						BaseCall: model.BaseCall{
							OpPath:   providedOpPath,
							RootOpID: "providedRootID",
						},
						OpID: "providedOpId",
					}

					providedOpCallSpec := &model.OpCallSpec{}
					errMsg := "errMsg"

					expectedEvent := model.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &model.OpEnded{
							Error: &model.CallEndedError{
								Message: errMsg,
							},
							OpID:     providedOpCall.OpID,
							OpRef:    providedOpPath,
							Outcome:  model.OpOutcomeFailed,
							RootOpID: providedOpCall.RootOpID,
							Outputs:  map[string]*model.Value{},
						},
					}

					fakePubSub := new(FakePubSub)
					eventChannel := make(chan model.Event)
					// close eventChannel to trigger immediate return
					close(eventChannel)
					fakePubSub.SubscribeReturns(eventChannel, nil)

					fakeOpFileGetter := new(FakeGetter)
					// err to trigger immediate return
					fakeOpFileGetter.GetReturns(nil, errors.New(errMsg))

					objectUnderTest := _opCaller{
						caller:       new(FakeCaller),
						callStore:    new(FakeCallStore),
						opFileGetter: fakeOpFileGetter,
						pubSub:       fakePubSub,
					}

					/* act */
					objectUnderTest.Call(
						context.Background(),
						providedOpCall,
						map[string]*model.Value{},
						nil,
						providedOpCallSpec,
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
		Context("caller.Call didn't error", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedOpPath := "providedOpPath"

				providedOpCall := &model.OpCall{
					BaseCall: model.BaseCall{
						OpPath:   providedOpPath,
						RootOpID: "providedRootID",
					},
					OpID: "providedOpId",
				}

				expectedOutputName := "expectedOutputName"

				providedOpCallSpec := &model.OpCallSpec{
					Outputs: map[string]string{
						expectedOutputName: "",
					},
				}

				fakeOutputsInterpreter := new(outputsFakes.FakeInterpreter)
				interpretedOutputs := map[string]*model.Value{
					expectedOutputName: new(model.Value),
					// include unbound output to ensure it's not added to scope
					"unexpectedOutputName": new(model.Value),
				}
				fakeOutputsInterpreter.InterpretReturns(interpretedOutputs, nil)

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					OpEnded: &model.OpEnded{
						OpID:     providedOpCall.OpID,
						OpRef:    providedOpPath,
						Outcome:  model.OpOutcomeSucceeded,
						RootOpID: providedOpCall.RootOpID,
						Outputs: map[string]*model.Value{
							expectedOutputName: interpretedOutputs[expectedOutputName],
						},
					},
				}

				fakeOpFileGetter := new(FakeGetter)
				fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

				fakePubSub := new(FakePubSub)
				eventChannel := make(chan model.Event)
				// close eventChannel to trigger immediate return
				close(eventChannel)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				objectUnderTest := _opCaller{
					caller:             new(FakeCaller),
					callStore:          new(FakeCallStore),
					opFileGetter:       fakeOpFileGetter,
					pubSub:             fakePubSub,
					outputsInterpreter: fakeOutputsInterpreter,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedOpCall,
					map[string]*model.Value{},
					nil,
					providedOpCallSpec,
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
