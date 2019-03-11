package core

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/outputs"
	"github.com/opctl/sdk-golang/opspec/opfile"
	"github.com/opctl/sdk-golang/util/pubsub"
)

var _ = Context("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(pubsub.Fake),
				newDCGNodeRepo(),
				new(fakeCaller),
				"",
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call dcgNodeRepo.add w/ expected args", func() {
			/* arrange */
			providedOpHandleRef := "dummyOpRef"
			fakeOpHandle := new(data.FakeHandle)
			fakeOpHandle.RefReturns(providedOpHandleRef)

			providedDCGOpCall := &model.DCGOpCall{
				DCGBaseCall: model.DCGBaseCall{
					OpHandle: fakeOpHandle,
				},
			}

			providedOpID := "dummyOpID"
			providedRootOpID := "dummyRootOpID"
			providedSCGOpCall := &model.SCGOpCall{}

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:       providedOpID,
				OpRef:    providedOpHandleRef,
				RootOpID: providedRootOpID,
				Op:       &dcgOpDescriptor{},
			}

			fakePubSub := new(pubsub.Fake)
			eventChannel := make(chan model.Event)
			// close eventChannel to trigger immediate return
			close(eventChannel)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeCaller := new(fakeCaller)
			// err to trigger immediate return
			fakeCaller.CallReturns(errors.New("dummyError"))

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			objectUnderTest := _opCaller{
				pubSub:      fakePubSub,
				dcgNodeRepo: fakeDCGNodeRepo,
				caller:      fakeCaller,
			}

			/* act */
			objectUnderTest.Call(
				providedDCGOpCall,
				map[string]*model.Value{},
				providedOpID,
				new(data.FakeHandle),
				providedRootOpID,
				providedSCGOpCall,
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
		It("should call pubSub.Publish w/ expected args", func() {
			/* arrange */
			providedOpHandleRef := "dummyOpRef"
			fakeOpHandle := new(data.FakeHandle)
			fakeOpHandle.RefReturns(providedOpHandleRef)

			providedDCGOpCall := &model.DCGOpCall{
				DCGBaseCall: model.DCGBaseCall{
					OpHandle: fakeOpHandle,
				},
			}

			providedOpID := "dummyOpID"
			providedRootOpID := "dummyRootOpID"
			providedSCGOpCall := &model.SCGOpCall{}

			expectedEvent := model.Event{
				Timestamp: time.Now().UTC(),
				OpStarted: &model.OpStartedEvent{
					OpID:     providedOpID,
					OpRef:    providedOpHandleRef,
					RootOpID: providedRootOpID,
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
				pubSub:      fakePubSub,
				dcgNodeRepo: new(fakeDCGNodeRepo),
				caller:      fakeCaller,
			}

			/* act */
			objectUnderTest.Call(
				providedDCGOpCall,
				map[string]*model.Value{},
				providedOpID,
				new(data.FakeHandle),
				providedRootOpID,
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
			providedRootOpID := "dummyRootOpID"

			dummyString := "dummyString"
			dcgOpCall := &model.DCGOpCall{
				DCGBaseCall: model.DCGBaseCall{
					OpHandle: new(data.FakeHandle),
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
			}

			fakePubSub := new(pubsub.Fake)
			eventChannel := make(chan model.Event)
			// close eventChannel to trigger immediate return
			close(eventChannel)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeCaller := new(fakeCaller)
			// err to trigger immediate return
			fakeCaller.CallReturns(errors.New("dummyErr"))

			objectUnderTest := _opCaller{
				pubSub:      fakePubSub,
				dcgNodeRepo: new(fakeDCGNodeRepo),
				caller:      fakeCaller,
			}

			/* act */
			objectUnderTest.Call(
				dcgOpCall,
				map[string]*model.Value{},
				"dummyOpID",
				new(data.FakeHandle),
				providedRootOpID,
				&model.SCGOpCall{},
			)

			/* assert */
			actualChildCallID,
				actualChildCallScope,
				actualChildSCG,
				actualOpRef,
				actualRootOpID := fakeCaller.CallArgsForCall(0)

			Expect(actualChildCallID).To(Equal(dcgOpCall.ChildCallID))
			Expect(actualChildCallScope).To(Equal(dcgOpCall.Inputs))
			Expect(actualChildSCG).To(Equal(dcgOpCall.ChildCallSCG))
			Expect(actualOpRef).To(Equal(dcgOpCall.OpHandle))
			Expect(actualRootOpID).To(Equal(providedRootOpID))
		})
		It("should call dcgNodeRepo.GetIfExists w/ expected args", func() {
			/* arrange */
			fakeHandle := new(data.FakeHandle)
			fakeHandle.PathReturns(new(string))

			providedDCGOpCall := &model.DCGOpCall{
				DCGBaseCall: model.DCGBaseCall{
					OpHandle: fakeHandle,
				},
			}

			providedRootOpID := "dummyRootOpID"

			fakeDotYmlGetter := new(dotyml.FakeGetter)
			fakeDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

			fakePubSub := new(pubsub.Fake)
			eventChannel := make(chan model.Event)
			// close eventChannel to trigger immediate return
			close(eventChannel)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			objectUnderTest := _opCaller{
				dotYmlGetter:       fakeDotYmlGetter,
				pubSub:             fakePubSub,
				outputsInterpreter: new(outputs.FakeInterpreter),
				dcgNodeRepo:        fakeDCGNodeRepo,
				caller:             new(fakeCaller),
			}

			/* act */
			objectUnderTest.Call(
				providedDCGOpCall,
				map[string]*model.Value{},
				"dummyOpID",
				new(data.FakeHandle),
				providedRootOpID,
				&model.SCGOpCall{},
			)

			/* assert */
			Expect(fakeDCGNodeRepo.GetIfExistsArgsForCall(0)).To(Equal(providedRootOpID))
		})
		Context("dcgNodeRepo.GetIfExists returns nil", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedOpHandleRef := "dummyOpRef"
				fakeOpHandle := new(data.FakeHandle)
				fakeOpHandle.RefReturns(providedOpHandleRef)
				fakeOpHandle.PathReturns(new(string))

				providedDCGOpCall := &model.DCGOpCall{
					DCGBaseCall: model.DCGBaseCall{
						OpHandle: fakeOpHandle,
					},
				}

				providedOpID := "dummyOpID"
				providedRootOpID := "dummyRootOpID"
				providedSCGOpCall := &model.SCGOpCall{}

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					OpEnded: &model.OpEndedEvent{
						OpID:     providedOpID,
						Outcome:  model.OpOutcomeKilled,
						RootOpID: providedRootOpID,
						OpRef:    providedOpHandleRef,
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
					dotYmlGetter:       fakeDotYmlGetter,
					pubSub:             fakePubSub,
					outputsInterpreter: new(outputs.FakeInterpreter),
					dcgNodeRepo:        new(fakeDCGNodeRepo),
					caller:             new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
					providedDCGOpCall,
					map[string]*model.Value{},
					providedOpID,
					new(data.FakeHandle),
					providedRootOpID,
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
		Context("dcgNodeRepo.GetIfExists doesn't return nil", func() {
			It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
				/* arrange */
				fakeHandle := new(data.FakeHandle)
				fakeHandle.PathReturns(new(string))

				providedDCGOpCall := &model.DCGOpCall{
					DCGBaseCall: model.DCGBaseCall{
						OpHandle: fakeHandle,
					},
				}

				providedOpID := "dummyOpID"

				fakeDotYmlGetter := new(dotyml.FakeGetter)
				fakeDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan model.Event)
				// close eventChannel to trigger immediate return
				close(eventChannel)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				objectUnderTest := _opCaller{
					dotYmlGetter:       fakeDotYmlGetter,
					pubSub:             fakePubSub,
					outputsInterpreter: new(outputs.FakeInterpreter),
					dcgNodeRepo:        fakeDCGNodeRepo,
					caller:             new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
					providedDCGOpCall,
					map[string]*model.Value{},
					providedOpID,
					new(data.FakeHandle),
					"dummyRootOpID",
					&model.SCGOpCall{},
				)

				/* assert */
				Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedOpID))
			})
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
						},
					}

					providedOpID := "dummyOpID"
					providedRootOpID := "dummyRootOpID"
					providedSCGOpCall := &model.SCGOpCall{}

					expectedEvent := model.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &model.OpEndedEvent{
							OpID:     providedOpID,
							OpRef:    providedOpHandleRef,
							Outcome:  model.OpOutcomeFailed,
							RootOpID: providedRootOpID,
							Outputs:  map[string]*model.Value{},
						},
					}

					fakeDCGNodeRepo := new(fakeDCGNodeRepo)
					fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

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
						pubSub:      fakePubSub,
						dcgNodeRepo: fakeDCGNodeRepo,
						caller:      fakeCaller,
					}

					/* act */
					objectUnderTest.Call(
						providedDCGOpCall,
						map[string]*model.Value{},
						providedOpID,
						new(data.FakeHandle),
						providedRootOpID,
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
						},
					}

					providedOpID := "dummyOpID"
					providedRootOpID := "dummyRootOpID"
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
							OpID:     providedOpID,
							OpRef:    providedOpHandleRef,
							Outcome:  model.OpOutcomeSucceeded,
							RootOpID: providedRootOpID,
							Outputs: map[string]*model.Value{
								expectedOutputName: interpretedOutputs[expectedOutputName],
							},
						},
					}

					fakeDotYmlGetter := new(dotyml.FakeGetter)
					fakeDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

					fakeDCGNodeRepo := new(fakeDCGNodeRepo)
					fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

					fakePubSub := new(pubsub.Fake)
					eventChannel := make(chan model.Event)
					// close eventChannel to trigger immediate return
					close(eventChannel)
					fakePubSub.SubscribeReturns(eventChannel, nil)

					objectUnderTest := _opCaller{
						dotYmlGetter:       fakeDotYmlGetter,
						pubSub:             fakePubSub,
						outputsInterpreter: fakeOutputsInterpreter,
						dcgNodeRepo:        fakeDCGNodeRepo,
						caller:             new(fakeCaller),
					}

					/* act */
					objectUnderTest.Call(
						providedDCGOpCall,
						map[string]*model.Value{},
						providedOpID,
						new(data.FakeHandle),
						providedRootOpID,
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
})
