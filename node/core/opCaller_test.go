package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/managepackages"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/validate"
	"github.com/pkg/errors"
	"time"
)

var _ = Context("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(managepackages.Fake),
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

			expectedDCGNodeDescriptor := &dcgNodeDescriptor{
				Id:       providedOpId,
				PkgRef:   providedPkgRef,
				RootOpId: providedRootOpId,
				Op:       &dcgOpDescriptor{},
			}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			objectUnderTest := newOpCaller(
				new(managepackages.Fake),
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
			)

			/* assert */
			Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
		})
		It("should call managepackages.GetPackage w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedOpId := "dummyOpId"
			providedPkgRef := "dummyPkgRef"
			providedRootOpId := "dummyRootOpId"

			expectedPkgRef := providedPkgRef

			fakeManagePackages := new(managepackages.Fake)

			objectUnderTest := newOpCaller(
				fakeManagePackages,
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
			)

			/* assert */
			Expect(fakeManagePackages.GetPackageArgsForCall(0)).To(Equal(expectedPkgRef))
		})
		Context("managepackages.GetPackage errors", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedOpId := "dummyOpId"
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					OpEncounteredError: &model.OpEncounteredErrorEvent{
						Msg:      "dummyError",
						OpId:     providedOpId,
						PkgRef:   providedPkgRef,
						RootOpId: providedRootOpId,
					},
				}

				fakeManagePackages := new(managepackages.Fake)
				fakeManagePackages.GetPackageReturns(
					model.PackageView{},
					errors.New(expectedEvent.OpEncounteredError.Msg),
				)

				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := newOpCaller(
					fakeManagePackages,
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
		Context("managepackages.GetPackage doesn't error", func() {
			It("should call validate.Param w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{
					"dummyVar1Name": {String: "dummyVar1Data"},
					"dummyVar2Name": {File: "dummyVar2Data"},
					"dummyVar3Name": {Dir: "dummyVar3Data"},
					"dummyVar4Name": {Socket: "dummyVar4Data"},
				}
				providedOpId := "dummyOpId"
				providedPkgRef := "dummyPkgRef"
				providedRootOpId := "dummyRootOpId"

				opReturnedFromPkg := model.PackageView{
					Inputs: map[string]*model.Param{
						"dummyVar1Name": {
							String: &model.StringParam{},
						},
						"dummyVar2Name": {
							File: &model.FileParam{},
						},
						"dummyVar3Name": {
							Dir: &model.DirParam{},
						},
						"dummyVar4Name": {
							Socket: &model.SocketParam{},
						},
					},
				}
				fakeManagePackages := new(managepackages.Fake)
				fakeManagePackages.GetPackageReturns(opReturnedFromPkg, nil)

				expectedCalls := map[*model.Data]*model.Param{}
				for inputName, input := range opReturnedFromPkg.Inputs {
					expectedCalls[providedInboundScope[inputName]] = input
				}

				fakeValidate := new(validate.Fake)

				objectUnderTest := newOpCaller(
					fakeManagePackages,
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
				)

				/* assert */
				actualCalls := map[*model.Data]*model.Param{}
				for i := 0; i < fakeValidate.ParamCallCount(); i++ {
					actualVarData, actualParam := fakeValidate.ParamArgsForCall(i)
					actualCalls[actualVarData] = actualParam
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
					fakeManagePackages := new(managepackages.Fake)
					fakeManagePackages.GetPackageReturns(opReturnedFromPkg, nil)

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

-`, "dummyVar1Name", "", errorReturnedFromValidate)

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
						fakeManagePackages,
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
						new(managepackages.Fake),
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

					opReturnedFromPkg := model.PackageView{
						Run: &model.Scg{
							Parallel: []*model.Scg{
								{
									Container: &model.ScgContainerCall{},
								},
							},
						},
					}
					fakeManagePackages := new(managepackages.Fake)
					fakeManagePackages.GetPackageReturns(opReturnedFromPkg, nil)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					expectedNodeId := "dummyNodeId"
					fakeUniqueStringFactory.ConstructReturns(expectedNodeId)

					fakeCaller := new(fakeCaller)

					objectUnderTest := newOpCaller(
						fakeManagePackages,
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

					fakeDCGNodeRepo := new(fakeDCGNodeRepo)

					objectUnderTest := newOpCaller(
						new(managepackages.Fake),
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
							new(managepackages.Fake),
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

						fakeDCGNodeRepo := new(fakeDCGNodeRepo)
						fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						objectUnderTest := newOpCaller(
							new(managepackages.Fake),
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

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEncounteredError: &model.OpEncounteredErrorEvent{
									Msg:      "dummyError",
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
								new(managepackages.Fake),
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
								new(managepackages.Fake),
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
								new(managepackages.Fake),
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
