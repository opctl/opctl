package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/pkg"
	"github.com/opspec-io/sdk-golang/pkg/validate"
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
			providedOpId := "dummyOpId"
			providedOpPkgRef := "dummyOpPkgRef"
			providedRootOpId := "dummyRootOpId"

			expectedDcgNodeDescriptor := &dcgNodeDescriptor{
				Id:       providedOpId,
				OpPkgRef: providedOpPkgRef,
				RootOpId: providedRootOpId,
				Op:       &dcgOpDescriptor{},
			}

			fakeDcgNodeRepo := new(fakeDcgNodeRepo)

			objectUnderTest := newOpCaller(
				new(pkg.Fake),
				new(pubsub.Fake),
				fakeDcgNodeRepo,
				new(fakeCaller),
				new(uniquestring.Fake),
				new(validate.Fake),
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedOpId,
				providedOpPkgRef,
				providedRootOpId,
			)

			/* assert */
			Expect(fakeDcgNodeRepo.AddArgsForCall(0)).To(Equal(expectedDcgNodeDescriptor))
		})
		It("should call pkg.GetOp w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedOpId := "dummyOpId"
			providedOpPkgRef := "dummyOpPkgRef"
			providedRootOpId := "dummyRootOpId"

			expectedOpPkgRef := providedOpPkgRef

			fakePkg := new(pkg.Fake)

			objectUnderTest := newOpCaller(
				fakePkg,
				new(pubsub.Fake),
				new(fakeDcgNodeRepo),
				new(fakeCaller),
				new(uniquestring.Fake),
				new(validate.Fake),
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedOpId,
				providedOpPkgRef,
				providedRootOpId,
			)

			/* assert */
			Expect(fakePkg.GetOpArgsForCall(0)).To(Equal(expectedOpPkgRef))
		})
		Context("pkg.GetOp errors", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedOpId := "dummyOpId"
				providedOpPkgRef := "dummyOpPkgRef"
				providedRootOpId := "dummyRootOpId"

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					OpEncounteredError: &model.OpEncounteredErrorEvent{
						Msg:      "dummyError",
						OpId:     providedOpId,
						OpPkgRef: providedOpPkgRef,
						RootOpId: providedRootOpId,
					},
				}

				fakePkg := new(pkg.Fake)
				fakePkg.GetOpReturns(
					model.OpView{},
					errors.New(expectedEvent.OpEncounteredError.Msg),
				)

				fakeDcgNodeRepo := new(fakeDcgNodeRepo)
				fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := newOpCaller(
					fakePkg,
					fakePubSub,
					fakeDcgNodeRepo,
					new(fakeCaller),
					new(uniquestring.Fake),
					new(validate.Fake),
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedOpId,
					providedOpPkgRef,
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
		Context("pkg.GetOp doesn't error", func() {
			It("should call validate.Param w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{
					"dummyVar1Name": {String: "dummyVar1Data"},
					"dummyVar2Name": {File: "dummyVar2Data"},
					"dummyVar3Name": {Dir: "dummyVar3Data"},
					"dummyVar4Name": {Socket: "dummyVar4Data"},
				}
				providedOpId := "dummyOpId"
				providedOpPkgRef := "dummyOpPkgRef"
				providedRootOpId := "dummyRootOpId"

				opReturnedFromPkg := model.OpView{
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
				fakePkg := new(pkg.Fake)
				fakePkg.GetOpReturns(opReturnedFromPkg, nil)

				expectedCalls := map[*model.Data]*model.Param{}
				for inputName, input := range opReturnedFromPkg.Inputs {
					expectedCalls[providedInboundScope[inputName]] = input
				}

				fakeValidate := new(validate.Fake)

				objectUnderTest := newOpCaller(
					fakePkg,
					new(pubsub.Fake),
					new(fakeDcgNodeRepo),
					new(fakeCaller),
					new(uniquestring.Fake),
					fakeValidate,
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedOpId,
					providedOpPkgRef,
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
					providedOpPkgRef := "dummyOpPkgRef"
					providedRootOpId := "dummyRootOpId"

					fakeDcgNodeRepo := new(fakeDcgNodeRepo)
					fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

					opReturnedFromPkg := model.OpView{
						Inputs: map[string]*model.Param{
							"dummyVar1Name": {
								String: &model.StringParam{
									IsSecret: true,
								},
							},
						},
					}
					fakePkg := new(pkg.Fake)
					fakePkg.GetOpReturns(opReturnedFromPkg, nil)

					fakeValidate := new(validate.Fake)

					errorReturnedFromValidate := "validationError0"
					fakeValidate.ParamReturns([]error{errors.New(errorReturnedFromValidate)})

					expectedMsg := fmt.Sprintf(`
-
  validation of the following op input(s) failed:

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
							OpPkgRef: providedOpPkgRef,
							RootOpId: providedRootOpId,
						},
					}

					objectUnderTest := newOpCaller(
						fakePkg,
						fakePubSub,
						fakeDcgNodeRepo,
						new(fakeCaller),
						new(uniquestring.Fake),
						fakeValidate,
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						providedOpPkgRef,
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
					providedOpPkgRef := "dummyOpPkgRef"
					providedRootOpId := "dummyRootOpId"

					expectedEvent := &model.Event{
						Timestamp: time.Now().UTC(),
						OpStarted: &model.OpStartedEvent{
							OpId:     providedOpId,
							OpPkgRef: providedOpPkgRef,
							RootOpId: providedRootOpId,
						},
					}

					fakePubSub := new(pubsub.Fake)

					objectUnderTest := newOpCaller(
						new(pkg.Fake),
						fakePubSub,
						new(fakeDcgNodeRepo),
						new(fakeCaller),
						new(uniquestring.Fake),
						new(validate.Fake),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						providedOpPkgRef,
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
					providedOpPkgRef := "dummyOpPkgRef"
					providedRootOpId := "dummyRootOpId"

					opReturnedFromPkg := model.OpView{
						Run: &model.Scg{
							Parallel: []*model.Scg{
								{
									Container: &model.ScgContainerCall{},
								},
							},
						},
					}
					fakePkg := new(pkg.Fake)
					fakePkg.GetOpReturns(opReturnedFromPkg, nil)

					fakeCaller := new(fakeCaller)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					expectedNodeId := "dummyNodeId"
					fakeUniqueStringFactory.ConstructReturns(expectedNodeId)

					objectUnderTest := newOpCaller(
						fakePkg,
						new(pubsub.Fake),
						new(fakeDcgNodeRepo),
						fakeCaller,
						fakeUniqueStringFactory,
						new(validate.Fake),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						providedOpPkgRef,
						providedRootOpId,
					)

					/* assert */
					actualNodeId,
						actualInboundScope,
						actualScg,
						actualOpPkgRef,
						actualRootOpId := fakeCaller.CallArgsForCall(0)

					Expect(actualNodeId).To(Equal(expectedNodeId))
					Expect(actualInboundScope).To(Equal(providedInboundScope))
					Expect(actualScg).To(Equal(opReturnedFromPkg.Run))
					Expect(actualOpPkgRef).To(Equal(providedOpPkgRef))
					Expect(actualRootOpId).To(Equal(providedRootOpId))
				})
				It("should call dcgNodeRepo.GetIfExists w/ expected args", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedOpId := "dummyOpId"
					providedOpPkgRef := "dummyOpPkgRef"
					providedRootOpId := "dummyRootOpId"

					fakeDcgNodeRepo := new(fakeDcgNodeRepo)

					objectUnderTest := newOpCaller(
						new(pkg.Fake),
						new(pubsub.Fake),
						fakeDcgNodeRepo,
						new(fakeCaller),
						new(uniquestring.Fake),
						new(validate.Fake),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						providedOpPkgRef,
						providedRootOpId,
					)

					/* assert */
					Expect(fakeDcgNodeRepo.GetIfExistsArgsForCall(0)).To(Equal(providedRootOpId))
				})
				Context("dcgNodeRepo.GetIfExists returns nil", func() {
					It("should call pubSub.Publish w/ expected args", func() {
						/* arrange */
						providedInboundScope := map[string]*model.Data{}
						providedOpId := "dummyOpId"
						providedOpPkgRef := "dummyOpPkgRef"
						providedRootOpId := "dummyRootOpId"

						expectedEvent := &model.Event{
							Timestamp: time.Now().UTC(),
							OpEnded: &model.OpEndedEvent{
								OpId:     providedOpId,
								Outcome:  model.OpOutcomeKilled,
								RootOpId: providedRootOpId,
								OpPkgRef: providedOpPkgRef,
							},
						}

						fakePubSub := new(pubsub.Fake)

						objectUnderTest := newOpCaller(
							new(pkg.Fake),
							fakePubSub,
							new(fakeDcgNodeRepo),
							new(fakeCaller),
							new(uniquestring.Fake),
							new(validate.Fake),
						)

						/* act */
						objectUnderTest.Call(
							providedInboundScope,
							providedOpId,
							providedOpPkgRef,
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
						providedOpPkgRef := "dummyOpPkgRef"
						providedRootOpId := "dummyRootOpId"

						fakeDcgNodeRepo := new(fakeDcgNodeRepo)
						fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						objectUnderTest := newOpCaller(
							new(pkg.Fake),
							new(pubsub.Fake),
							fakeDcgNodeRepo,
							new(fakeCaller),
							new(uniquestring.Fake),
							new(validate.Fake),
						)

						/* act */
						objectUnderTest.Call(
							providedInboundScope,
							providedOpId,
							providedOpPkgRef,
							providedRootOpId,
						)

						/* assert */
						Expect(fakeDcgNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedOpId))
					})
					Context("caller.Call errored", func() {
						It("should call pubSub.Publish w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOpId := "dummyOpId"
							providedOpPkgRef := "dummyOpPkgRef"
							providedRootOpId := "dummyRootOpId"

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEncounteredError: &model.OpEncounteredErrorEvent{
									Msg:      "dummyError",
									OpId:     providedOpId,
									OpPkgRef: providedOpPkgRef,
									RootOpId: providedRootOpId,
								},
							}

							fakeDcgNodeRepo := new(fakeDcgNodeRepo)
							fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							fakeCaller := new(fakeCaller)
							fakeCaller.CallReturns(
								map[string]*model.Data{},
								errors.New(expectedEvent.OpEncounteredError.Msg),
							)

							objectUnderTest := newOpCaller(
								new(pkg.Fake),
								fakePubSub,
								fakeDcgNodeRepo,
								fakeCaller,
								new(uniquestring.Fake),
								new(validate.Fake),
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								providedOpPkgRef,
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
							providedOpPkgRef := "dummyOpPkgRef"
							providedRootOpId := "dummyRootOpId"

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEnded: &model.OpEndedEvent{
									OpId:     providedOpId,
									OpPkgRef: providedOpPkgRef,
									Outcome:  model.OpOutcomeFailed,
									RootOpId: providedRootOpId,
								},
							}

							fakeDcgNodeRepo := new(fakeDcgNodeRepo)
							fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							fakeCaller := new(fakeCaller)
							fakeCaller.CallReturns(
								map[string]*model.Data{},
								errors.New("dummyError"),
							)

							objectUnderTest := newOpCaller(
								new(pkg.Fake),
								fakePubSub,
								fakeDcgNodeRepo,
								fakeCaller,
								new(uniquestring.Fake),
								new(validate.Fake),
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								providedOpPkgRef,
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
							providedOpPkgRef := "dummyOpPkgRef"
							providedRootOpId := "dummyRootOpId"

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEnded: &model.OpEndedEvent{
									OpId:     providedOpId,
									OpPkgRef: providedOpPkgRef,
									Outcome:  model.OpOutcomeSucceeded,
									RootOpId: providedRootOpId,
								},
							}

							fakeDcgNodeRepo := new(fakeDcgNodeRepo)
							fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							objectUnderTest := newOpCaller(
								new(pkg.Fake),
								fakePubSub,
								fakeDcgNodeRepo,
								new(fakeCaller),
								new(uniquestring.Fake),
								new(validate.Fake),
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								providedOpPkgRef,
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
