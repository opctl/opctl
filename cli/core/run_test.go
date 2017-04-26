package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/util/clicolorer"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/clioutput"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/sdk-golang/consumenodeapi"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/virtual-go/vioutil"
	"github.com/virtual-go/vos"
	"time"
)

var _ = Context("Run", func() {
	Context("Execute", func() {
		Context("os.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeVOS := new(vos.Fake)
				expectedError := errors.New("dummyError")
				fakeVOS.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:          new(pkg.Fake),
					cliExiter:    fakeCliExiter,
					nodeProvider: new(nodeprovider.Fake),
					os:           fakeVOS,
				}

				/* act */
				objectUnderTest.Run("dummyName", &RunOpts{})

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("os.Getwd doesn't error", func() {
			It("should call pkg.Resolve w/ expected args", func() {
				/* arrange */
				providedPkgRef := "dummyPkgName"

				expectedPkgBasePath := "dummyWorkDir"
				expectedPkgRef := providedPkgRef

				fakePkg := new(pkg.Fake)

				fakeCliExiter := new(cliexiter.Fake)

				fakeVOS := new(vos.Fake)
				fakeVOS.GetwdReturns(expectedPkgBasePath, nil)

				objectUnderTest := _core{
					pkg:          fakePkg,
					cliExiter:    fakeCliExiter,
					nodeProvider: new(nodeprovider.Fake),
					os:           fakeVOS,
					ioutil:       new(vioutil.Fake),
				}

				/* act */
				objectUnderTest.Run(providedPkgRef, &RunOpts{})

				/* assert */
				actualPkgBasePath, actualPkgRef := fakePkg.ResolveArgsForCall(0)
				Expect(actualPkgBasePath).Should(Equal(expectedPkgBasePath))
				Expect(actualPkgRef).Should(Equal(expectedPkgRef))
			})
			Context("pkg.Resolve fails", func() {

			})
			Context("pkg.Resolve succeeds", func() {
				It("should call pkg.Get w/ expected args", func() {
					/* arrange */
					resolvedPkgRef := "dummyPkgName"
					wdReturnedFromVOS := "dummyWorkDir"

					expectedPkgRef := resolvedPkgRef

					fakePkg := new(pkg.Fake)
					// err to trigger immediate return
					fakePkg.GetReturns(&model.PkgManifest{}, errors.New("dummyError"))

					fakeConsumeNodeApi := new(consumenodeapi.Fake)

					fakeCliExiter := new(cliexiter.Fake)

					fakeVOS := new(vos.Fake)
					fakePkg.ResolveReturns(resolvedPkgRef, true)
					fakeVOS.GetwdReturns(wdReturnedFromVOS, nil)

					objectUnderTest := _core{
						pkg:               fakePkg,
						consumeNodeApi:    fakeConsumeNodeApi,
						cliExiter:         fakeCliExiter,
						cliParamSatisfier: new(cliparamsatisfier.Fake),
						nodeProvider:      new(nodeprovider.Fake),
						os:                fakeVOS,
						ioutil:            new(vioutil.Fake),
					}

					/* act */
					objectUnderTest.Run("", &RunOpts{})

					/* assert */
					Expect(fakePkg.GetArgsForCall(0)).Should(Equal(expectedPkgRef))
				})
				Context("pkg.Get errors", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						fakeCliExiter := new(cliexiter.Fake)
						returnedError := errors.New("dummyError")

						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns("", true)
						fakePkg.GetReturns(&model.PkgManifest{}, returnedError)

						objectUnderTest := _core{
							pkg:               fakePkg,
							cliExiter:         fakeCliExiter,
							cliParamSatisfier: new(cliparamsatisfier.Fake),
							nodeProvider:      new(nodeprovider.Fake),
							os:                new(vos.Fake),
							ioutil:            new(vioutil.Fake),
						}

						/* act */
						objectUnderTest.Run("", &RunOpts{})

						/* assert */
						Expect(fakeCliExiter.ExitArgsForCall(0)).
							Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
					})
				})
				Context("pkg.Get doesn't error", func() {
					It("should call paramSatisfier.Satisfy w/ expected args", func() {
						/* arrange */
						param1Name := "DUMMY_PARAM1_NAME"
						pkgManifest := &model.PkgManifest{
							Inputs: map[string]*model.Param{
								param1Name: {
									String: &model.StringParam{},
								},
							},
						}

						expectedParams := pkgManifest.Inputs

						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns("", true)
						fakePkg.GetReturns(
							pkgManifest,
							nil,
						)

						// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
						fakeConsumeNodeApi := new(consumenodeapi.Fake)
						eventChannel := make(chan model.Event)
						close(eventChannel)
						fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)

						objectUnderTest := _core{
							pkg:               fakePkg,
							consumeNodeApi:    fakeConsumeNodeApi,
							cliExiter:         new(cliexiter.Fake),
							cliParamSatisfier: fakeCliParamSatisfier,
							nodeProvider:      new(nodeprovider.Fake),
							os:                new(vos.Fake),
							ioutil:            new(vioutil.Fake),
						}

						/* act */
						objectUnderTest.Run("", &RunOpts{})

						/* assert */
						_, actualParams := fakeCliParamSatisfier.SatisfyArgsForCall(0)

						Expect(actualParams).To(Equal(expectedParams))
					})
					It("should call consumeNodeApi.StartOp w/ expected args", func() {
						/* arrange */
						cwd := "dummyWorkDir"
						fakeVOS := new(vos.Fake)
						fakeVOS.GetwdReturns(cwd, nil)

						resolvedPkgRef := "dummyPkgRef"

						expectedArg1ValueString := "dummyArg1Value"
						expectedArgs := model.StartOpReq{
							Args: map[string]*model.Data{
								"dummyArg1Name": {String: &expectedArg1ValueString},
							},
							PkgRef: resolvedPkgRef,
						}

						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns(resolvedPkgRef, true)
						fakePkg.GetReturns(&model.PkgManifest{}, nil)

						// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
						fakeConsumeNodeApi := new(consumenodeapi.Fake)
						eventChannel := make(chan model.Event)
						close(eventChannel)
						fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args)

						objectUnderTest := _core{
							pkg:               fakePkg,
							consumeNodeApi:    fakeConsumeNodeApi,
							cliExiter:         new(cliexiter.Fake),
							cliParamSatisfier: fakeCliParamSatisfier,
							nodeProvider:      new(nodeprovider.Fake),
							os:                fakeVOS,
							ioutil:            new(vioutil.Fake),
						}

						/* act */
						objectUnderTest.Run("", &RunOpts{})

						/* assert */
						actualArgs := fakeConsumeNodeApi.StartOpArgsForCall(0)
						Expect(actualArgs).To(Equal(expectedArgs))
					})
					Context("consumeNodeApi.StartOp errors", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeCliExiter := new(cliexiter.Fake)
							returnedError := errors.New("dummyError")

							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns("", true)
							fakePkg.GetReturns(&model.PkgManifest{}, nil)

							fakeConsumeNodeApi := new(consumenodeapi.Fake)
							fakeConsumeNodeApi.StartOpReturns("dummyOpId", returnedError)

							objectUnderTest := _core{
								pkg:               fakePkg,
								consumeNodeApi:    fakeConsumeNodeApi,
								cliExiter:         fakeCliExiter,
								cliParamSatisfier: new(cliparamsatisfier.Fake),
								nodeProvider:      new(nodeprovider.Fake),
								os:                new(vos.Fake),
								ioutil:            new(vioutil.Fake),
							}

							/* act */
							objectUnderTest.Run("", &RunOpts{})

							/* assert */
							Expect(fakeCliExiter.ExitArgsForCall(0)).
								Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
						})
					})
					Context("consumeNodeApi.StartOp doesn't error", func() {
						It("should call consumeNodeApi.GetEventStream w/ expected args", func() {
							/* arrange */
							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns("", true)
							fakePkg.GetReturns(&model.PkgManifest{}, nil)
							rootOpIdReturnedFromStartOp := "dummyRootOpId"
							startTime := time.Now().UTC()
							expectedReq := &model.GetEventStreamReq{
								Filter: &model.EventFilter{
									RootOpIds: []string{rootOpIdReturnedFromStartOp},
									Since:     &startTime,
								},
							}

							fakeConsumeNodeApi := new(consumenodeapi.Fake)
							fakeConsumeNodeApi.StartOpReturns(rootOpIdReturnedFromStartOp, nil)
							eventChannel := make(chan model.Event)
							close(eventChannel)
							fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

							objectUnderTest := _core{
								pkg:               fakePkg,
								consumeNodeApi:    fakeConsumeNodeApi,
								cliExiter:         new(cliexiter.Fake),
								cliParamSatisfier: new(cliparamsatisfier.Fake),
								nodeProvider:      new(nodeprovider.Fake),
								os:                new(vos.Fake),
								ioutil:            new(vioutil.Fake),
							}

							/* act */
							objectUnderTest.Run("", &RunOpts{})

							/* assert */
							actualReq := fakeConsumeNodeApi.GetEventStreamArgsForCall(0)

							// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
							Expect(*actualReq.Filter.Since).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
							// set temporal fields to expected vals since they're already asserted
							actualReq.Filter.Since = &startTime

							Expect(actualReq).Should(Equal(expectedReq))
						})
						Context("consumeNodeApi.GetEventStream errors", func() {
							It("should call exiter w/ expected args", func() {
								/* arrange */
								fakeCliExiter := new(cliexiter.Fake)
								returnedError := errors.New("dummyError")

								fakePkg := new(pkg.Fake)
								fakePkg.ResolveReturns("", true)
								fakePkg.GetReturns(&model.PkgManifest{}, nil)

								fakeConsumeNodeApi := new(consumenodeapi.Fake)
								fakeConsumeNodeApi.GetEventStreamReturns(nil, returnedError)

								objectUnderTest := _core{
									pkg:               fakePkg,
									consumeNodeApi:    fakeConsumeNodeApi,
									cliExiter:         fakeCliExiter,
									cliParamSatisfier: new(cliparamsatisfier.Fake),
									nodeProvider:      new(nodeprovider.Fake),
									os:                new(vos.Fake),
									ioutil:            new(vioutil.Fake),
								}

								/* act */
								objectUnderTest.Run("", &RunOpts{})

								/* assert */
								Expect(fakeCliExiter.ExitArgsForCall(0)).
									Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
							})
						})
						Context("consumeNodeApi.GetEventStream doesn't error", func() {
							Context("event channel closes", func() {
								It("should call exiter w/ expected args", func() {
									/* arrange */
									fakeCliExiter := new(cliexiter.Fake)

									fakePkg := new(pkg.Fake)
									fakePkg.ResolveReturns("", true)
									fakePkg.GetReturns(&model.PkgManifest{}, nil)

									fakeConsumeNodeApi := new(consumenodeapi.Fake)
									eventChannel := make(chan model.Event)
									close(eventChannel)
									fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

									objectUnderTest := _core{
										pkg:               fakePkg,
										consumeNodeApi:    fakeConsumeNodeApi,
										cliExiter:         fakeCliExiter,
										cliParamSatisfier: new(cliparamsatisfier.Fake),
										nodeProvider:      new(nodeprovider.Fake),
										os:                new(vos.Fake),
										ioutil:            new(vioutil.Fake),
									}

									/* act */
									objectUnderTest.Run("", &RunOpts{})

									/* assert */
									Expect(fakeCliExiter.ExitArgsForCall(0)).
										Should(Equal(cliexiter.ExitReq{Message: "Event channel closed unexpectedly", Code: 1}))
								})
							})
							Context("event channel doesn't close", func() {
								Context("event received", func() {
									rootOpId := "dummyRootOpId"
									Context("OpEndedEvent", func() {
										Context("Outcome==SUCCEEDED", func() {
											It("should call exiter w/ expected args", func() {
												/* arrange */
												opEndedEvent := model.Event{
													Timestamp: time.Now(),
													OpEnded: &model.OpEndedEvent{
														OpId:     rootOpId,
														PkgRef:   "dummyPkgRef",
														Outcome:  model.OpOutcomeSucceeded,
														RootOpId: rootOpId,
													},
												}

												fakeCliExiter := new(cliexiter.Fake)

												fakePkg := new(pkg.Fake)
												fakePkg.ResolveReturns("", true)
												fakePkg.GetReturns(&model.PkgManifest{}, nil)

												fakeConsumeNodeApi := new(consumenodeapi.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)
												fakeConsumeNodeApi.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

												objectUnderTest := _core{
													pkg:               fakePkg,
													cliColorer:        clicolorer.New(),
													consumeNodeApi:    fakeConsumeNodeApi,
													cliExiter:         fakeCliExiter,
													cliOutput:         new(clioutput.Fake),
													cliParamSatisfier: new(cliparamsatisfier.Fake),
													nodeProvider:      new(nodeprovider.Fake),
													os:                new(vos.Fake),
													ioutil:            new(vioutil.Fake),
												}

												/* act/assert */
												objectUnderTest.Run("", &RunOpts{})
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
														OpId:     rootOpId,
														PkgRef:   "dummyPkgRef",
														Outcome:  model.OpOutcomeKilled,
														RootOpId: rootOpId,
													},
												}

												fakeCliExiter := new(cliexiter.Fake)

												fakePkg := new(pkg.Fake)
												fakePkg.ResolveReturns("", true)
												fakePkg.GetReturns(&model.PkgManifest{}, nil)

												fakeConsumeNodeApi := new(consumenodeapi.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)
												fakeConsumeNodeApi.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

												objectUnderTest := _core{
													pkg:               fakePkg,
													cliColorer:        clicolorer.New(),
													consumeNodeApi:    fakeConsumeNodeApi,
													cliExiter:         fakeCliExiter,
													cliOutput:         new(clioutput.Fake),
													cliParamSatisfier: new(cliparamsatisfier.Fake),
													nodeProvider:      new(nodeprovider.Fake),
													os:                new(vos.Fake),
													ioutil:            new(vioutil.Fake),
												}

												/* act/assert */
												objectUnderTest.Run("", &RunOpts{})
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
														OpId:     rootOpId,
														PkgRef:   "dummyPkgRef",
														Outcome:  model.OpOutcomeFailed,
														RootOpId: rootOpId,
													},
												}

												fakeCliExiter := new(cliexiter.Fake)

												fakePkg := new(pkg.Fake)
												fakePkg.ResolveReturns("", true)
												fakePkg.GetReturns(&model.PkgManifest{}, nil)

												fakeConsumeNodeApi := new(consumenodeapi.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)
												fakeConsumeNodeApi.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

												objectUnderTest := _core{
													pkg:               fakePkg,
													cliColorer:        clicolorer.New(),
													consumeNodeApi:    fakeConsumeNodeApi,
													cliExiter:         fakeCliExiter,
													cliOutput:         new(clioutput.Fake),
													cliParamSatisfier: new(cliparamsatisfier.Fake),
													nodeProvider:      new(nodeprovider.Fake),
													os:                new(vos.Fake),
													ioutil:            new(vioutil.Fake),
												}

												/* act/assert */
												objectUnderTest.Run("", &RunOpts{})
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
														OpId:     rootOpId,
														PkgRef:   "dummyPkgRef",
														Outcome:  "some unexpected outcome",
														RootOpId: rootOpId,
													},
												}

												fakeCliExiter := new(cliexiter.Fake)

												fakePkg := new(pkg.Fake)
												fakePkg.ResolveReturns("", true)
												fakePkg.GetReturns(&model.PkgManifest{}, nil)

												fakeConsumeNodeApi := new(consumenodeapi.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)
												fakeConsumeNodeApi.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

												objectUnderTest := _core{
													pkg:               fakePkg,
													cliColorer:        clicolorer.New(),
													consumeNodeApi:    fakeConsumeNodeApi,
													cliExiter:         fakeCliExiter,
													cliOutput:         new(clioutput.Fake),
													cliParamSatisfier: new(cliparamsatisfier.Fake),
													nodeProvider:      new(nodeprovider.Fake),
													os:                new(vos.Fake),
													ioutil:            new(vioutil.Fake),
												}

												/* act/assert */
												objectUnderTest.Run("", &RunOpts{})
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
})
