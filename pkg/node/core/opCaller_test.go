package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
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
				new(bundle.Fake),
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
			providedOpRef := "dummyOpRef"
			providedOpGraphId := "dummyOpGraphId"

			expectedDcgNodeDescriptor := &dcgNodeDescriptor{
				Id:        providedOpId,
				OpRef:     providedOpRef,
				OpGraphId: providedOpGraphId,
				Op:        &dcgOpDescriptor{},
			}

			fakeDcgNodeRepo := new(fakeDcgNodeRepo)

			objectUnderTest := newOpCaller(
				new(bundle.Fake),
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
				providedOpRef,
				providedOpGraphId,
			)

			/* assert */
			Expect(fakeDcgNodeRepo.AddArgsForCall(0)).To(Equal(expectedDcgNodeDescriptor))
		})
		It("should call bundle.GetOp w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedOpId := "dummyOpId"
			providedOpRef := "dummyOpRef"
			providedOpGraphId := "dummyOpGraphId"

			expectedOpRef := providedOpRef

			fakeBundle := new(bundle.Fake)

			objectUnderTest := newOpCaller(
				fakeBundle,
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
				providedOpRef,
				providedOpGraphId,
			)

			/* assert */
			Expect(fakeBundle.GetOpArgsForCall(0)).To(Equal(expectedOpRef))
		})
		Context("bundle.GetOp errors", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedOpId := "dummyOpId"
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					OpEncounteredError: &model.OpEncounteredErrorEvent{
						Msg:       "dummyError",
						OpId:      providedOpId,
						OpRef:     providedOpRef,
						OpGraphId: providedOpGraphId,
					},
				}

				fakeBundle := new(bundle.Fake)
				fakeBundle.GetOpReturns(
					model.OpView{},
					errors.New(expectedEvent.OpEncounteredError.Msg),
				)

				fakeDcgNodeRepo := new(fakeDcgNodeRepo)
				fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := newOpCaller(
					fakeBundle,
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
					providedOpRef,
					providedOpGraphId,
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
		Context("bundle.GetOp doesn't error", func() {
			It("should call validate.Param w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{
					"dummyVar1Name": {String: "dummyVar1Data"},
					"dummyVar2Name": {File: "dummyVar2Data"},
					"dummyVar3Name": {Dir: "dummyVar3Data"},
					"dummyVar4Name": {Socket: "dummyVar4Data"},
				}
				providedOpId := "dummyOpId"
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				opReturnedFromBundle := model.OpView{
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
				fakeBundle := new(bundle.Fake)
				fakeBundle.GetOpReturns(opReturnedFromBundle, nil)

				expectedCalls := map[*model.Data]*model.Param{}
				for inputName, input := range opReturnedFromBundle.Inputs {
					expectedCalls[providedInboundScope[inputName]] = input
				}

				fakeValidate := new(validate.Fake)

				objectUnderTest := newOpCaller(
					fakeBundle,
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
					providedOpRef,
					providedOpGraphId,
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
					providedOpRef := "dummyOpRef"
					providedOpGraphId := "dummyOpGraphId"

					fakeDcgNodeRepo := new(fakeDcgNodeRepo)
					fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

					opReturnedFromBundle := model.OpView{
						Inputs: map[string]*model.Param{
							"dummyVar1Name": {
								String: &model.StringParam{
									IsSecret: true,
								},
							},
						},
					}
					fakeBundle := new(bundle.Fake)
					fakeBundle.GetOpReturns(opReturnedFromBundle, nil)

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
							Msg:       expectedMsg,
							OpId:      providedOpId,
							OpRef:     providedOpRef,
							OpGraphId: providedOpGraphId,
						},
					}

					objectUnderTest := newOpCaller(
						fakeBundle,
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
						providedOpRef,
						providedOpGraphId,
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
					providedOpRef := "dummyOpRef"
					providedOpGraphId := "dummyOpGraphId"

					expectedEvent := &model.Event{
						Timestamp: time.Now().UTC(),
						OpStarted: &model.OpStartedEvent{
							OpId:      providedOpId,
							OpRef:     providedOpRef,
							OpGraphId: providedOpGraphId,
						},
					}

					fakePubSub := new(pubsub.Fake)

					objectUnderTest := newOpCaller(
						new(bundle.Fake),
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
						providedOpRef,
						providedOpGraphId,
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
					providedOpRef := "dummyOpRef"
					providedOpGraphId := "dummyOpGraphId"

					opReturnedFromBundle := model.OpView{
						Run: &model.Scg{
							Parallel: []*model.Scg{
								{
									Container: &model.ScgContainerCall{},
								},
							},
						},
					}
					fakeBundle := new(bundle.Fake)
					fakeBundle.GetOpReturns(opReturnedFromBundle, nil)

					fakeCaller := new(fakeCaller)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					expectedNodeId := "dummyNodeId"
					fakeUniqueStringFactory.ConstructReturns(expectedNodeId)

					objectUnderTest := newOpCaller(
						fakeBundle,
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
						providedOpRef,
						providedOpGraphId,
					)

					/* assert */
					actualNodeId,
						actualInboundScope,
						actualScg,
						actualOpRef,
						actualOpGraphId := fakeCaller.CallArgsForCall(0)

					Expect(actualNodeId).To(Equal(expectedNodeId))
					Expect(actualInboundScope).To(Equal(providedInboundScope))
					Expect(actualScg).To(Equal(opReturnedFromBundle.Run))
					Expect(actualOpRef).To(Equal(providedOpRef))
					Expect(actualOpGraphId).To(Equal(providedOpGraphId))
				})
				It("should call dcgNodeRepo.GetIfExists w/ expected args", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedOpId := "dummyOpId"
					providedOpRef := "dummyOpRef"
					providedOpGraphId := "dummyOpGraphId"

					fakeDcgNodeRepo := new(fakeDcgNodeRepo)

					objectUnderTest := newOpCaller(
						new(bundle.Fake),
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
						providedOpRef,
						providedOpGraphId,
					)

					/* assert */
					Expect(fakeDcgNodeRepo.GetIfExistsArgsForCall(0)).To(Equal(providedOpGraphId))
				})
				Context("dcgNodeRepo.GetIfExists returns nil", func() {
					It("should call pubSub.Publish w/ expected args", func() {
						/* arrange */
						providedInboundScope := map[string]*model.Data{}
						providedOpId := "dummyOpId"
						providedOpRef := "dummyOpRef"
						providedOpGraphId := "dummyOpGraphId"

						expectedEvent := &model.Event{
							Timestamp: time.Now().UTC(),
							OpEnded: &model.OpEndedEvent{
								OpId:      providedOpId,
								Outcome:   model.OpOutcomeKilled,
								OpGraphId: providedOpGraphId,
								OpRef:     providedOpRef,
							},
						}

						fakePubSub := new(pubsub.Fake)

						objectUnderTest := newOpCaller(
							new(bundle.Fake),
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
							providedOpRef,
							providedOpGraphId,
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
						providedOpRef := "dummyOpRef"
						providedOpGraphId := "dummyOpGraphId"

						fakeDcgNodeRepo := new(fakeDcgNodeRepo)
						fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						objectUnderTest := newOpCaller(
							new(bundle.Fake),
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
							providedOpRef,
							providedOpGraphId,
						)

						/* assert */
						Expect(fakeDcgNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedOpId))
					})
					Context("caller.Call errored", func() {
						It("should call pubSub.Publish w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOpId := "dummyOpId"
							providedOpRef := "dummyOpRef"
							providedOpGraphId := "dummyOpGraphId"

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEncounteredError: &model.OpEncounteredErrorEvent{
									Msg:       "dummyError",
									OpId:      providedOpId,
									OpRef:     providedOpRef,
									OpGraphId: providedOpGraphId,
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
								new(bundle.Fake),
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
								providedOpRef,
								providedOpGraphId,
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
							providedOpRef := "dummyOpRef"
							providedOpGraphId := "dummyOpGraphId"

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEnded: &model.OpEndedEvent{
									OpId:      providedOpId,
									OpRef:     providedOpRef,
									Outcome:   model.OpOutcomeFailed,
									OpGraphId: providedOpGraphId,
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
								new(bundle.Fake),
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
								providedOpRef,
								providedOpGraphId,
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
							providedOpRef := "dummyOpRef"
							providedOpGraphId := "dummyOpGraphId"

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpEnded: &model.OpEndedEvent{
									OpId:      providedOpId,
									OpRef:     providedOpRef,
									Outcome:   model.OpOutcomeSucceeded,
									OpGraphId: providedOpGraphId,
								},
							}

							fakeDcgNodeRepo := new(fakeDcgNodeRepo)
							fakeDcgNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							fakePubSub := new(pubsub.Fake)

							objectUnderTest := newOpCaller(
								new(bundle.Fake),
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
								providedOpRef,
								providedOpGraphId,
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
