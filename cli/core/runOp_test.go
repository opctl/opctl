package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/colorer"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/engineclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
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

				fakeExiter := new(fakeExiter)

				objectUnderTest := _core{
					bundle: new(bundle.Fake),
					exiter: fakeExiter,
					vos:    fakeVos,
				}

				/* act */
				objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyName")

				/* assert */
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
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

				fakeExiter := new(fakeExiter)

				providedName := "dummyOpName"
				providedCollection := "dummyCollection"
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.Fake)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)
				expectedPath := filepath.Join(wdReturnedFromVos, providedCollection, providedName)

				objectUnderTest := _core{
					bundle:         fakeBundle,
					engineClient:   fakeEngineClient,
					exiter:         fakeExiter,
					paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
					vos:            fakeVos,
				}

				/* act */
				objectUnderTest.RunOp([]string{}, providedCollection, providedName)

				/* assert */
				Expect(fakeBundle.GetOpArgsForCall(0)).Should(Equal(expectedPath))
			})
			Context("bundle.GetOp errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeExiter := new(fakeExiter)
					returnedError := errors.New("dummyError")

					fakeBundle := new(bundle.Fake)
					fakeBundle.GetOpReturns(model.OpView{}, returnedError)

					objectUnderTest := _core{
						bundle:         fakeBundle,
						exiter:         fakeExiter,
						paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
						vos:            new(vos.Fake),
					}

					/* act */
					objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyName")

					/* assert */
					Expect(fakeExiter.ExitArgsForCall(0)).
						Should(Equal(ExitReq{Message: returnedError.Error(), Code: 1}))
				})
			})
			Context("bundle.GetOp doesn't error", func() {
				Context("op has params", func() {
					Context("op args provided explicitly w/ values", func() {
						It("should call engineClient.StartOp with provided arg values", func() {
							/* arrange */
							fakeExiter := new(fakeExiter)
							param1Name := "DUMMY_PARAM1_NAME"
							param1Value := &model.Data{String: "dummyParam1Value"}

							fakeBundle := new(bundle.Fake)
							fakeBundle.GetOpReturns(
								model.OpView{
									Inputs: []*model.Param{
										{
											String: &model.StringParam{
												Name: param1Name,
											},
										},
									},
								},
								nil,
							)

							fakeEngineClient := new(engineclient.Fake)
							fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

							objectUnderTest := _core{
								bundle:         fakeBundle,
								engineClient:   fakeEngineClient,
								exiter:         fakeExiter,
								paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
								vos:            new(vos.Fake),
							}

							expectedArgs := map[string]*model.Data{param1Name: param1Value}
							providedArgs := []string{fmt.Sprintf("%v=%v", param1Name, param1Value.String)}

							/* act */
							objectUnderTest.RunOp(providedArgs, "dummyCollection", "dummyOpName")

							/* assert */
							Expect(fakeEngineClient.StartOpArgsForCall(0).Args).To(BeEquivalentTo(expectedArgs))
						})
					})
					Context("op args provided explicitly w/out values", func() {
						It("should call engineClient.StartOp with arg values obtained from the environment", func() {
							/* arrange */
							fakeExiter := new(fakeExiter)
							param1Name := "DUMMY_PARAM1_NAME"
							param1Value := &model.Data{String: "dummyParam1Value"}

							fakeVos := new(vos.Fake)
							fakeVos.GetenvReturns(param1Value.String)

							fakeBundle := new(bundle.Fake)
							fakeBundle.GetOpReturns(
								model.OpView{
									Inputs: []*model.Param{
										{
											String: &model.StringParam{
												Name: param1Name,
											},
										},
									},
								},
								nil,
							)

							fakeEngineClient := new(engineclient.Fake)
							fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

							objectUnderTest := _core{
								bundle:         fakeBundle,
								engineClient:   fakeEngineClient,
								exiter:         fakeExiter,
								paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), fakeVos),
								vos:            new(vos.Fake),
							}

							expectedArgs := map[string]*model.Data{param1Name: param1Value}
							providedArgs := []string{param1Name}

							/* act */
							objectUnderTest.RunOp(providedArgs, "dummyCollection", "dummyOpName")

							/* assert */
							Expect(fakeEngineClient.StartOpArgsForCall(0).Args).To(BeEquivalentTo(expectedArgs))
						})
					})
					Context("op args not provided", func() {
						Context("op params don't have defaults", func() {
							It("should call bundle.RunOp with arg values obtained from the environment", func() {
								/* arrange */
								fakeExiter := new(fakeExiter)
								param1Name := "DUMMY_PARAM1_NAME"
								param1Value := &model.Data{String: "dummyParam1Value"}

								fakeVos := new(vos.Fake)
								fakeVos.GetenvReturns(param1Value.String)

								fakeBundle := new(bundle.Fake)
								fakeBundle.GetOpReturns(
									model.OpView{
										Inputs: []*model.Param{
											{
												String: &model.StringParam{
													Name: param1Name,
												},
											},
										},
									},
									nil,
								)

								fakeEngineClient := new(engineclient.Fake)
								fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

								objectUnderTest := _core{
									bundle:         fakeBundle,
									engineClient:   fakeEngineClient,
									exiter:         fakeExiter,
									paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), fakeVos),
									vos:            new(vos.Fake),
								}

								expectedArgs := map[string]*model.Data{param1Name: param1Value}
								providedArgs := []string{}

								/* act */
								objectUnderTest.RunOp(providedArgs, "dummyCollection", "dummyOpName")

								/* assert */
								Expect(fakeEngineClient.StartOpArgsForCall(0).Args).To(BeEquivalentTo(expectedArgs))
							})
						})
						Context("op params have defaults", func() {
							It("should call bundle.RunOp without args for defaulted params", func() {
								/* arrange */
								fakeExiter := new(fakeExiter)
								// unique name to ensure conflicting env var not present
								param1Name := string(time.Now().Unix())

								fakeBundle := new(bundle.Fake)
								fakeBundle.GetOpReturns(
									model.OpView{
										Inputs: []*model.Param{
											{
												String: &model.StringParam{
													Name:    param1Name,
													Default: "dummyDefault",
												},
											},
										},
									},
									nil,
								)

								fakeEngineClient := new(engineclient.Fake)
								fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

								objectUnderTest := _core{
									bundle:         fakeBundle,
									engineClient:   fakeEngineClient,
									exiter:         fakeExiter,
									paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
									vos:            new(vos.Fake),
								}

								expectedArgs := map[string]*model.Data{}
								providedArgs := []string{}

								/* act */
								objectUnderTest.RunOp(providedArgs, "dummyCollection", "dummyOpName")

								/* assert */
								Expect(fakeEngineClient.StartOpArgsForCall(0).Args).To(BeEquivalentTo(expectedArgs))
							})
						})
					})
				})
				Context("op doesn't have params", func() {
					Context("engineClient.StartOp errors", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeExiter := new(fakeExiter)
							returnedError := errors.New("dummyError")

							fakeBundle := new(bundle.Fake)
							fakeBundle.GetOpReturns(model.OpView{}, nil)

							fakeEngineClient := new(engineclient.Fake)
							fakeEngineClient.StartOpReturns("dummyOpId", returnedError)

							objectUnderTest := _core{
								bundle:         fakeBundle,
								engineClient:   fakeEngineClient,
								exiter:         fakeExiter,
								paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
								vos:            new(vos.Fake),
							}

							/* act */
							objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

							/* assert */
							Expect(fakeExiter.ExitArgsForCall(0)).
								Should(Equal(ExitReq{Message: returnedError.Error(), Code: 1}))
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
								bundle:         fakeBundle,
								engineClient:   fakeEngineClient,
								exiter:         new(fakeExiter),
								paramSatisfier: newParamSatisfier(colorer.New(), new(fakeExiter), validate.New(), new(vos.Fake)),
								vos:            new(vos.Fake),
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
								fakeExiter := new(fakeExiter)
								returnedError := errors.New("dummyError")

								fakeBundle := new(bundle.Fake)
								fakeBundle.GetOpReturns(model.OpView{}, nil)

								fakeEngineClient := new(engineclient.Fake)
								fakeEngineClient.GetEventStreamReturns(nil, returnedError)

								objectUnderTest := _core{
									bundle:         fakeBundle,
									engineClient:   fakeEngineClient,
									exiter:         fakeExiter,
									paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
									vos:            new(vos.Fake),
								}

								/* act */
								objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

								/* assert */
								Expect(fakeExiter.ExitArgsForCall(0)).
									Should(Equal(ExitReq{Message: returnedError.Error(), Code: 1}))
							})
						})
						Context("engineClient.GetEventStream doesn't error", func() {
							Context("event channel closes", func() {
								It("should call exiter w/ expected args", func() {
									/* arrange */
									fakeExiter := new(fakeExiter)

									fakeBundle := new(bundle.Fake)
									fakeBundle.GetOpReturns(model.OpView{}, nil)

									fakeEngineClient := new(engineclient.Fake)
									eventChannel := make(chan model.Event)
									close(eventChannel)
									fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

									objectUnderTest := _core{
										bundle:         fakeBundle,
										engineClient:   fakeEngineClient,
										exiter:         fakeExiter,
										paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
										vos:            new(vos.Fake),
									}

									/* act */
									objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

									/* assert */
									Expect(fakeExiter.ExitArgsForCall(0)).
										Should(Equal(ExitReq{Message: "Event channel closed unexpectedly", Code: 1}))
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

												fakeExiter := new(fakeExiter)

												fakeBundle := new(bundle.Fake)
												fakeBundle.GetOpReturns(model.OpView{}, nil)

												fakeEngineClient := new(engineclient.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
												fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

												objectUnderTest := _core{
													bundle:         fakeBundle,
													colorer:        colorer.New(),
													engineClient:   fakeEngineClient,
													exiter:         fakeExiter,
													output:         new(fakeOutput),
													paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
													vos:            new(vos.Fake),
												}

												/* act/assert */
												objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
												Expect(fakeExiter.ExitArgsForCall(0)).
													Should(Equal(ExitReq{Code: 0}))
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

												fakeExiter := new(fakeExiter)

												fakeBundle := new(bundle.Fake)
												fakeBundle.GetOpReturns(model.OpView{}, nil)

												fakeEngineClient := new(engineclient.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
												fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

												objectUnderTest := _core{
													bundle:         fakeBundle,
													colorer:        colorer.New(),
													engineClient:   fakeEngineClient,
													exiter:         fakeExiter,
													output:         new(fakeOutput),
													paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
													vos:            new(vos.Fake),
												}

												/* act/assert */
												objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
												Expect(fakeExiter.ExitArgsForCall(0)).
													Should(Equal(ExitReq{Code: 137}))
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

												fakeExiter := new(fakeExiter)

												fakeBundle := new(bundle.Fake)
												fakeBundle.GetOpReturns(model.OpView{}, nil)

												fakeEngineClient := new(engineclient.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
												fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

												objectUnderTest := _core{
													bundle:         fakeBundle,
													colorer:        colorer.New(),
													engineClient:   fakeEngineClient,
													exiter:         fakeExiter,
													output:         new(fakeOutput),
													paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
													vos:            new(vos.Fake),
												}

												/* act/assert */
												objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
												Expect(fakeExiter.ExitArgsForCall(0)).
													Should(Equal(ExitReq{Code: 1}))
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

												fakeExiter := new(fakeExiter)

												fakeBundle := new(bundle.Fake)
												fakeBundle.GetOpReturns(model.OpView{}, nil)

												fakeEngineClient := new(engineclient.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
												fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

												objectUnderTest := _core{
													bundle:         fakeBundle,
													colorer:        colorer.New(),
													engineClient:   fakeEngineClient,
													exiter:         fakeExiter,
													output:         new(fakeOutput),
													paramSatisfier: newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.Fake)),
													vos:            new(vos.Fake),
												}

												/* act/assert */
												objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
												Expect(fakeExiter.ExitArgsForCall(0)).
													Should(Equal(ExitReq{Code: 1}))
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
})
