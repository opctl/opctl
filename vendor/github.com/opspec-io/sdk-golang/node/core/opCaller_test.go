package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/dotyml"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall/outputs"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"time"
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
			providedOpID := "dummyOpID"
			providedRootOpID := "dummyRootOpID"
			providedSCGOpCall := &model.SCGOpCall{
				Ref: "dummyOpRef",
			}

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:       providedOpID,
				OpRef:    providedSCGOpCall.Ref,
				RootOpID: providedRootOpID,
				Op:       &dcgOpDescriptor{},
			}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			fakeOpCallInterpreter := new(opcall.FakeInterpreter)
			// error to trigger immediate return
			fakeOpCallInterpreter.InterpretReturns(nil, errors.New("dummyError"))

			objectUnderTest := _opCaller{
				opCallInterpreter: fakeOpCallInterpreter,
				pubSub:            new(pubsub.Fake),
				dcgNodeRepo:       fakeDCGNodeRepo,
				caller:            new(fakeCaller),
			}

			/* act */
			objectUnderTest.Call(
				map[string]*model.Value{},
				providedOpID,
				new(data.FakeHandle),
				providedRootOpID,
				providedSCGOpCall,
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
		It("should call opCall.Construct w/ expected args & return errors", func() {
			/* arrange */
			providedScope := map[string]*model.Value{}
			providedOpID := "dummyOpID"
			providedRootOpID := "dummyRootOpID"
			providedOpHandle := new(data.FakeHandle)
			providedSCGOpCall := &model.SCGOpCall{
				Ref: "dummyOpRef",
			}

			fakeOpCallInterpreter := new(opcall.FakeInterpreter)

			// error to trigger immediate return
			expectedErr := errors.New("dummyError")
			fakeOpCallInterpreter.InterpretReturns(
				&model.DCGOpCall{
					DCGBaseCall: model.DCGBaseCall{},
				},
				expectedErr,
			)

			objectUnderTest := _opCaller{
				opCallInterpreter: fakeOpCallInterpreter,
				pubSub:            new(pubsub.Fake),
				dcgNodeRepo:       new(fakeDCGNodeRepo),
				caller:            new(fakeCaller),
			}

			/* act */
			actualErr := objectUnderTest.Call(
				providedScope,
				providedOpID,
				providedOpHandle,
				providedRootOpID,
				providedSCGOpCall,
			)

			/* assert */
			actualScope,
				actualSCGOpCall,
				actualOpID,
				actualOpHandle,
				actualRootOpID := fakeOpCallInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualSCGOpCall).To(Equal(providedSCGOpCall))
			Expect(actualOpID).To(Equal(providedOpID))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
			Expect(actualRootOpID).To(Equal(providedRootOpID))

			Expect(actualErr).To(Equal(expectedErr))
		})
		Context("opCall.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				providedOpID := "dummyOpID"
				providedRootOpID := "dummyRootOpID"
				providedSCGOpCall := &model.SCGOpCall{
					Ref: "dummyOpRef",
				}

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					OpErred: &model.OpErredEvent{
						Msg:      "dummyError",
						OpID:     providedOpID,
						OpRef:    providedSCGOpCall.Ref,
						RootOpID: providedRootOpID,
					},
				}

				fakeOpCallInterpreter := new(opcall.FakeInterpreter)
				fakeOpCallInterpreter.InterpretReturns(
					nil,
					errors.New(expectedEvent.OpErred.Msg),
				)

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := _opCaller{
					opCallInterpreter: fakeOpCallInterpreter,
					pubSub:            fakePubSub,
					dcgNodeRepo:       fakeDCGNodeRepo,
					caller:            new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
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
		})
		Context("opCall.Construct doesn't error", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedOpID := "dummyOpID"
				providedRootOpID := "dummyRootOpID"
				providedSCGOpCall := &model.SCGOpCall{
					Ref: "dummyOpRef",
				}

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					OpStarted: &model.OpStartedEvent{
						OpID:     providedOpID,
						OpRef:    providedSCGOpCall.Ref,
						RootOpID: providedRootOpID,
					},
				}

				fakeOpCallInterpreter := new(opcall.FakeInterpreter)
				fakeOpCallInterpreter.InterpretReturns(
					&model.DCGOpCall{
						DCGBaseCall: model.DCGBaseCall{},
					},
					nil,
				)

				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan model.Event)
				// close eventChannel to trigger immediate return
				close(eventChannel)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				fakeCaller := new(fakeCaller)
				// err to trigger immediate return
				fakeCaller.CallReturns(errors.New("dummyError"))

				objectUnderTest := _opCaller{
					opCallInterpreter: fakeOpCallInterpreter,
					pubSub:            fakePubSub,
					dcgNodeRepo:       new(fakeDCGNodeRepo),
					caller:            fakeCaller,
				}

				/* act */
				objectUnderTest.Call(
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
				dcgOpCall := model.DCGOpCall{
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
				fakeOpCallInterpreter := new(opcall.FakeInterpreter)
				fakeOpCallInterpreter.InterpretReturns(
					&dcgOpCall,
					nil,
				)

				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan model.Event)
				// close eventChannel to trigger immediate return
				close(eventChannel)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				fakeCaller := new(fakeCaller)
				// err to trigger immediate return
				fakeCaller.CallReturns(errors.New("dummyErr"))

				objectUnderTest := _opCaller{
					opCallInterpreter: fakeOpCallInterpreter,
					pubSub:            fakePubSub,
					dcgNodeRepo:       new(fakeDCGNodeRepo),
					caller:            fakeCaller,
				}

				/* act */
				objectUnderTest.Call(
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
				providedRootOpID := "dummyRootOpID"

				fakeDataHandle := new(data.FakeHandle)
				fakeDataHandle.PathReturns(new(string))

				fakeOpCallInterpreter := new(opcall.FakeInterpreter)
				fakeOpCallInterpreter.InterpretReturns(
					&model.DCGOpCall{
						DCGBaseCall: model.DCGBaseCall{
							OpHandle: fakeDataHandle,
						},
					},
					nil,
				)

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
					opCallInterpreter:  fakeOpCallInterpreter,
					outputsInterpreter: new(outputs.FakeInterpreter),
					dcgNodeRepo:        fakeDCGNodeRepo,
					caller:             new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
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
					providedOpID := "dummyOpID"
					providedRootOpID := "dummyRootOpID"
					providedSCGOpCall := &model.SCGOpCall{
						Ref: "dummyOpRef",
					}

					expectedEvent := model.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &model.OpEndedEvent{
							OpID:     providedOpID,
							Outcome:  model.OpOutcomeKilled,
							RootOpID: providedRootOpID,
							OpRef:    providedSCGOpCall.Ref,
						},
					}

					fakeDataHandle := new(data.FakeHandle)
					fakeDataHandle.PathReturns(new(string))

					fakeOpCallInterpreter := new(opcall.FakeInterpreter)
					fakeOpCallInterpreter.InterpretReturns(
						&model.DCGOpCall{
							DCGBaseCall: model.DCGBaseCall{
								OpHandle: fakeDataHandle,
							},
						},
						nil,
					)

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
						opCallInterpreter:  fakeOpCallInterpreter,
						outputsInterpreter: new(outputs.FakeInterpreter),
						dcgNodeRepo:        new(fakeDCGNodeRepo),
						caller:             new(fakeCaller),
					}

					/* act */
					objectUnderTest.Call(
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
					providedOpID := "dummyOpID"

					fakeDataHandle := new(data.FakeHandle)
					fakeDataHandle.PathReturns(new(string))

					fakeOpCallInterpreter := new(opcall.FakeInterpreter)
					fakeOpCallInterpreter.InterpretReturns(
						&model.DCGOpCall{
							DCGBaseCall: model.DCGBaseCall{
								OpHandle: fakeDataHandle,
							},
						},
						nil,
					)

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
						opCallInterpreter:  fakeOpCallInterpreter,
						outputsInterpreter: new(outputs.FakeInterpreter),
						dcgNodeRepo:        fakeDCGNodeRepo,
						caller:             new(fakeCaller),
					}

					/* act */
					objectUnderTest.Call(
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
						providedOpID := "dummyOpID"
						providedRootOpID := "dummyRootOpID"
						providedSCGOpCall := &model.SCGOpCall{
							Ref: "dummyOpRef",
						}

						expectedEvent := model.Event{
							Timestamp: time.Now().UTC(),
							OpEnded: &model.OpEndedEvent{
								OpID:     providedOpID,
								OpRef:    providedSCGOpCall.Ref,
								Outcome:  model.OpOutcomeFailed,
								RootOpID: providedRootOpID,
								Outputs:  map[string]*model.Value{},
							},
						}

						fakeOpCallInterpreter := new(opcall.FakeInterpreter)
						fakeOpCallInterpreter.InterpretReturns(
							&model.DCGOpCall{
								DCGBaseCall: model.DCGBaseCall{},
							},
							nil,
						)

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
							pubSub:            fakePubSub,
							opCallInterpreter: fakeOpCallInterpreter,
							dcgNodeRepo:       fakeDCGNodeRepo,
							caller:            fakeCaller,
						}

						/* act */
						objectUnderTest.Call(
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
						providedOpID := "dummyOpID"
						providedRootOpID := "dummyRootOpID"
						expectedOutputName := "expectedOutputName"

						providedSCGOpCall := &model.SCGOpCall{
							Ref: "dummyOpRef",
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
								OpRef:    providedSCGOpCall.Ref,
								Outcome:  model.OpOutcomeSucceeded,
								RootOpID: providedRootOpID,
								Outputs: map[string]*model.Value{
									expectedOutputName: interpretedOutputs[expectedOutputName],
								},
							},
						}

						fakeDataHandle := new(data.FakeHandle)
						fakeDataHandle.PathReturns(new(string))

						fakeOpCallInterpreter := new(opcall.FakeInterpreter)
						fakeOpCallInterpreter.InterpretReturns(
							&model.DCGOpCall{
								DCGBaseCall: model.DCGBaseCall{
									OpHandle: fakeDataHandle,
								},
							},
							nil,
						)

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
							opCallInterpreter:  fakeOpCallInterpreter,
							outputsInterpreter: fakeOutputsInterpreter,
							dcgNodeRepo:        fakeDCGNodeRepo,
							caller:             new(fakeCaller),
						}

						/* act */
						objectUnderTest.Call(
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
})
