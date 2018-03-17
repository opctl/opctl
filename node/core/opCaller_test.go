package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall/outputs"
	"github.com/opspec-io/sdk-golang/pkg"
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
			providedOpId := "dummyOpId"
			providedRootOpId := "dummyRootOpId"
			providedSCGOpCall := &model.SCGOpCall{
				Pkg: &model.SCGOpCallPkg{
					Ref: "dummyPkgRef",
				}}

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:       providedOpId,
				PkgRef:   providedSCGOpCall.Pkg.Ref,
				RootOpId: providedRootOpId,
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
				providedOpId,
				new(data.FakeHandle),
				providedRootOpId,
				providedSCGOpCall,
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
		It("should call opCall.Construct w/ expected args & return errors", func() {
			/* arrange */
			providedScope := map[string]*model.Value{}
			providedOpId := "dummyOpId"
			providedRootOpId := "dummyRootOpId"
			providedOpDirHandle := new(data.FakeHandle)
			providedSCGOpCall := &model.SCGOpCall{
				Pkg: &model.SCGOpCallPkg{
					Ref: "dummyPkgRef",
				},
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
				providedOpId,
				providedOpDirHandle,
				providedRootOpId,
				providedSCGOpCall,
			)

			/* assert */
			actualScope,
				actualSCGOpCall,
				actualOpId,
				actualOpDirHandle,
				actualRootOpId := fakeOpCallInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualSCGOpCall).To(Equal(providedSCGOpCall))
			Expect(actualOpId).To(Equal(providedOpId))
			Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))
			Expect(actualRootOpId).To(Equal(providedRootOpId))

			Expect(actualErr).To(Equal(expectedErr))
		})
		Context("opCall.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				providedOpId := "dummyOpId"
				providedRootOpId := "dummyRootOpId"
				providedSCGOpCall := &model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref: "dummyPkgRef",
					},
				}

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					OpErred: &model.OpErredEvent{
						Msg:      "dummyError",
						OpId:     providedOpId,
						PkgRef:   providedSCGOpCall.Pkg.Ref,
						RootOpId: providedRootOpId,
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
					providedOpId,
					new(data.FakeHandle),
					providedRootOpId,
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
				providedOpId := "dummyOpId"
				providedRootOpId := "dummyRootOpId"
				providedSCGOpCall := &model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref: "dummyPkgRef",
					},
				}

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					OpStarted: &model.OpStartedEvent{
						OpId:     providedOpId,
						PkgRef:   providedSCGOpCall.Pkg.Ref,
						RootOpId: providedRootOpId,
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
					providedOpId,
					new(data.FakeHandle),
					providedRootOpId,
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
				providedRootOpId := "dummyRootOpId"

				dummyString := "dummyString"
				dcgOpCall := model.DCGOpCall{
					DCGBaseCall: model.DCGBaseCall{
						DataHandle: new(data.FakeHandle),
					},
					ChildCallId: "dummyChildCallId",
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
					"dummyOpId",
					new(data.FakeHandle),
					providedRootOpId,
					&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
				)

				/* assert */
				actualChildCallId,
					actualChildCallScope,
					actualChildSCG,
					actualPkgRef,
					actualRootOpId := fakeCaller.CallArgsForCall(0)

				Expect(actualChildCallId).To(Equal(dcgOpCall.ChildCallId))
				Expect(actualChildCallScope).To(Equal(dcgOpCall.Inputs))
				Expect(actualChildSCG).To(Equal(dcgOpCall.ChildCallSCG))
				Expect(actualPkgRef).To(Equal(dcgOpCall.DataHandle))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			})
			It("should call dcgNodeRepo.GetIfExists w/ expected args", func() {
				/* arrange */
				providedRootOpId := "dummyRootOpId"

				fakeDataHandle := new(data.FakeHandle)
				fakeDataHandle.PathReturns(new(string))

				fakeOpCallInterpreter := new(opcall.FakeInterpreter)
				fakeOpCallInterpreter.InterpretReturns(
					&model.DCGOpCall{
						DCGBaseCall: model.DCGBaseCall{
							DataHandle: fakeDataHandle,
						},
					},
					nil,
				)

				fakePkg := new(pkg.Fake)
				fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan model.Event)
				// close eventChannel to trigger immediate return
				close(eventChannel)
				fakePubSub.SubscribeReturns(eventChannel, nil)

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)

				objectUnderTest := _opCaller{
					pkg:                fakePkg,
					pubSub:             fakePubSub,
					opCallInterpreter:  fakeOpCallInterpreter,
					outputsInterpreter: new(outputs.FakeInterpreter),
					dcgNodeRepo:        fakeDCGNodeRepo,
					caller:             new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
					map[string]*model.Value{},
					"dummyOpId",
					new(data.FakeHandle),
					providedRootOpId,
					&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
				)

				/* assert */
				Expect(fakeDCGNodeRepo.GetIfExistsArgsForCall(0)).To(Equal(providedRootOpId))
			})
			Context("dcgNodeRepo.GetIfExists returns nil", func() {
				It("should call pubSub.Publish w/ expected args", func() {
					/* arrange */
					providedOpId := "dummyOpId"
					providedRootOpId := "dummyRootOpId"
					providedSCGOpCall := &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: "dummyPkgRef",
						},
					}

					expectedEvent := model.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &model.OpEndedEvent{
							OpId:     providedOpId,
							Outcome:  model.OpOutcomeKilled,
							RootOpId: providedRootOpId,
							PkgRef:   providedSCGOpCall.Pkg.Ref,
						},
					}

					fakeDataHandle := new(data.FakeHandle)
					fakeDataHandle.PathReturns(new(string))

					fakeOpCallInterpreter := new(opcall.FakeInterpreter)
					fakeOpCallInterpreter.InterpretReturns(
						&model.DCGOpCall{
							DCGBaseCall: model.DCGBaseCall{
								DataHandle: fakeDataHandle,
							},
						},
						nil,
					)

					fakePkg := new(pkg.Fake)
					fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

					fakePubSub := new(pubsub.Fake)
					eventChannel := make(chan model.Event)
					// close eventChannel to trigger immediate return
					close(eventChannel)
					fakePubSub.SubscribeReturns(eventChannel, nil)

					objectUnderTest := _opCaller{
						pkg:                fakePkg,
						pubSub:             fakePubSub,
						opCallInterpreter:  fakeOpCallInterpreter,
						outputsInterpreter: new(outputs.FakeInterpreter),
						dcgNodeRepo:        new(fakeDCGNodeRepo),
						caller:             new(fakeCaller),
					}

					/* act */
					objectUnderTest.Call(
						map[string]*model.Value{},
						providedOpId,
						new(data.FakeHandle),
						providedRootOpId,
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
					providedOpId := "dummyOpId"

					fakeDataHandle := new(data.FakeHandle)
					fakeDataHandle.PathReturns(new(string))

					fakeOpCallInterpreter := new(opcall.FakeInterpreter)
					fakeOpCallInterpreter.InterpretReturns(
						&model.DCGOpCall{
							DCGBaseCall: model.DCGBaseCall{
								DataHandle: fakeDataHandle,
							},
						},
						nil,
					)

					fakePkg := new(pkg.Fake)
					fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

					fakeDCGNodeRepo := new(fakeDCGNodeRepo)
					fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

					fakePubSub := new(pubsub.Fake)
					eventChannel := make(chan model.Event)
					// close eventChannel to trigger immediate return
					close(eventChannel)
					fakePubSub.SubscribeReturns(eventChannel, nil)

					objectUnderTest := _opCaller{
						pkg:                fakePkg,
						pubSub:             fakePubSub,
						opCallInterpreter:  fakeOpCallInterpreter,
						outputsInterpreter: new(outputs.FakeInterpreter),
						dcgNodeRepo:        fakeDCGNodeRepo,
						caller:             new(fakeCaller),
					}

					/* act */
					objectUnderTest.Call(
						map[string]*model.Value{},
						providedOpId,
						new(data.FakeHandle),
						"dummyRootOpId",
						&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
					)

					/* assert */
					Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedOpId))
				})
				Context("caller.Call errs", func() {
					It("should call pubSub.Publish w/ expected args", func() {
						/* arrange */
						providedOpId := "dummyOpId"
						providedRootOpId := "dummyRootOpId"
						providedSCGOpCall := &model.SCGOpCall{
							Pkg: &model.SCGOpCallPkg{
								Ref: "dummyPkgRef",
							},
						}

						expectedEvent := model.Event{
							Timestamp: time.Now().UTC(),
							OpEnded: &model.OpEndedEvent{
								OpId:     providedOpId,
								PkgRef:   providedSCGOpCall.Pkg.Ref,
								Outcome:  model.OpOutcomeFailed,
								RootOpId: providedRootOpId,
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
							providedOpId,
							new(data.FakeHandle),
							providedRootOpId,
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
						providedOpId := "dummyOpId"
						providedRootOpId := "dummyRootOpId"
						providedSCGOpCall := &model.SCGOpCall{
							Pkg: &model.SCGOpCallPkg{
								Ref: "dummyPkgRef",
							},
						}

						fakeOutputsInterpreter := new(outputs.FakeInterpreter)
						interpretedOutputs := map[string]*model.Value{"dummyOutputName": new(model.Value)}
						fakeOutputsInterpreter.InterpretReturns(interpretedOutputs, nil)

						expectedEvent := model.Event{
							Timestamp: time.Now().UTC(),
							OpEnded: &model.OpEndedEvent{
								OpId:     providedOpId,
								PkgRef:   providedSCGOpCall.Pkg.Ref,
								Outcome:  model.OpOutcomeSucceeded,
								RootOpId: providedRootOpId,
								Outputs:  interpretedOutputs,
							},
						}

						fakeDataHandle := new(data.FakeHandle)
						fakeDataHandle.PathReturns(new(string))

						fakeOpCallInterpreter := new(opcall.FakeInterpreter)
						fakeOpCallInterpreter.InterpretReturns(
							&model.DCGOpCall{
								DCGBaseCall: model.DCGBaseCall{
									DataHandle: fakeDataHandle,
								},
							},
							nil,
						)

						fakePkg := new(pkg.Fake)
						fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

						fakeDCGNodeRepo := new(fakeDCGNodeRepo)
						fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						fakePubSub := new(pubsub.Fake)
						eventChannel := make(chan model.Event)
						// close eventChannel to trigger immediate return
						close(eventChannel)
						fakePubSub.SubscribeReturns(eventChannel, nil)

						objectUnderTest := _opCaller{
							pkg:                fakePkg,
							pubSub:             fakePubSub,
							opCallInterpreter:  fakeOpCallInterpreter,
							outputsInterpreter: fakeOutputsInterpreter,
							dcgNodeRepo:        fakeDCGNodeRepo,
							caller:             new(fakeCaller),
						}

						/* act */
						objectUnderTest.Call(
							map[string]*model.Value{},
							providedOpId,
							new(data.FakeHandle),
							providedRootOpId,
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
