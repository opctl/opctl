package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/managepackages"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
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
				newDcgNodeRepo(),
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
			providedOutputs := make(chan *variable, 150)
			providedOpId := "dummyOpId"
			providedPkgRef := "dummyPkgRef"
			providedRootOpId := "dummyRootOpId"

			expectedDcgNodeDescriptor := &dcgNodeDescriptor{
				Id:       providedOpId,
				PkgRef:   providedPkgRef,
				RootOpId: providedRootOpId,
				Op:       &dcgOpDescriptor{},
			}

			fakeDcgNodeRepo := new(fakeDcgNodeRepo)

			fakeCaller := new(fakeCaller)
			// outputs chan must be closed for method under test to return
			fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
				close(outputs)
				return
			}

			objectUnderTest := newOpCaller(
				new(managepackages.Fake),
				new(pubsub.Fake),
				fakeDcgNodeRepo,
				fakeCaller,
				new(uniquestring.Fake),
				new(validate.Fake),
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedOutputs,
				providedOpId,
				providedPkgRef,
				providedRootOpId,
			)

			/* assert */
			Expect(fakeDcgNodeRepo.AddArgsForCall(0)).To(Equal(expectedDcgNodeDescriptor))
		})
		It("should call managepackages.GetPackage w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedOutputs := make(chan *variable, 150)
			providedOpId := "dummyOpId"
			providedPkgRef := "dummyPkgRef"
			providedRootOpId := "dummyRootOpId"

			expectedPkgRef := providedPkgRef

			fakeManagePackages := new(managepackages.Fake)

			fakeCaller := new(fakeCaller)
			// outputs chan must be closed for method under test to return
			fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
				close(outputs)
				return
			}

			objectUnderTest := newOpCaller(
				fakeManagePackages,
				new(pubsub.Fake),
				new(fakeDcgNodeRepo),
				fakeCaller,
				new(uniquestring.Fake),
				new(validate.Fake),
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedOutputs,
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
				providedOutputs := make(chan *variable, 150)
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

				fakeDcgNodeRepo := new(fakeDcgNodeRepo)
				fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := newOpCaller(
					fakeManagePackages,
					fakePubSub,
					fakeDcgNodeRepo,
					new(fakeCaller),
					new(uniquestring.Fake),
					new(validate.Fake),
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedOutputs,
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
				providedOutputs := make(chan *variable, 150)
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

				fakeCaller := new(fakeCaller)
				// outputs chan must be closed for method under test to return
				fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
					close(outputs)
					return
				}

				objectUnderTest := newOpCaller(
					fakeManagePackages,
					new(pubsub.Fake),
					new(fakeDcgNodeRepo),
					fakeCaller,
					new(uniquestring.Fake),
					fakeValidate,
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedOutputs,
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
					providedOutputs := make(chan *variable, 150)
					providedOpId := "dummyOpId"
					providedPkgRef := "dummyPkgRef"
					providedRootOpId := "dummyRootOpId"

					fakeDcgNodeRepo := new(fakeDcgNodeRepo)
					fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

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
						fakeDcgNodeRepo,
						new(fakeCaller),
						new(uniquestring.Fake),
						fakeValidate,
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOutputs,
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
					providedOutputs := make(chan *variable, 150)
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

					fakeCaller := new(fakeCaller)
					// outputs chan must be closed for method under test to return
					fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
						close(outputs)
						return
					}

					objectUnderTest := newOpCaller(
						new(managepackages.Fake),
						fakePubSub,
						new(fakeDcgNodeRepo),
						fakeCaller,
						new(uniquestring.Fake),
						new(validate.Fake),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOutputs,
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
					providedOutputs := make(chan *variable, 150)
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

					fakeCaller := new(fakeCaller)
					// outputs chan must be closed for method under test to return
					fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
						close(outputs)
						return
					}

					fakeUniqueStringFactory := new(uniquestring.Fake)
					expectedNodeId := "dummyNodeId"
					fakeUniqueStringFactory.ConstructReturns(expectedNodeId)

					objectUnderTest := newOpCaller(
						fakeManagePackages,
						new(pubsub.Fake),
						new(fakeDcgNodeRepo),
						fakeCaller,
						fakeUniqueStringFactory,
						new(validate.Fake),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOutputs,
						providedOpId,
						providedPkgRef,
						providedRootOpId,
					)

					/* assert */
					actualNodeId,
						actualInboundScope,
						_,
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
					providedOutputs := make(chan *variable, 150)
					providedOpId := "dummyOpId"
					providedPkgRef := "dummyPkgRef"
					providedRootOpId := "dummyRootOpId"

					fakeDcgNodeRepo := new(fakeDcgNodeRepo)

					fakeCaller := new(fakeCaller)
					// outputs chan must be closed for method under test to return
					fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
						close(outputs)
						return
					}

					objectUnderTest := newOpCaller(
						new(managepackages.Fake),
						new(pubsub.Fake),
						fakeDcgNodeRepo,
						fakeCaller,
						new(uniquestring.Fake),
						new(validate.Fake),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOutputs,
						providedOpId,
						providedPkgRef,
						providedRootOpId,
					)

					/* assert */
					Expect(fakeDcgNodeRepo.GetIfExistsArgsForCall(0)).To(Equal(providedRootOpId))
				})
				Context("dcgNodeRepo.GetIfExists returns nil", func() {
					It("should call pubSub.Publish w/ expected args", func() {
						/* arrange */
						providedInboundScope := map[string]*model.Data{}
						providedOutputs := make(chan *variable, 150)
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

						fakeCaller := new(fakeCaller)
						// outputs chan must be closed for method under test to return
						fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
							close(outputs)
							return
						}

						objectUnderTest := newOpCaller(
							new(managepackages.Fake),
							fakePubSub,
							new(fakeDcgNodeRepo),
							fakeCaller,
							new(uniquestring.Fake),
							new(validate.Fake),
						)

						/* act */
						objectUnderTest.Call(
							providedInboundScope,
							providedOutputs,
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
						providedOutputs := make(chan *variable, 150)
						providedOpId := "dummyOpId"
						providedPkgRef := "dummyPkgRef"
						providedRootOpId := "dummyRootOpId"

						fakeDcgNodeRepo := new(fakeDcgNodeRepo)
						fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						fakeCaller := new(fakeCaller)
						// outputs chan must be closed for method under test to return
						fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
							close(outputs)
							return
						}

						objectUnderTest := newOpCaller(
							new(managepackages.Fake),
							new(pubsub.Fake),
							fakeDcgNodeRepo,
							fakeCaller,
							new(uniquestring.Fake),
							new(validate.Fake),
						)

						/* act */
						objectUnderTest.Call(
							providedInboundScope,
							providedOutputs,
							providedOpId,
							providedPkgRef,
							providedRootOpId,
						)

						/* assert */
						Expect(fakeDcgNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedOpId))
					})
					Context("caller.Call errored", func() {
						It("should call pubSub.Publish w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOutputs := make(chan *variable, 150)
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

							fakeDcgNodeRepo := new(fakeDcgNodeRepo)
							fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							fakeCaller := new(fakeCaller)
							fakeCaller.CallReturns(
								errors.New(expectedEvent.OpEncounteredError.Msg),
							)

							objectUnderTest := newOpCaller(
								new(managepackages.Fake),
								fakePubSub,
								fakeDcgNodeRepo,
								fakeCaller,
								new(uniquestring.Fake),
								new(validate.Fake),
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOutputs,
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
							providedOutputs := make(chan *variable, 150)
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

							fakeDcgNodeRepo := new(fakeDcgNodeRepo)
							fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							fakeCaller := new(fakeCaller)
							fakeCaller.CallReturns(
								errors.New("dummyError"),
							)

							objectUnderTest := newOpCaller(
								new(managepackages.Fake),
								fakePubSub,
								fakeDcgNodeRepo,
								fakeCaller,
								new(uniquestring.Fake),
								new(validate.Fake),
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOutputs,
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
							providedOutputs := make(chan *variable, 150)
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

							fakeDcgNodeRepo := new(fakeDcgNodeRepo)
							fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							fakeCaller := new(fakeCaller)
							// outputs chan must be closed for method under test to return
							fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
								close(outputs)
								return
							}

							objectUnderTest := newOpCaller(
								new(managepackages.Fake),
								fakePubSub,
								fakeDcgNodeRepo,
								fakeCaller,
								new(uniquestring.Fake),
								new(validate.Fake),
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOutputs,
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
