package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/validate"
	"github.com/pkg/errors"
	"time"
)

var _ = Context("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(pkg.Fake),
				new(pubsub.Fake),
				newDCGNodeRepo(),
				new(fakeCaller),
				new(uniquestring.Fake),
				new(validate.Fake),
			)).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call dcgNodeRepo.add w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedOpId := "dummyOpId"
			providedPkgRef := "dummyPkgRef"
			providedRootOpId := "dummyRootOpId"
			providedScgOpCall := &model.ScgOpCall{}

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:       providedOpId,
				PkgRef:   providedPkgRef,
				RootOpId: providedRootOpId,
				Op:       &dcgOpDescriptor{},
			}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			objectUnderTest := newOpCaller(
				new(pkg.Fake),
				new(pubsub.Fake),
				fakeDCGNodeRepo,
				new(fakeCaller),
				new(uniquestring.Fake),
				new(validate.Fake),
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedOpId,
				providedPkgRef,
				providedRootOpId,
				providedScgOpCall,
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
		It("should call pkg.Get w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedOpId := "dummyOpId"
			providedPkgRef := "dummyPkgRef"
			providedRootOpId := "dummyRootOpId"
			providedScgOpCall := &model.ScgOpCall{}

			expectedPkgRef := providedPkgRef

			fakePkg := new(pkg.Fake)

			objectUnderTest := newOpCaller(
				fakePkg,
				new(pubsub.Fake),
				new(fakeDCGNodeRepo),
				new(fakeCaller),
				new(uniquestring.Fake),
				new(validate.Fake),
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedOpId,
				providedPkgRef,
				providedRootOpId,
				providedScgOpCall,
			)

			/* assert */
			Expect(fakePkg.GetArgsForCall(0)).To(Equal(expectedPkgRef))
		})
		Context("pkg.Get errors", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedOpId := "dummyOpId"
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"
				providedScgOpCall := &model.ScgOpCall{}

				errorsReturnedFromValidate := []error{
					errors.New("dummyError1"),
					errors.New("dummyError2"),
				}

				expectedErrorMsg :=
					fmt.Errorf(`

-
  Package manifest %v invalid.
  Error(s):
    - %v
    - %v
-`,
						providedPkgRef,
						errorsReturnedFromValidate[0],
						errorsReturnedFromValidate[1],
					)

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					OpEncounteredError: &model.OpEncounteredErrorEvent{
						Msg:      expectedErrorMsg.Error(),
						OpId:     providedOpId,
						PkgRef:   providedPkgRef,
						RootOpId: providedRootOpId,
					},
				}

				fakePkg := new(pkg.Fake)
				fakePkg.ValidateReturns(errorsReturnedFromValidate)

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := newOpCaller(
					fakePkg,
					fakePubSub,
					fakeDCGNodeRepo,
					new(fakeCaller),
					new(uniquestring.Fake),
					new(validate.Fake),
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedOpId,
					providedPkgRef,
					providedRootOpId,
					providedScgOpCall,
				)

				/* assert */
				actualEvent := fakePubSub.PublishArgsForCall(0)

				// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
				Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
				// set temporal fields to expected vals since they're already asserted
				actualEvent.Timestamp = expectedEvent.Timestamp

				Expect(actualEvent).To(Equal(expectedEvent))
			})
		})
		Context("pkg.Get errors", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedOpId := "dummyOpId"
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"
				providedScgOpCall := &model.ScgOpCall{}

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					OpEncounteredError: &model.OpEncounteredErrorEvent{
						Msg:      "dummyError",
						OpId:     providedOpId,
						PkgRef:   providedPkgRef,
						RootOpId: providedRootOpId,
					},
				}

				fakePkg := new(pkg.Fake)
				fakePkg.GetReturns(
					model.PackageView{},
					errors.New(expectedEvent.OpEncounteredError.Msg),
				)

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := newOpCaller(
					fakePkg,
					fakePubSub,
					fakeDCGNodeRepo,
					new(fakeCaller),
					new(uniquestring.Fake),
					new(validate.Fake),
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedOpId,
					providedPkgRef,
					providedRootOpId,
					providedScgOpCall,
				)

				/* assert */
				actualEvent := fakePubSub.PublishArgsForCall(0)

				// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
				Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
				// set temporal fields to expected vals since they're already asserted
				actualEvent.Timestamp = expectedEvent.Timestamp

				Expect(actualEvent).To(Equal(expectedEvent))
			})
		})
		Context("pkg.Get doesn't error", func() {
			It("should call validate.Param w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{
					"name1": {String: "val1"},
					"name2": {File: "val2"},
					"name3": {Dir: "val3"},
					"name4": {Socket: "val4"},
					"name5": {Number: 5},
				}
				providedOpId := "dummyOpId"
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"
				providedScgOpCall := &model.ScgOpCall{
					Inputs: map[string]string{
						"name1": "",
						"name2": "",
						"name3": "",
						"name4": "",
						"name5": "",
						"name6": "",
						"name7": "",
					},
				}

				returnedPkg := model.PackageView{
					Inputs: map[string]*model.Param{
						"name1": {String: &model.StringParam{}},
						"name2": {File: &model.FileParam{}},
						"name3": {Dir: &model.DirParam{}},
						"name4": {Socket: &model.SocketParam{}},
						"name5": {Number: &model.NumberParam{}},
						"name6": {Number: &model.NumberParam{Default: 6}},
						"name7": {String: &model.StringParam{Default: "seven"}},
					},
				}
				fakePkg := new(pkg.Fake)
				fakePkg.GetReturns(returnedPkg, nil)

				expectedCalls := map[model.Data]*model.Param{
					// from scope
					*providedInboundScope["name1"]: returnedPkg.Inputs["name1"],
					*providedInboundScope["name2"]: returnedPkg.Inputs["name2"],
					*providedInboundScope["name3"]: returnedPkg.Inputs["name3"],
					*providedInboundScope["name4"]: returnedPkg.Inputs["name4"],
					*providedInboundScope["name5"]: returnedPkg.Inputs["name5"],
					// from defaults
					model.Data{
						Number: returnedPkg.Inputs["name6"].Number.Default,
					}: returnedPkg.Inputs["name6"],
					model.Data{
						String: returnedPkg.Inputs["name7"].String.Default,
					}: returnedPkg.Inputs["name7"],
				}

				fakeValidate := new(validate.Fake)

				objectUnderTest := newOpCaller(
					fakePkg,
					new(pubsub.Fake),
					new(fakeDCGNodeRepo),
					new(fakeCaller),
					new(uniquestring.Fake),
					fakeValidate,
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedOpId,
					providedPkgRef,
					providedRootOpId,
					providedScgOpCall,
				)

				/* assert */
				actualCalls := map[model.Data]*model.Param{}
				for i := 0; i < fakeValidate.ParamCallCount(); i++ {
					actualVarData, actualParam := fakeValidate.ParamArgsForCall(i)
					actualCalls[*actualVarData] = actualParam
				}
				Expect(actualCalls).To(Equal(expectedCalls))
			})
			Context("validate.Param errors", func() {
				It("should call pubSub.Publish w/ expected args", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedOpId := "dummyOpId"
					providedPkgRef := "dummyPkgRef"
					providedRootOpId := "dummyRootOpId"
					providedScgOpCall := &model.ScgOpCall{}

					fakeDCGNodeRepo := new(fakeDCGNodeRepo)
					fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

					opReturnedFromPkg := model.PackageView{
						Inputs: map[string]*model.Param{
							"dummyVar1Name": {
								String: &model.StringParam{
									IsSecret: true,
								},
							},
						},
					}
					fakePkg := new(pkg.Fake)
					fakePkg.GetReturns(opReturnedFromPkg, nil)

					fakeValidate := new(validate.Fake)

					errorReturnedFromValidate := "validationError0"
					fakeValidate.ParamReturns([]error{errors.New(errorReturnedFromValidate)})

					expectedMsg := fmt.Sprintf(`
-
  validation of the following input failed:

  Name: %v
  Value: %v
  Error(s):
    - %v

-`, "dummyVar1Name", "************", errorReturnedFromValidate)

					fakePubSub := new(pubsub.Fake)
					expectedEvent := &model.Event{
						Timestamp: time.Now().UTC(),
						OpEncounteredError: &model.OpEncounteredErrorEvent{
							Msg:      expectedMsg,
							OpId:     providedOpId,
							PkgRef:   providedPkgRef,
							RootOpId: providedRootOpId,
						},
					}

					objectUnderTest := newOpCaller(
						fakePkg,
						fakePubSub,
						fakeDCGNodeRepo,
						new(fakeCaller),
						new(uniquestring.Fake),
						fakeValidate,
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						providedPkgRef,
						providedRootOpId,
						providedScgOpCall,
					)

					/* assert */
					actualEvent := fakePubSub.PublishArgsForCall(0)

					// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
					Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
					// set temporal fields to expected vals since they're already asserted
					actualEvent.Timestamp = expectedEvent.Timestamp

					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
			Context("validate.Param doesn't error", func() {
				It("should call pubSub.Publish w/ expected args", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedOpId := "dummyOpId"
					providedPkgRef := "dummyPkgRef"
					providedRootOpId := "dummyRootOpId"
					providedScgOpCall := &model.ScgOpCall{}

					expectedEvent := &model.Event{
						Timestamp: time.Now().UTC(),
						OpStarted: &model.OpStartedEvent{
							OpId:     providedOpId,
							PkgRef:   providedPkgRef,
							RootOpId: providedRootOpId,
						},
					}

					fakePubSub := new(pubsub.Fake)

					objectUnderTest := newOpCaller(
						new(pkg.Fake),
						fakePubSub,
						new(fakeDCGNodeRepo),
						new(fakeCaller),
						new(uniquestring.Fake),
						new(validate.Fake),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						providedPkgRef,
						providedRootOpId,
						providedScgOpCall,
					)

					/* assert */
					actualEvent := fakePubSub.PublishArgsForCall(0)

					// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
					Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
					// set temporal fields to expected vals since they're already asserted
					actualEvent.Timestamp = expectedEvent.Timestamp

					Expect(actualEvent).To(Equal(expectedEvent))
				})
				It("should call caller.Call w/ expected args", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedOpId := "dummyOpId"
					providedPkgRef := "dummyPkgRef"
					providedRootOpId := "dummyRootOpId"
					providedScgOpCall := &model.ScgOpCall{}

					opReturnedFromPkg := model.PackageView{
						Run: &model.Scg{
							Parallel: []*model.Scg{
								{
									Container: &model.ScgContainerCall{},
								},
							},
						},
					}
					fakePkg := new(pkg.Fake)
					fakePkg.GetReturns(opReturnedFromPkg, nil)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					expectedNodeId := "dummyNodeId"
					fakeUniqueStringFactory.ConstructReturns(expectedNodeId)

					fakeCaller := new(fakeCaller)

					objectUnderTest := newOpCaller(
						fakePkg,
						new(pubsub.Fake),
						new(fakeDCGNodeRepo),
						fakeCaller,
						fakeUniqueStringFactory,
						new(validate.Fake),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						providedPkgRef,
						providedRootOpId,
						providedScgOpCall,
					)

					/* assert */
					actualNodeId,
						actualInboundScope,
						actualScg,
						actualPkgRef,
						actualRootOpId := fakeCaller.CallArgsForCall(0)

					Expect(actualNodeId).To(Equal(expectedNodeId))
					Expect(actualInboundScope).To(Equal(providedInboundScope))
					Expect(actualScg).To(Equal(opReturnedFromPkg.Run))
					Expect(actualPkgRef).To(Equal(providedPkgRef))
					Expect(actualRootOpId).To(Equal(providedRootOpId))
				})
				It("should call dcgNodeRepo.GetIfExists w/ expected args", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedOpId := "dummyOpId"
					providedPkgRef := "dummyPkgRef"
					providedRootOpId := "dummyRootOpId"
					providedScgOpCall := &model.ScgOpCall{}

					fakeDCGNodeRepo := new(fakeDCGNodeRepo)

					objectUnderTest := newOpCaller(
						new(pkg.Fake),
						new(pubsub.Fake),
						fakeDCGNodeRepo,
						new(fakeCaller),
						new(uniquestring.Fake),
						new(validate.Fake),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						providedPkgRef,
						providedRootOpId,
						providedScgOpCall,
					)

					/* assert */
					Expect(fakeDCGNodeRepo.GetIfExistsArgsForCall(0)).To(Equal(providedRootOpId))
				})
				Context("dcgNodeRepo.GetIfExists returns nil", func() {
					It("should call pubSub.Publish w/ expected args", func() {
						/* arrange */
						providedInboundScope := map[string]*model.Data{}
						providedOpId := "dummyOpId"
						providedPkgRef := "dummyPkgRef"
						providedRootOpId := "dummyRootOpId"
						providedScgOpCall := &model.ScgOpCall{}

						expectedEvent := &model.Event{
							Timestamp: time.Now().UTC(),
							OpEnded: &model.OpEndedEvent{
								OpId:     providedOpId,
								Outcome:  model.OpOutcomeKilled,
								RootOpId: providedRootOpId,
								PkgRef:   providedPkgRef,
							},
						}

						fakePubSub := new(pubsub.Fake)

						objectUnderTest := newOpCaller(
							new(pkg.Fake),
							fakePubSub,
							new(fakeDCGNodeRepo),
							new(fakeCaller),
							new(uniquestring.Fake),
							new(validate.Fake),
						)

						/* act */
						objectUnderTest.Call(
							providedInboundScope,
							providedOpId,
							providedPkgRef,
							providedRootOpId,
							providedScgOpCall,
						)

						/* assert */
						actualEvent := fakePubSub.PublishArgsForCall(1)

						// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
						Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
						// set temporal fields to expected vals since they're already asserted
						actualEvent.Timestamp = expectedEvent.Timestamp

						Expect(actualEvent).To(Equal(expectedEvent))
					})
				})
				Context("dcgNodeRepo.GetIfExists doesn't return nil", func() {
					It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
						/* arrange */
						providedInboundScope := map[string]*model.Data{}
						providedOpId := "dummyOpId"
						providedPkgRef := "dummyPkgRef"
						providedRootOpId := "dummyRootOpId"
						providedScgOpCall := &model.ScgOpCall{}

						fakeDCGNodeRepo := new(fakeDCGNodeRepo)
						fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						objectUnderTest := newOpCaller(
							new(pkg.Fake),
							new(pubsub.Fake),
							fakeDCGNodeRepo,
							new(fakeCaller),
							new(uniquestring.Fake),
							new(validate.Fake),
						)

						/* act */
						objectUnderTest.Call(
							providedInboundScope,
							providedOpId,
							providedPkgRef,
							providedRootOpId,
							providedScgOpCall,
						)

						/* assert */
						Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedOpId))
					})
					Context("caller.Call errored", func() {
						It("should call pubSub.Publish w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOpId := "dummyOpId"
							providedPkgRef := "dummyPkgRef"
							providedRootOpId := "dummyRootOpId"
							providedScgOpCall := &model.ScgOpCall{}

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEncounteredError: &model.OpEncounteredErrorEvent{
									Msg:      "Error encountered during call",
									OpId:     providedOpId,
									PkgRef:   providedPkgRef,
									RootOpId: providedRootOpId,
								},
							}

							fakeDCGNodeRepo := new(fakeDCGNodeRepo)
							fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							fakeCaller := new(fakeCaller)
							fakeCaller.CallReturns(
								errors.New(expectedEvent.OpEncounteredError.Msg),
							)

							objectUnderTest := newOpCaller(
								new(pkg.Fake),
								fakePubSub,
								fakeDCGNodeRepo,
								fakeCaller,
								new(uniquestring.Fake),
								new(validate.Fake),
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								providedPkgRef,
								providedRootOpId,
								providedScgOpCall,
							)

							/* assert */
							actualEvent := fakePubSub.PublishArgsForCall(1)

							// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
							Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
							// set temporal fields to expected vals since they're already asserted
							actualEvent.Timestamp = expectedEvent.Timestamp

							Expect(actualEvent).To(Equal(expectedEvent))
						})
						It("should call pubSub.Publish w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOpId := "dummyOpId"
							providedPkgRef := "dummyPkgRef"
							providedRootOpId := "dummyRootOpId"
							providedScgOpCall := &model.ScgOpCall{}

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEnded: &model.OpEndedEvent{
									OpId:     providedOpId,
									PkgRef:   providedPkgRef,
									Outcome:  model.OpOutcomeFailed,
									RootOpId: providedRootOpId,
								},
							}

							fakeDCGNodeRepo := new(fakeDCGNodeRepo)
							fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							fakeCaller := new(fakeCaller)
							fakeCaller.CallReturns(
								errors.New("dummyError"),
							)

							objectUnderTest := newOpCaller(
								new(pkg.Fake),
								fakePubSub,
								fakeDCGNodeRepo,
								fakeCaller,
								new(uniquestring.Fake),
								new(validate.Fake),
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								providedPkgRef,
								providedRootOpId,
								providedScgOpCall,
							)

							/* assert */
							actualEvent := fakePubSub.PublishArgsForCall(2)

							// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
							Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
							// set temporal fields to expected vals since they're already asserted
							actualEvent.Timestamp = expectedEvent.Timestamp

							Expect(actualEvent).To(Equal(expectedEvent))
						})
					})
					Context("caller.Call didn't error", func() {
						It("should call pubSub.Publish w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOpId := "dummyOpId"
							providedPkgRef := "dummyPkgRef"
							providedRootOpId := "dummyRootOpId"
							providedScgOpCall := &model.ScgOpCall{}

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEnded: &model.OpEndedEvent{
									OpId:     providedOpId,
									PkgRef:   providedPkgRef,
									Outcome:  model.OpOutcomeSucceeded,
									RootOpId: providedRootOpId,
								},
							}

							fakeDCGNodeRepo := new(fakeDCGNodeRepo)
							fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							objectUnderTest := newOpCaller(
								new(pkg.Fake),
								fakePubSub,
								fakeDCGNodeRepo,
								new(fakeCaller),
								new(uniquestring.Fake),
								new(validate.Fake),
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								providedPkgRef,
								providedRootOpId,
								providedScgOpCall,
							)

							/* assert */
							actualEvent := fakePubSub.PublishArgsForCall(1)

							// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
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
})
