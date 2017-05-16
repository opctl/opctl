package core

import (
	"errors"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/util/clicolorer"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/clioutput"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api/client"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"path/filepath"
	"time"
)

var _ = Context("Run", func() {
	Context("Execute", func() {
		Context("os.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeIOS := new(ios.Fake)
				expectedError := errors.New("dummyError")
				fakeIOS.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:          new(pkg.Fake),
					cliExiter:    fakeCliExiter,
					nodeProvider: new(nodeprovider.Fake),
					os:           fakeIOS,
				}

				/* act */
				objectUnderTest.Run("dummyName", &RunOpts{})

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("os.Getwd doesn't error", func() {
			It("should call pkg.ParseRef w/ expected args", func() {
				/* arrange */
				providedPkgRef := "dummyPkgName"
				expectedPkgRef := providedPkgRef

				fakePkg := new(pkg.Fake)
				// error to trigger immediate return
				fakePkg.ParseRefReturns(nil, errors.New("dummyError"))

				objectUnderTest := _core{
					pkg:          fakePkg,
					cliExiter:    new(cliexiter.Fake),
					nodeProvider: new(nodeprovider.Fake),
					os:           new(ios.Fake),
				}

				/* act */
				objectUnderTest.Run(providedPkgRef, &RunOpts{})

				/* assert */
				actualPkgRef := fakePkg.ParseRefArgsForCall(0)
				Expect(actualPkgRef).To(Equal(expectedPkgRef))

			})
			Context("pkg.ParseRef errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					expectedError := errors.New("dummyError")

					fakePkg := new(pkg.Fake)
					fakePkg.ParseRefReturns(nil, expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:          fakePkg,
						cliExiter:    fakeCliExiter,
						nodeProvider: new(nodeprovider.Fake),
						os:           new(ios.Fake),
						ioutil:       new(iioutil.Fake),
					}

					/* act */
					objectUnderTest.Run("dummyPkgRef", &RunOpts{})

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
				})
			})
			Context("pkg.ParseRef doesn't error", func() {
				It("should call pkg.Resolve w/ expected args", func() {
					/* arrange */
					providedPkgRef := "dummyPkgName"

					expectedLookPaths := []string{"dummyWorkDir"}
					expectedPkgRef := &pkg.PkgRef{
						FullyQualifiedName: "dummyName",
						Version:            "dummyVersion",
					}

					fakePkg := new(pkg.Fake)
					fakePkg.ParseRefReturns(expectedPkgRef, nil)

					fakeCliExiter := new(cliexiter.Fake)

					fakeIOS := new(ios.Fake)
					fakeIOS.GetwdReturns(expectedLookPaths[0], nil)

					objectUnderTest := _core{
						pkg:          fakePkg,
						cliExiter:    fakeCliExiter,
						nodeProvider: new(nodeprovider.Fake),
						os:           fakeIOS,
						ioutil:       new(iioutil.Fake),
					}

					/* act */
					objectUnderTest.Run(providedPkgRef, &RunOpts{})

					/* assert */
					actualPkgRef, actualLookPaths := fakePkg.ResolveArgsForCall(0)
					Expect(actualLookPaths).To(Equal(expectedLookPaths))
					Expect(actualPkgRef).To(Equal(expectedPkgRef))
				})
				Context("pkg.Resolve fails", func() {

				})
				Context("pkg.Resolve succeeds", func() {
					It("should call manifest.Unmarshal w/ expected args", func() {
						/* arrange */
						pkgPath := "dummyPkgName"
						wdReturnedFromIOS := "dummyWorkDir"

						expectedPkgRef := filepath.Join(pkgPath, pkg.OpDotYmlFileName)

						fakeManifest := new(manifest.Fake)
						// err to trigger immediate return
						fakeManifest.UnmarshalReturns(&model.PkgManifest{}, errors.New("dummyError"))

						fakeOpspecNodeAPIClient := new(client.Fake)

						fakeCliExiter := new(cliexiter.Fake)

						fakeIOS := new(ios.Fake)
						fakeIOS.GetwdReturns(wdReturnedFromIOS, nil)

						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns(pkgPath, true)

						objectUnderTest := _core{
							manifest:            fakeManifest,
							pkg:                 fakePkg,
							opspecNodeAPIClient: fakeOpspecNodeAPIClient,
							cliExiter:           fakeCliExiter,
							cliParamSatisfier:   new(cliparamsatisfier.Fake),
							nodeProvider:        new(nodeprovider.Fake),
							os:                  fakeIOS,
							ioutil:              new(iioutil.Fake),
						}

						/* act */
						objectUnderTest.Run("", &RunOpts{})

						/* assert */
						Expect(fakeManifest.UnmarshalArgsForCall(0)).To(Equal(expectedPkgRef))
					})
					Context("manifest.Unmarshal errors", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeCliExiter := new(cliexiter.Fake)
							returnedError := errors.New("dummyError")

							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns("", true)

							fakeManifest := new(manifest.Fake)
							fakeManifest.UnmarshalReturns(&model.PkgManifest{}, returnedError)

							objectUnderTest := _core{
								manifest:          fakeManifest,
								pkg:               fakePkg,
								cliExiter:         fakeCliExiter,
								cliParamSatisfier: new(cliparamsatisfier.Fake),
								nodeProvider:      new(nodeprovider.Fake),
								os:                new(ios.Fake),
								ioutil:            new(iioutil.Fake),
							}

							/* act */
							objectUnderTest.Run("", &RunOpts{})

							/* assert */
							Expect(fakeCliExiter.ExitArgsForCall(0)).
								To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
						})
					})
					Context("manifest.Unmarshal doesn't error", func() {
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

							fakeManifest := new(manifest.Fake)
							fakeManifest.UnmarshalReturns(
								pkgManifest,
								nil,
							)

							// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
							fakeOpspecNodeAPIClient := new(client.Fake)
							eventChannel := make(chan model.Event)
							close(eventChannel)
							fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

							fakeCliParamSatisfier := new(cliparamsatisfier.Fake)

							objectUnderTest := _core{
								manifest:            fakeManifest,
								pkg:                 fakePkg,
								opspecNodeAPIClient: fakeOpspecNodeAPIClient,
								cliExiter:           new(cliexiter.Fake),
								cliParamSatisfier:   fakeCliParamSatisfier,
								nodeProvider:        new(nodeprovider.Fake),
								os:                  new(ios.Fake),
								ioutil:              new(iioutil.Fake),
							}

							/* act */
							objectUnderTest.Run("", &RunOpts{})

							/* assert */
							_, actualParams := fakeCliParamSatisfier.SatisfyArgsForCall(0)

							Expect(actualParams).To(Equal(expectedParams))
						})
						It("should call opspecNodeAPIClient.StartOp w/ expected args", func() {
							/* arrange */
							cwd := "dummyWorkDir"
							fakeIOS := new(ios.Fake)
							fakeIOS.GetwdReturns(cwd, nil)

							resolvedPkgRef := "dummyPkgRef"

							expectedArg1ValueString := "dummyArg1Value"
							expectedArgs := model.StartOpReq{
								Args: map[string]*model.Data{
									"dummyArg1Name": {String: &expectedArg1ValueString},
								},
								Pkg: &model.DCGOpCallPkg{
									Ref: resolvedPkgRef,
								},
							}

							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns(resolvedPkgRef, true)

							fakeManifest := new(manifest.Fake)
							fakeManifest.UnmarshalReturns(&model.PkgManifest{}, nil)

							// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
							fakeOpspecNodeAPIClient := new(client.Fake)
							eventChannel := make(chan model.Event)
							close(eventChannel)
							fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

							fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
							fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args)

							objectUnderTest := _core{
								manifest:            fakeManifest,
								pkg:                 fakePkg,
								opspecNodeAPIClient: fakeOpspecNodeAPIClient,
								cliExiter:           new(cliexiter.Fake),
								cliParamSatisfier:   fakeCliParamSatisfier,
								nodeProvider:        new(nodeprovider.Fake),
								os:                  fakeIOS,
								ioutil:              new(iioutil.Fake),
							}

							/* act */
							objectUnderTest.Run("", &RunOpts{})

							/* assert */
							actualArgs := fakeOpspecNodeAPIClient.StartOpArgsForCall(0)
							Expect(actualArgs).To(Equal(expectedArgs))
						})
						Context("opspecNodeAPIClient.StartOp errors", func() {
							It("should call exiter w/ expected args", func() {
								/* arrange */
								fakeCliExiter := new(cliexiter.Fake)
								returnedError := errors.New("dummyError")

								fakePkg := new(pkg.Fake)
								fakePkg.ResolveReturns("", true)

								fakeManifest := new(manifest.Fake)
								fakeManifest.UnmarshalReturns(&model.PkgManifest{}, nil)

								fakeOpspecNodeAPIClient := new(client.Fake)
								fakeOpspecNodeAPIClient.StartOpReturns("dummyOpId", returnedError)

								objectUnderTest := _core{
									manifest:            fakeManifest,
									pkg:                 fakePkg,
									opspecNodeAPIClient: fakeOpspecNodeAPIClient,
									cliExiter:           fakeCliExiter,
									cliParamSatisfier:   new(cliparamsatisfier.Fake),
									nodeProvider:        new(nodeprovider.Fake),
									os:                  new(ios.Fake),
									ioutil:              new(iioutil.Fake),
								}

								/* act */
								objectUnderTest.Run("", &RunOpts{})

								/* assert */
								Expect(fakeCliExiter.ExitArgsForCall(0)).
									To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
							})
						})
						Context("opspecNodeAPIClient.StartOp doesn't error", func() {
							It("should call opspecNodeAPIClient.GetEventStream w/ expected args", func() {
								/* arrange */
								fakePkg := new(pkg.Fake)
								fakePkg.ResolveReturns("", true)

								fakeManifest := new(manifest.Fake)
								fakeManifest.UnmarshalReturns(&model.PkgManifest{}, nil)

								rootOpIdReturnedFromStartOp := "dummyRootOpId"
								startTime := time.Now().UTC()
								expectedReq := &model.GetEventStreamReq{
									Filter: &model.EventFilter{
										RootOpIds: []string{rootOpIdReturnedFromStartOp},
										Since:     &startTime,
									},
								}

								fakeOpspecNodeAPIClient := new(client.Fake)
								fakeOpspecNodeAPIClient.StartOpReturns(rootOpIdReturnedFromStartOp, nil)
								eventChannel := make(chan model.Event)
								close(eventChannel)
								fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

								objectUnderTest := _core{
									manifest:            fakeManifest,
									pkg:                 fakePkg,
									opspecNodeAPIClient: fakeOpspecNodeAPIClient,
									cliExiter:           new(cliexiter.Fake),
									cliParamSatisfier:   new(cliparamsatisfier.Fake),
									nodeProvider:        new(nodeprovider.Fake),
									os:                  new(ios.Fake),
									ioutil:              new(iioutil.Fake),
								}

								/* act */
								objectUnderTest.Run("", &RunOpts{})

								/* assert */
								actualReq := fakeOpspecNodeAPIClient.GetEventStreamArgsForCall(0)

								// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
								Expect(*actualReq.Filter.Since).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
								// set temporal fields to expected vals since they're already asserted
								actualReq.Filter.Since = &startTime

								Expect(actualReq).To(Equal(expectedReq))
							})
							Context("opspecNodeAPIClient.GetEventStream errors", func() {
								It("should call exiter w/ expected args", func() {
									/* arrange */
									fakeCliExiter := new(cliexiter.Fake)
									returnedError := errors.New("dummyError")

									fakePkg := new(pkg.Fake)
									fakePkg.ResolveReturns("", true)

									fakeManifest := new(manifest.Fake)
									fakeManifest.UnmarshalReturns(&model.PkgManifest{}, nil)

									fakeOpspecNodeAPIClient := new(client.Fake)
									fakeOpspecNodeAPIClient.GetEventStreamReturns(nil, returnedError)

									objectUnderTest := _core{
										manifest:            fakeManifest,
										pkg:                 fakePkg,
										opspecNodeAPIClient: fakeOpspecNodeAPIClient,
										cliExiter:           fakeCliExiter,
										cliParamSatisfier:   new(cliparamsatisfier.Fake),
										nodeProvider:        new(nodeprovider.Fake),
										os:                  new(ios.Fake),
										ioutil:              new(iioutil.Fake),
									}

									/* act */
									objectUnderTest.Run("", &RunOpts{})

									/* assert */
									Expect(fakeCliExiter.ExitArgsForCall(0)).
										To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
								})
							})
							Context("opspecNodeAPIClient.GetEventStream doesn't error", func() {
								Context("event channel closes", func() {
									It("should call exiter w/ expected args", func() {
										/* arrange */
										fakeCliExiter := new(cliexiter.Fake)

										fakePkg := new(pkg.Fake)
										fakePkg.ResolveReturns("", true)

										fakeManifest := new(manifest.Fake)
										fakeManifest.UnmarshalReturns(&model.PkgManifest{}, nil)

										fakeOpspecNodeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event)
										close(eventChannel)
										fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

										objectUnderTest := _core{
											manifest:            fakeManifest,
											pkg:                 fakePkg,
											opspecNodeAPIClient: fakeOpspecNodeAPIClient,
											cliExiter:           fakeCliExiter,
											cliParamSatisfier:   new(cliparamsatisfier.Fake),
											nodeProvider:        new(nodeprovider.Fake),
											os:                  new(ios.Fake),
											ioutil:              new(iioutil.Fake),
										}

										/* act */
										objectUnderTest.Run("", &RunOpts{})

										/* assert */
										Expect(fakeCliExiter.ExitArgsForCall(0)).
											To(Equal(cliexiter.ExitReq{Message: "Event channel closed unexpectedly", Code: 1}))
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

													fakeManifest := new(manifest.Fake)
													fakeManifest.UnmarshalReturns(&model.PkgManifest{}, nil)

													fakeOpspecNodeAPIClient := new(client.Fake)
													eventChannel := make(chan model.Event, 10)
													eventChannel <- opEndedEvent
													defer close(eventChannel)
													fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
													fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

													objectUnderTest := _core{
														manifest:            fakeManifest,
														pkg:                 fakePkg,
														cliColorer:          clicolorer.New(),
														opspecNodeAPIClient: fakeOpspecNodeAPIClient,
														cliExiter:           fakeCliExiter,
														cliOutput:           new(clioutput.Fake),
														cliParamSatisfier:   new(cliparamsatisfier.Fake),
														nodeProvider:        new(nodeprovider.Fake),
														os:                  new(ios.Fake),
														ioutil:              new(iioutil.Fake),
													}

													/* act/assert */
													objectUnderTest.Run("", &RunOpts{})
													Expect(fakeCliExiter.ExitArgsForCall(0)).
														To(Equal(cliexiter.ExitReq{Code: 0}))
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

													fakeManifest := new(manifest.Fake)
													fakeManifest.UnmarshalReturns(&model.PkgManifest{}, nil)

													fakeOpspecNodeAPIClient := new(client.Fake)
													eventChannel := make(chan model.Event, 10)
													eventChannel <- opEndedEvent
													defer close(eventChannel)
													fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
													fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

													objectUnderTest := _core{
														manifest:            fakeManifest,
														pkg:                 fakePkg,
														cliColorer:          clicolorer.New(),
														opspecNodeAPIClient: fakeOpspecNodeAPIClient,
														cliExiter:           fakeCliExiter,
														cliOutput:           new(clioutput.Fake),
														cliParamSatisfier:   new(cliparamsatisfier.Fake),
														nodeProvider:        new(nodeprovider.Fake),
														os:                  new(ios.Fake),
														ioutil:              new(iioutil.Fake),
													}

													/* act/assert */
													objectUnderTest.Run("", &RunOpts{})
													Expect(fakeCliExiter.ExitArgsForCall(0)).
														To(Equal(cliexiter.ExitReq{Code: 137}))
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

													fakeManifest := new(manifest.Fake)
													fakeManifest.UnmarshalReturns(&model.PkgManifest{}, nil)

													fakeOpspecNodeAPIClient := new(client.Fake)
													eventChannel := make(chan model.Event, 10)
													eventChannel <- opEndedEvent
													defer close(eventChannel)
													fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
													fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

													objectUnderTest := _core{
														manifest:            fakeManifest,
														pkg:                 fakePkg,
														cliColorer:          clicolorer.New(),
														opspecNodeAPIClient: fakeOpspecNodeAPIClient,
														cliExiter:           fakeCliExiter,
														cliOutput:           new(clioutput.Fake),
														cliParamSatisfier:   new(cliparamsatisfier.Fake),
														nodeProvider:        new(nodeprovider.Fake),
														os:                  new(ios.Fake),
														ioutil:              new(iioutil.Fake),
													}

													/* act/assert */
													objectUnderTest.Run("", &RunOpts{})
													Expect(fakeCliExiter.ExitArgsForCall(0)).
														To(Equal(cliexiter.ExitReq{Code: 1}))
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

													fakeManifest := new(manifest.Fake)
													fakeManifest.UnmarshalReturns(&model.PkgManifest{}, nil)

													fakeOpspecNodeAPIClient := new(client.Fake)
													eventChannel := make(chan model.Event, 10)
													eventChannel <- opEndedEvent
													defer close(eventChannel)
													fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
													fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

													objectUnderTest := _core{
														manifest:            fakeManifest,
														pkg:                 fakePkg,
														cliColorer:          clicolorer.New(),
														opspecNodeAPIClient: fakeOpspecNodeAPIClient,
														cliExiter:           fakeCliExiter,
														cliOutput:           new(clioutput.Fake),
														cliParamSatisfier:   new(cliparamsatisfier.Fake),
														nodeProvider:        new(nodeprovider.Fake),
														os:                  new(ios.Fake),
														ioutil:              new(iioutil.Fake),
													}

													/* act/assert */
													objectUnderTest.Run("", &RunOpts{})
													Expect(fakeCliExiter.ExitArgsForCall(0)).
														To(Equal(cliexiter.ExitReq{Code: 1}))
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
})
