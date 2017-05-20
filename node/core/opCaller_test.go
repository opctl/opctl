package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
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
				"dummyRootFSPath",
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

				objectUnderTest := newOpCaller(
					new(pubsub.Fake),
					fakeDCGNodeRepo,
					new(fakeCaller),
					"dummyRootFSPath",
				)

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

			objectUnderTest := newOpCaller(
				new(pubsub.Fake),
				fakeDCGNodeRepo,
				new(fakeCaller),
				"dummyRootFSPath",
			)

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
		It("should call dcgOpCallFactory.Construct w/ expected args & return errors", func() {
			/* arrange */

			/* act */

			/* assert */
		})
		Context("dcgOpCallFactory.Construct errors", func() {
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
					OpErred: &model.OpErredEvent{
						Msg:      "dummyError",
						OpId:     providedOpId,
						PkgRef:   providedSCGOpCall.Pkg.Ref,
						RootOpId: providedRootOpId,
					},
				}

				fakeDCGOpCallFactory := new(fakeDCGOpCallFactory)
				fakeDCGOpCallFactory.ConstructReturns(
					nil,
					errors.New(expectedEvent.OpErred.Msg),
				)

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := _opCaller{
					dcgOpCallFactory: fakeDCGOpCallFactory,
					pubSub:           fakePubSub,
					dcgNodeRepo:      fakeDCGNodeRepo,
					caller:           new(fakeCaller),
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
		Context("dcgOpCallFactory.Construct doesn't error", func() {
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

				fakeDCGOpCallFactory := new(fakeDCGOpCallFactory)
				fakeDCGOpCallFactory.ConstructReturns(
					&model.DCGOpCall{
						DCGBaseCall: &model.DCGBaseCall{},
					},
					nil,
				)

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := _opCaller{
					dcgOpCallFactory: fakeDCGOpCallFactory,
					pubSub:           fakePubSub,
					dcgNodeRepo:      new(fakeDCGNodeRepo),
					caller:           new(fakeCaller),
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
					ChildCallScope: map[string]*model.Data{
						"dummyScopeName": {String: &dummyString},
					},
				}
				fakeDCGOpCallFactory := new(fakeDCGOpCallFactory)
				fakeDCGOpCallFactory.ConstructReturns(
					&dcgOpCall,
					nil,
				)

				fakeCaller := new(fakeCaller)

				objectUnderTest := _opCaller{
					dcgOpCallFactory: fakeDCGOpCallFactory,
					pubSub:           new(pubsub.Fake),
					dcgNodeRepo:      new(fakeDCGNodeRepo),
					caller:           fakeCaller,
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
				Expect(actualChildCallScope).To(Equal(dcgOpCall.ChildCallScope))
				Expect(actualChildSCG).To(Equal(dcgOpCall.ChildCallSCG))
				Expect(actualPkgRef).To(Equal(dcgOpCall.PkgRef))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			})
			It("should call dcgNodeRepo.GetIfExists w/ expected args", func() {
				/* arrange */
				providedRootOpId := "dummyRootOpId"

				fakeDCGOpCallFactory := new(fakeDCGOpCallFactory)
				fakeDCGOpCallFactory.ConstructReturns(
					&model.DCGOpCall{
						DCGBaseCall: &model.DCGBaseCall{},
					},
					nil,
				)

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)

				objectUnderTest := _opCaller{
					pubSub:           new(pubsub.Fake),
					dcgOpCallFactory: fakeDCGOpCallFactory,
					dcgNodeRepo:      fakeDCGNodeRepo,
					caller:           new(fakeCaller),
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

					fakeDCGOpCallFactory := new(fakeDCGOpCallFactory)
					fakeDCGOpCallFactory.ConstructReturns(
						&model.DCGOpCall{
							DCGBaseCall: &model.DCGBaseCall{},
						},
						nil,
					)

					fakePubSub := new(pubsub.Fake)

					objectUnderTest := _opCaller{
						pubSub:           fakePubSub,
						dcgOpCallFactory: fakeDCGOpCallFactory,
						dcgNodeRepo:      new(fakeDCGNodeRepo),
						caller:           new(fakeCaller),
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

					fakeDCGOpCallFactory := new(fakeDCGOpCallFactory)
					fakeDCGOpCallFactory.ConstructReturns(
						&model.DCGOpCall{
							DCGBaseCall: &model.DCGBaseCall{},
						},
						nil,
					)

					fakeDCGNodeRepo := new(fakeDCGNodeRepo)
					fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

					objectUnderTest := _opCaller{
						pubSub:           new(pubsub.Fake),
						dcgOpCallFactory: fakeDCGOpCallFactory,
						dcgNodeRepo:      fakeDCGNodeRepo,
						caller:           new(fakeCaller),
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

						callErr := errors.New("dummyError")

						expectedEvent := &model.Event{
							Timestamp: time.Now().UTC(),
							OpErred: &model.OpErredEvent{
								Msg:      callErr.Error(),
								OpId:     providedOpId,
								PkgRef:   providedSCGOpCall.Pkg.Ref,
								RootOpId: providedRootOpId,
							},
						}

						fakeDCGOpCallFactory := new(fakeDCGOpCallFactory)
						fakeDCGOpCallFactory.ConstructReturns(
							&model.DCGOpCall{
								DCGBaseCall: &model.DCGBaseCall{},
							},
							nil,
						)

						fakeDCGNodeRepo := new(fakeDCGNodeRepo)
						fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						fakePubSub := new(pubsub.Fake)

						fakeCaller := new(fakeCaller)
						fakeCaller.CallReturns(callErr)

						objectUnderTest := _opCaller{
							pubSub:           fakePubSub,
							dcgOpCallFactory: fakeDCGOpCallFactory,
							dcgNodeRepo:      fakeDCGNodeRepo,
							caller:           fakeCaller,
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

						fakeDCGOpCallFactory := new(fakeDCGOpCallFactory)
						fakeDCGOpCallFactory.ConstructReturns(
							&model.DCGOpCall{
								DCGBaseCall: &model.DCGBaseCall{},
							},
							nil,
						)

						fakeDCGNodeRepo := new(fakeDCGNodeRepo)
						fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						fakePubSub := new(pubsub.Fake)

						fakeCaller := new(fakeCaller)
						fakeCaller.CallReturns(
							errors.New("dummyError"),
						)

						objectUnderTest := _opCaller{
							pubSub:           fakePubSub,
							dcgOpCallFactory: fakeDCGOpCallFactory,
							dcgNodeRepo:      fakeDCGNodeRepo,
							caller:           fakeCaller,
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
							},
						}

						fakeDCGOpCallFactory := new(fakeDCGOpCallFactory)
						fakeDCGOpCallFactory.ConstructReturns(
							&model.DCGOpCall{
								DCGBaseCall: &model.DCGBaseCall{},
							},
							nil,
						)

						fakeDCGNodeRepo := new(fakeDCGNodeRepo)
						fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						fakePubSub := new(pubsub.Fake)

						objectUnderTest := _opCaller{
							pubSub:           fakePubSub,
							dcgOpCallFactory: fakeDCGOpCallFactory,
							dcgNodeRepo:      fakeDCGNodeRepo,
							caller:           new(fakeCaller),
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
