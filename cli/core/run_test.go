package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/clioutput"
	"github.com/opspec-io/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/opctl/util/colorer"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/engineclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
	"path/filepath"
	"time"
)

var _ = Context("runOp", func() {
	Context("Execute", func() {
		Context("vos.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeVos := new(vos.Fake)
				expectedError := errors.New("dummyError")
				fakeVos.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					bundle:    new(bundle.Fake),
					cliExiter: fakeCliExiter,
					vos:       fakeVos,
				}

				/* act */
				objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyName")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call bundle.GetOp w/ expected args", func() {
				/* arrange */
				fakeBundle := new(bundle.Fake)

				fakeEngineClient := new(engineclient.Fake)
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

				fakeCliExiter := new(cliexiter.Fake)

				providedName := "dummyOpName"
				providedCollection := "dummyCollection"
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.Fake)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)
				expectedPath := filepath.Join(wdReturnedFromVos, providedCollection, providedName)

				objectUnderTest := _core{
					bundle:            fakeBundle,
					engineClient:      fakeEngineClient,
					cliExiter:         fakeCliExiter,
					cliParamSatisfier: new(cliparamsatisfier.Fake),
					vos:               fakeVos,
				}

				/* act */
				objectUnderTest.RunOp([]string{}, providedCollection, providedName)

				/* assert */
				Expect(fakeBundle.GetOpArgsForCall(0)).Should(Equal(expectedPath))
			})
			Context("bundle.GetOp errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiter.Fake)
					returnedError := errors.New("dummyError")

					fakeBundle := new(bundle.Fake)
					fakeBundle.GetOpReturns(model.OpView{}, returnedError)

					objectUnderTest := _core{
						bundle:            fakeBundle,
						cliExiter:         fakeCliExiter,
						cliParamSatisfier: new(cliparamsatisfier.Fake),
						vos:               new(vos.Fake),
					}

					/* act */
					objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyName")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
				})
			})
			Context("bundle.GetOp doesn't error", func() {
				It("should call paramSatisfier.Satisfy w/ expected args", func() {
					/* arrange */
					param1Name := "DUMMY_PARAM1_NAME"
					arg1Value := &model.Data{String: "dummyParam1Value"}

					providedArgs := []string{fmt.Sprintf("%v=%v", param1Name, arg1Value.String)}

					expectedParams := map[string]*model.Param{
						param1Name: {
							String: &model.StringParam{},
						},
					}

					fakeBundle := new(bundle.Fake)
					fakeBundle.GetOpReturns(
						model.OpView{
							Inputs: expectedParams,
						},
						nil,
					)

					// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
					fakeEngineClient := new(engineclient.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

					fakeCliParamSatisfier := new(cliparamsatisfier.Fake)

					objectUnderTest := _core{
						bundle:            fakeBundle,
						engineClient:      fakeEngineClient,
						cliExiter:         new(cliexiter.Fake),
						cliParamSatisfier: fakeCliParamSatisfier,
						vos:               new(vos.Fake),
					}

					/* act */
					objectUnderTest.RunOp(providedArgs, "dummyCollection", "dummyOpName")

					/* assert */
					actualArgs, actualParams := fakeCliParamSatisfier.SatisfyArgsForCall(0)
					Expect(actualArgs).To(Equal(providedArgs))
					Expect(actualParams).To(Equal(expectedParams))
				})
				It("should call engineClient.StartOp w/ expected args", func() {
					/* arrange */
					pwd := "dummyWorkDir"
					fakeVos := new(vos.Fake)
					fakeVos.GetwdReturns(pwd, nil)

					providedCollection := "dummyCollection"
					providedOp := "dummyOp"
					expectedArgs := model.StartOpReq{
						Args: map[string]*model.Data{
							"dummyArg1Name": {String: "dummyArg1Value"},
						},
						OpRef: path.Join(pwd, providedCollection, providedOp),
					}

					// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
					fakeEngineClient := new(engineclient.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

					fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
					fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args)

					objectUnderTest := _core{
						bundle:            new(bundle.Fake),
						engineClient:      fakeEngineClient,
						cliExiter:         new(cliexiter.Fake),
						cliParamSatisfier: fakeCliParamSatisfier,
						vos:               fakeVos,
					}

					/* act */
					objectUnderTest.RunOp([]string{}, providedCollection, providedOp)

					/* assert */
					actualArgs := fakeEngineClient.StartOpArgsForCall(0)
					Expect(actualArgs).To(Equal(expectedArgs))
				})
				Context("engineClient.StartOp errors", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						fakeCliExiter := new(cliexiter.Fake)
						returnedError := errors.New("dummyError")

						fakeBundle := new(bundle.Fake)
						fakeBundle.GetOpReturns(model.OpView{}, nil)

						fakeEngineClient := new(engineclient.Fake)
						fakeEngineClient.StartOpReturns("dummyOpId", returnedError)

						objectUnderTest := _core{
							bundle:            fakeBundle,
							engineClient:      fakeEngineClient,
							cliExiter:         fakeCliExiter,
							cliParamSatisfier: new(cliparamsatisfier.Fake),
							vos:               new(vos.Fake),
						}

						/* act */
						objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

						/* assert */
						Expect(fakeCliExiter.ExitArgsForCall(0)).
							Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
					})
				})
				Context("engineClient.StartOp doesn't error", func() {
					It("should call engineClient.GetEventStream w/ expected args", func() {
						/* arrange */
						fakeBundle := new(bundle.Fake)
						fakeBundle.GetOpReturns(model.OpView{}, nil)
						opGraphIdReturnedFromStartOp := "dummyOpGraphId"
						expectedEventFilter := &model.GetEventStreamReq{
							Filter: &model.EventFilter{
								OpGraphIds: []string{opGraphIdReturnedFromStartOp},
							},
						}

						fakeEngineClient := new(engineclient.Fake)
						fakeEngineClient.StartOpReturns(opGraphIdReturnedFromStartOp, nil)
						eventChannel := make(chan model.Event)
						close(eventChannel)
						fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

						objectUnderTest := _core{
							bundle:            fakeBundle,
							engineClient:      fakeEngineClient,
							cliExiter:         new(cliexiter.Fake),
							cliParamSatisfier: new(cliparamsatisfier.Fake),
							vos:               new(vos.Fake),
						}

						/* act */
						objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

						/* assert */
						Expect(fakeEngineClient.GetEventStreamArgsForCall(0)).
							Should(Equal(expectedEventFilter))
					})
					Context("engineClient.GetEventStream errors", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeCliExiter := new(cliexiter.Fake)
							returnedError := errors.New("dummyError")

							fakeBundle := new(bundle.Fake)
							fakeBundle.GetOpReturns(model.OpView{}, nil)

							fakeEngineClient := new(engineclient.Fake)
							fakeEngineClient.GetEventStreamReturns(nil, returnedError)

							objectUnderTest := _core{
								bundle:            fakeBundle,
								engineClient:      fakeEngineClient,
								cliExiter:         fakeCliExiter,
								cliParamSatisfier: new(cliparamsatisfier.Fake),
								vos:               new(vos.Fake),
							}

							/* act */
							objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

							/* assert */
							Expect(fakeCliExiter.ExitArgsForCall(0)).
								Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
						})
					})
					Context("engineClient.GetEventStream doesn't error", func() {
						Context("event channel closes", func() {
							It("should call exiter w/ expected args", func() {
								/* arrange */
								fakeCliExiter := new(cliexiter.Fake)

								fakeBundle := new(bundle.Fake)
								fakeBundle.GetOpReturns(model.OpView{}, nil)

								fakeEngineClient := new(engineclient.Fake)
								eventChannel := make(chan model.Event)
								close(eventChannel)
								fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

								objectUnderTest := _core{
									bundle:            fakeBundle,
									engineClient:      fakeEngineClient,
									cliExiter:         fakeCliExiter,
									cliParamSatisfier: new(cliparamsatisfier.Fake),
									vos:               new(vos.Fake),
								}

								/* act */
								objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

								/* assert */
								Expect(fakeCliExiter.ExitArgsForCall(0)).
									Should(Equal(cliexiter.ExitReq{Message: "Event channel closed unexpectedly", Code: 1}))
							})
						})
						Context("event channel doesn't close", func() {
							Context("event received", func() {
								opGraphId := "dummyOpGraphId"
								Context("OpEndedEvent", func() {
									Context("Outcome==SUCCEEDED", func() {
										It("should call exiter w/ expected args", func() {
											/* arrange */
											opEndedEvent := model.Event{
												Timestamp: time.Now(),
												OpEnded: &model.OpEndedEvent{
													OpId:      opGraphId,
													OpRef:     "dummyOpRef",
													Outcome:   model.OpOutcomeSucceeded,
													OpGraphId: opGraphId,
												},
											}

											fakeCliExiter := new(cliexiter.Fake)

											fakeBundle := new(bundle.Fake)
											fakeBundle.GetOpReturns(model.OpView{}, nil)

											fakeEngineClient := new(engineclient.Fake)
											eventChannel := make(chan model.Event, 10)
											eventChannel <- opEndedEvent
											defer close(eventChannel)
											fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
											fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

											objectUnderTest := _core{
												bundle:            fakeBundle,
												colorer:           colorer.New(),
												engineClient:      fakeEngineClient,
												cliExiter:         fakeCliExiter,
												cliOutput:         new(clioutput.Fake),
												cliParamSatisfier: new(cliparamsatisfier.Fake),
												vos:               new(vos.Fake),
											}

											/* act/assert */
											objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
											Expect(fakeCliExiter.ExitArgsForCall(0)).
												Should(Equal(cliexiter.ExitReq{Code: 0}))
										})
									})
									Context("Outcome==KILLED", func() {
										It("should call exiter w/ expected args", func() {
											/* arrange */
											opEndedEvent := model.Event{
												Timestamp: time.Now(),
												OpEnded: &model.OpEndedEvent{
													OpId:      opGraphId,
													OpRef:     "dummyOpRef",
													Outcome:   model.OpOutcomeKilled,
													OpGraphId: opGraphId,
												},
											}

											fakeCliExiter := new(cliexiter.Fake)

											fakeBundle := new(bundle.Fake)
											fakeBundle.GetOpReturns(model.OpView{}, nil)

											fakeEngineClient := new(engineclient.Fake)
											eventChannel := make(chan model.Event, 10)
											eventChannel <- opEndedEvent
											defer close(eventChannel)
											fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
											fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

											objectUnderTest := _core{
												bundle:            fakeBundle,
												colorer:           colorer.New(),
												engineClient:      fakeEngineClient,
												cliExiter:         fakeCliExiter,
												cliOutput:         new(clioutput.Fake),
												cliParamSatisfier: new(cliparamsatisfier.Fake),
												vos:               new(vos.Fake),
											}

											/* act/assert */
											objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
											Expect(fakeCliExiter.ExitArgsForCall(0)).
												Should(Equal(cliexiter.ExitReq{Code: 137}))
										})

									})
									Context("Outcome==FAILED", func() {
										It("should call exiter w/ expected args", func() {
											/* arrange */
											opEndedEvent := model.Event{
												Timestamp: time.Now(),
												OpEnded: &model.OpEndedEvent{
													OpId:      opGraphId,
													OpRef:     "dummyOpRef",
													Outcome:   model.OpOutcomeFailed,
													OpGraphId: opGraphId,
												},
											}

											fakeCliExiter := new(cliexiter.Fake)

											fakeBundle := new(bundle.Fake)
											fakeBundle.GetOpReturns(model.OpView{}, nil)

											fakeEngineClient := new(engineclient.Fake)
											eventChannel := make(chan model.Event, 10)
											eventChannel <- opEndedEvent
											defer close(eventChannel)
											fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
											fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

											objectUnderTest := _core{
												bundle:            fakeBundle,
												colorer:           colorer.New(),
												engineClient:      fakeEngineClient,
												cliExiter:         fakeCliExiter,
												cliOutput:         new(clioutput.Fake),
												cliParamSatisfier: new(cliparamsatisfier.Fake),
												vos:               new(vos.Fake),
											}

											/* act/assert */
											objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
											Expect(fakeCliExiter.ExitArgsForCall(0)).
												Should(Equal(cliexiter.ExitReq{Code: 1}))
										})
									})
									Context("Outcome==?", func() {
										It("should call exiter w/ expected args", func() {
											/* arrange */
											opEndedEvent := model.Event{
												Timestamp: time.Now(),
												OpEnded: &model.OpEndedEvent{
													OpId:      opGraphId,
													OpRef:     "dummyOpRef",
													Outcome:   "some unexpected outcome",
													OpGraphId: opGraphId,
												},
											}

											fakeCliExiter := new(cliexiter.Fake)

											fakeBundle := new(bundle.Fake)
											fakeBundle.GetOpReturns(model.OpView{}, nil)

											fakeEngineClient := new(engineclient.Fake)
											eventChannel := make(chan model.Event, 10)
											eventChannel <- opEndedEvent
											defer close(eventChannel)
											fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
											fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

											objectUnderTest := _core{
												bundle:            fakeBundle,
												colorer:           colorer.New(),
												engineClient:      fakeEngineClient,
												cliExiter:         fakeCliExiter,
												cliOutput:         new(clioutput.Fake),
												cliParamSatisfier: new(cliparamsatisfier.Fake),
												vos:               new(vos.Fake),
											}

											/* act/assert */
											objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
											Expect(fakeCliExiter.ExitArgsForCall(0)).
												Should(Equal(cliexiter.ExitReq{Code: 1}))
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})
})
