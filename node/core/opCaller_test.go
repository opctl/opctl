package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall"
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
		Context("deprecated pkg format", func() {
			It("should call dcgNodeRepo.add w/ expected args", func() {
				/* arrange */
				providedOpId := "dummyOpId"
				providedRootOpId := "dummyRootOpId"
				providedSCGOpCall := &model.SCGOpCall{Ref: "dummyPkgRef"}

				expectedDCGNodeDescriptor := &dcgNodeDescriptor{
					Id:       providedOpId,
					PkgRef:   providedSCGOpCall.Ref,
					RootOpId: providedRootOpId,
					Op:       &dcgOpDescriptor{},
				}

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)

				fakeOpCall := new(opcall.Fake)
				// error to trigger immediate return
				fakeOpCall.InterpretReturns(nil, errors.New("dummyError"))

				objectUnderTest := _opCaller{
					opCall:      fakeOpCall,
					pubSub:      new(pubsub.Fake),
					dcgNodeRepo: fakeDCGNodeRepo,
					caller:      new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
					map[string]*model.Data{},
					providedOpId,
					"dummyPkgBasePath",
					providedRootOpId,
					providedSCGOpCall,
				)

				/* assert */
				Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
			})
		})
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

			fakeOpCall := new(opcall.Fake)
			// error to trigger immediate return
			fakeOpCall.InterpretReturns(nil, errors.New("dummyError"))

			objectUnderTest := _opCaller{
				opCall:      fakeOpCall,
				pubSub:      new(pubsub.Fake),
				dcgNodeRepo: fakeDCGNodeRepo,
				caller:      new(fakeCaller),
			}

			/* act */
			objectUnderTest.Call(
				map[string]*model.Data{},
				providedOpId,
				"dummyPkgBasePath",
				providedRootOpId,
				providedSCGOpCall,
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
		It("should call opCall.Construct w/ expected args & return errors", func() {
			/* arrange */
			providedScope := map[string]*model.Data{}
			providedOpId := "dummyOpId"
			providedRootOpId := "dummyRootOpId"
			providedPkgBasePath := "dummyPkgBasePath"
			providedSCGOpCall := &model.SCGOpCall{
				Pkg: &model.SCGOpCallPkg{
					Ref: "dummyPkgRef",
				},
			}

			fakeOpCall := new(opcall.Fake)

			// error to trigger immediate return
			expectedErr := errors.New("dummyError")
			fakeOpCall.InterpretReturns(
				&model.DCGOpCall{
					DCGBaseCall: &model.DCGBaseCall{},
				},
				expectedErr,
			)

			objectUnderTest := _opCaller{
				opCall:      fakeOpCall,
				pubSub:      new(pubsub.Fake),
				dcgNodeRepo: new(fakeDCGNodeRepo),
				caller:      new(fakeCaller),
			}

			/* act */
			actualErr := objectUnderTest.Call(
				providedScope,
				providedOpId,
				providedPkgBasePath,
				providedRootOpId,
				providedSCGOpCall,
			)

			/* assert */
			actualScope,
				actualSCGOpCall,
				actualOpId,
				actualPkgBasePath,
				actualRootOpId := fakeOpCall.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualSCGOpCall).To(Equal(providedSCGOpCall))
			Expect(actualOpId).To(Equal(providedOpId))
			Expect(actualPkgBasePath).To(Equal(providedPkgBasePath))
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

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					OpErred: &model.OpErredEvent{
						Msg:      "dummyError",
						OpId:     providedOpId,
						PkgRef:   providedSCGOpCall.Pkg.Ref,
						RootOpId: providedRootOpId,
					},
				}

				fakeOpCall := new(opcall.Fake)
				fakeOpCall.InterpretReturns(
					nil,
					errors.New(expectedEvent.OpErred.Msg),
				)

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := _opCaller{
					opCall:      fakeOpCall,
					pubSub:      fakePubSub,
					dcgNodeRepo: fakeDCGNodeRepo,
					caller:      new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
					map[string]*model.Data{},
					providedOpId,
					"dummyPkgBasePath",
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

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					OpStarted: &model.OpStartedEvent{
						OpId:     providedOpId,
						PkgRef:   providedSCGOpCall.Pkg.Ref,
						RootOpId: providedRootOpId,
					},
				}

				fakeOpCall := new(opcall.Fake)
				fakeOpCall.InterpretReturns(
					&model.DCGOpCall{
						DCGBaseCall: &model.DCGBaseCall{},
					},
					nil,
				)

				fakePubSub := new(pubsub.Fake)
				fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
					// close eventChannel to trigger immediate return
					close(eventChannel)
				}

				objectUnderTest := _opCaller{
					opCall:      fakeOpCall,
					pubSub:      fakePubSub,
					dcgNodeRepo: new(fakeDCGNodeRepo),
					caller:      new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
					map[string]*model.Data{},
					providedOpId,
					"dummyPkgBasePath",
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
					DCGBaseCall: &model.DCGBaseCall{
						PkgRef: "dummyPkgRef",
					},
					ChildCallId: "dummyChildCallId",
					ChildCallSCG: &model.SCG{
						Parallel: []*model.SCG{
							{
								Container: &model.SCGContainerCall{},
							},
						},
					},
					Inputs: map[string]*model.Data{
						"dummyScopeName": {String: &dummyString},
					},
				}
				fakeOpCall := new(opcall.Fake)
				fakeOpCall.InterpretReturns(
					&dcgOpCall,
					nil,
				)

				fakePubSub := new(pubsub.Fake)
				fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
					// close eventChannel to trigger immediate return
					close(eventChannel)
				}

				fakeCaller := new(fakeCaller)

				objectUnderTest := _opCaller{
					opCall:      fakeOpCall,
					pubSub:      fakePubSub,
					dcgNodeRepo: new(fakeDCGNodeRepo),
					caller:      fakeCaller,
				}

				/* act */
				objectUnderTest.Call(
					map[string]*model.Data{},
					"dummyOpId",
					"dummyPkgBasePath",
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
				Expect(actualPkgRef).To(Equal(dcgOpCall.PkgRef))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			})
			It("should call dcgNodeRepo.GetIfExists w/ expected args", func() {
				/* arrange */
				providedRootOpId := "dummyRootOpId"

				fakeOpCall := new(opcall.Fake)
				fakeOpCall.InterpretReturns(
					&model.DCGOpCall{
						DCGBaseCall: &model.DCGBaseCall{},
					},
					nil,
				)

				fakePubSub := new(pubsub.Fake)
				fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
					// close eventChannel to trigger immediate return
					close(eventChannel)
				}

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)

				objectUnderTest := _opCaller{
					pubSub:      fakePubSub,
					opCall:      fakeOpCall,
					dcgNodeRepo: fakeDCGNodeRepo,
					caller:      new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
					map[string]*model.Data{},
					"dummyOpId",
					"dummyPkgBasePath",
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

					expectedEvent := &model.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &model.OpEndedEvent{
							OpId:     providedOpId,
							Outcome:  model.OpOutcomeKilled,
							RootOpId: providedRootOpId,
							PkgRef:   providedSCGOpCall.Pkg.Ref,
						},
					}

					fakeOpCall := new(opcall.Fake)
					fakeOpCall.InterpretReturns(
						&model.DCGOpCall{
							DCGBaseCall: &model.DCGBaseCall{},
						},
						nil,
					)

					fakePubSub := new(pubsub.Fake)
					fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
						// close eventChannel to trigger immediate return
						close(eventChannel)
					}

					objectUnderTest := _opCaller{
						pubSub:      fakePubSub,
						opCall:      fakeOpCall,
						dcgNodeRepo: new(fakeDCGNodeRepo),
						caller:      new(fakeCaller),
					}

					/* act */
					objectUnderTest.Call(
						map[string]*model.Data{},
						providedOpId,
						"dummyPkgBasePath",
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

					fakeOpCall := new(opcall.Fake)
					fakeOpCall.InterpretReturns(
						&model.DCGOpCall{
							DCGBaseCall: &model.DCGBaseCall{},
						},
						nil,
					)

					fakeDCGNodeRepo := new(fakeDCGNodeRepo)
					fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

					fakePubSub := new(pubsub.Fake)
					fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
						// close eventChannel to trigger immediate return
						close(eventChannel)
					}

					objectUnderTest := _opCaller{
						pubSub:      fakePubSub,
						opCall:      fakeOpCall,
						dcgNodeRepo: fakeDCGNodeRepo,
						caller:      new(fakeCaller),
					}

					/* act */
					objectUnderTest.Call(
						map[string]*model.Data{},
						providedOpId,
						"dummyPkgBasePath",
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

						expectedEvent := &model.Event{
							Timestamp: time.Now().UTC(),
							OpEnded: &model.OpEndedEvent{
								OpId:     providedOpId,
								PkgRef:   providedSCGOpCall.Pkg.Ref,
								Outcome:  model.OpOutcomeFailed,
								RootOpId: providedRootOpId,
							},
						}

						fakeOpCall := new(opcall.Fake)
						fakeOpCall.InterpretReturns(
							&model.DCGOpCall{
								DCGBaseCall: &model.DCGBaseCall{},
							},
							nil,
						)

						fakeDCGNodeRepo := new(fakeDCGNodeRepo)
						fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						fakePubSub := new(pubsub.Fake)
						fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
							// close eventChannel to trigger immediate return
							close(eventChannel)
						}

						fakeCaller := new(fakeCaller)
						fakeCaller.CallReturns(
							errors.New("dummyError"),
						)

						objectUnderTest := _opCaller{
							pubSub:      fakePubSub,
							opCall:      fakeOpCall,
							dcgNodeRepo: fakeDCGNodeRepo,
							caller:      fakeCaller,
						}

						/* act */
						objectUnderTest.Call(
							map[string]*model.Data{},
							providedOpId,
							"dummyPkgBasePath",
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

						expectedEvent := &model.Event{
							Timestamp: time.Now().UTC(),
							OpEnded: &model.OpEndedEvent{
								OpId:     providedOpId,
								PkgRef:   providedSCGOpCall.Pkg.Ref,
								Outcome:  model.OpOutcomeSucceeded,
								RootOpId: providedRootOpId,
								Outputs:  map[string]*model.Data{},
							},
						}

						fakeOpCall := new(opcall.Fake)
						fakeOpCall.InterpretReturns(
							&model.DCGOpCall{
								DCGBaseCall: &model.DCGBaseCall{},
							},
							nil,
						)

						fakeDCGNodeRepo := new(fakeDCGNodeRepo)
						fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						fakePubSub := new(pubsub.Fake)
						fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
							// close eventChannel to trigger immediate return
							close(eventChannel)
						}

						objectUnderTest := _opCaller{
							pubSub:      fakePubSub,
							opCall:      fakeOpCall,
							dcgNodeRepo: fakeDCGNodeRepo,
							caller:      new(fakeCaller),
						}

						/* act */
						objectUnderTest.Call(
							map[string]*model.Data{},
							providedOpId,
							"dummyPkgBasePath",
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
