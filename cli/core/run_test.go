package core

import (
	"context"
	"errors"
	"fmt"
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
				objectUnderTest.Run(context.TODO(), "dummyName", &RunOpts{})

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("os.Getwd doesn't error", func() {
			It("should call pkg.Resolve w/ expected args", func() {
				/* arrange */
				providedPkgRef := "dummyPkgRef"

				fakePkg := new(pkg.Fake)
				// error to trigger immediate return
				fakePkg.ResolveReturns(nil, errors.New("dummyError"))

				fakeIOS := new(ios.Fake)
				workDir := "dummyWorkDir"
				fakeIOS.GetwdReturns(workDir, nil)

				objectUnderTest := _core{
					pkg:          fakePkg,
					cliExiter:    new(cliexiter.Fake),
					nodeProvider: new(nodeprovider.Fake),
					os:           fakeIOS,
					ioutil:       new(iioutil.Fake),
				}

				/* act */
				objectUnderTest.Run(context.TODO(), providedPkgRef, &RunOpts{})

				/* assert */
				actualPkgRef, actualResolveOpts := fakePkg.ResolveArgsForCall(0)
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualResolveOpts).To(Equal(&pkg.ResolveOpts{BasePath: workDir}))
			})
			Context("pkg.Resolve errs", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					providedPkgRef := "dummyPkgRef"
					wdReturnedFromIOS := "dummyWorkDir"

					fakeIOS := new(ios.Fake)
					fakeIOS.GetwdReturns(wdReturnedFromIOS, nil)

					resolveErr := errors.New("dummyError")

					expectedMsg := fmt.Sprintf(
						"Unable to resolve package '%v' from '%v'; error was: %v",
						providedPkgRef,
						wdReturnedFromIOS,
						resolveErr.Error(),
					)

					fakePkg := new(pkg.Fake)
					fakePkg.ResolveReturns(nil, resolveErr)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:          fakePkg,
						cliExiter:    fakeCliExiter,
						os:           fakeIOS,
						nodeProvider: new(nodeprovider.Fake),
					}

					/* act */
					objectUnderTest.Run(context.TODO(), providedPkgRef, &RunOpts{})

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: expectedMsg, Code: 1}))
				})
			})
			Context("pkg.Resolve doesn't err", func() {
				It("should call pkg.GetManifest w/ expected args", func() {
					/* arrange */
					fakePkgHandle := new(pkg.FakeHandle)
					fakePkg := new(pkg.Fake)
					fakePkg.ResolveReturns(fakePkgHandle, nil)
					// error to trigger immediate return
					fakePkg.GetManifestReturns(nil, errors.New("dummyError"))

					objectUnderTest := _core{
						pkg:                 fakePkg,
						opspecNodeAPIClient: new(client.Fake),
						cliExiter:           new(cliexiter.Fake),
						cliParamSatisfier:   new(cliparamsatisfier.Fake),
						nodeProvider:        new(nodeprovider.Fake),
						os:                  new(ios.Fake),
						ioutil:              new(iioutil.Fake),
					}

					/* act */
					objectUnderTest.Run(context.TODO(), "", &RunOpts{})

					/* assert */
					Expect(fakePkg.GetManifestArgsForCall(0)).To(Equal(fakePkgHandle))
				})
				Context("pkg.GetManifest errors", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						fakeCliExiter := new(cliexiter.Fake)
						returnedError := errors.New("dummyError")

						fakePkgHandle := new(pkg.FakeHandle)
						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns(fakePkgHandle, nil)
						fakePkg.GetManifestReturns(nil, returnedError)

						objectUnderTest := _core{
							pkg:               fakePkg,
							cliExiter:         fakeCliExiter,
							cliParamSatisfier: new(cliparamsatisfier.Fake),
							nodeProvider:      new(nodeprovider.Fake),
							os:                new(ios.Fake),
						}

						/* act */
						objectUnderTest.Run(context.TODO(), "", &RunOpts{})

						/* assert */
						Expect(fakeCliExiter.ExitArgsForCall(0)).
							To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
					})
				})
				Context("pkg.GetManifest doesn't error", func() {
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

						fakePkgHandle := new(pkg.FakeHandle)
						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns(fakePkgHandle, nil)
						fakePkg.GetManifestReturns(pkgManifest, nil)

						// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
						fakeOpspecNodeAPIClient := new(client.Fake)
						eventChannel := make(chan model.Event)
						close(eventChannel)
						fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)

						objectUnderTest := _core{
							pkg:                 fakePkg,
							opspecNodeAPIClient: fakeOpspecNodeAPIClient,
							cliExiter:           new(cliexiter.Fake),
							cliParamSatisfier:   fakeCliParamSatisfier,
							nodeProvider:        new(nodeprovider.Fake),
							os:                  new(ios.Fake),
							ioutil:              new(iioutil.Fake),
						}

						/* act */
						objectUnderTest.Run(context.TODO(), "", &RunOpts{})

						/* assert */
						_, actualParams := fakeCliParamSatisfier.SatisfyArgsForCall(0)

						Expect(actualParams).To(Equal(expectedParams))
					})
					It("should call opspecNodeAPIClient.StartOp w/ expected args", func() {
						/* arrange */
						resolvedPkgRef := "dummyPkgRef"

						providedContext := context.TODO()
						expectedCtx := providedContext

						expectedArg1ValueString := "dummyArg1Value"
						expectedArgs := model.StartOpReq{
							Args: map[string]*model.Value{
								"dummyArg1Name": {String: &expectedArg1ValueString},
							},
							Pkg: &model.DCGOpCallPkg{
								Ref: resolvedPkgRef,
							},
						}

						fakePkgHandle := new(pkg.FakeHandle)
						fakePkgHandle.RefReturns(resolvedPkgRef)
						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns(fakePkgHandle, nil)
						fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

						// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
						fakeOpspecNodeAPIClient := new(client.Fake)
						eventChannel := make(chan model.Event)
						close(eventChannel)
						fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args)

						objectUnderTest := _core{
							pkg:                 fakePkg,
							opspecNodeAPIClient: fakeOpspecNodeAPIClient,
							cliExiter:           new(cliexiter.Fake),
							cliParamSatisfier:   fakeCliParamSatisfier,
							nodeProvider:        new(nodeprovider.Fake),
							os:                  new(ios.Fake),
							ioutil:              new(iioutil.Fake),
						}

						/* act */
						objectUnderTest.Run(providedContext, "", &RunOpts{})

						/* assert */
						actualCtx, actualArgs := fakeOpspecNodeAPIClient.StartOpArgsForCall(0)
						Expect(actualCtx).To(Equal(expectedCtx))
						Expect(actualArgs).To(Equal(expectedArgs))
					})
					Context("opspecNodeAPIClient.StartOp errors", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeCliExiter := new(cliexiter.Fake)
							returnedError := errors.New("dummyError")

							fakePkgHandle := new(pkg.FakeHandle)
							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns(fakePkgHandle, nil)
							fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

							fakeOpspecNodeAPIClient := new(client.Fake)
							fakeOpspecNodeAPIClient.StartOpReturns("dummyOpId", returnedError)

							objectUnderTest := _core{
								pkg:                 fakePkg,
								opspecNodeAPIClient: fakeOpspecNodeAPIClient,
								cliExiter:           fakeCliExiter,
								cliParamSatisfier:   new(cliparamsatisfier.Fake),
								nodeProvider:        new(nodeprovider.Fake),
								os:                  new(ios.Fake),
								ioutil:              new(iioutil.Fake),
							}

							/* act */
							objectUnderTest.Run(context.TODO(), "", &RunOpts{})

							/* assert */
							Expect(fakeCliExiter.ExitArgsForCall(0)).
								To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
						})
					})
					Context("opspecNodeAPIClient.StartOp doesn't error", func() {
						It("should call opspecNodeAPIClient.GetEventStream w/ expected args", func() {
							/* arrange */
							fakePkgHandle := new(pkg.FakeHandle)
							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns(fakePkgHandle, nil)
							fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

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
								pkg:                 fakePkg,
								opspecNodeAPIClient: fakeOpspecNodeAPIClient,
								cliExiter:           new(cliexiter.Fake),
								cliParamSatisfier:   new(cliparamsatisfier.Fake),
								nodeProvider:        new(nodeprovider.Fake),
								os:                  new(ios.Fake),
								ioutil:              new(iioutil.Fake),
							}

							/* act */
							objectUnderTest.Run(context.TODO(), "", &RunOpts{})

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

								fakePkgHandle := new(pkg.FakeHandle)
								fakePkg := new(pkg.Fake)
								fakePkg.ResolveReturns(fakePkgHandle, nil)
								fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

								fakeOpspecNodeAPIClient := new(client.Fake)
								fakeOpspecNodeAPIClient.GetEventStreamReturns(nil, returnedError)

								objectUnderTest := _core{
									pkg:                 fakePkg,
									opspecNodeAPIClient: fakeOpspecNodeAPIClient,
									cliExiter:           fakeCliExiter,
									cliParamSatisfier:   new(cliparamsatisfier.Fake),
									nodeProvider:        new(nodeprovider.Fake),
									os:                  new(ios.Fake),
									ioutil:              new(iioutil.Fake),
								}

								/* act */
								objectUnderTest.Run(context.TODO(), "", &RunOpts{})

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

									fakePkgHandle := new(pkg.FakeHandle)
									fakePkg := new(pkg.Fake)
									fakePkg.ResolveReturns(fakePkgHandle, nil)
									fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

									fakeOpspecNodeAPIClient := new(client.Fake)
									eventChannel := make(chan model.Event)
									close(eventChannel)
									fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

									objectUnderTest := _core{
										pkg:                 fakePkg,
										opspecNodeAPIClient: fakeOpspecNodeAPIClient,
										cliExiter:           fakeCliExiter,
										cliParamSatisfier:   new(cliparamsatisfier.Fake),
										nodeProvider:        new(nodeprovider.Fake),
										os:                  new(ios.Fake),
										ioutil:              new(iioutil.Fake),
									}

									/* act */
									objectUnderTest.Run(context.TODO(), "", &RunOpts{})

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

												fakePkgHandle := new(pkg.FakeHandle)
												fakePkg := new(pkg.Fake)
												fakePkg.ResolveReturns(fakePkgHandle, nil)
												fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

												fakeOpspecNodeAPIClient := new(client.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
												fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

												objectUnderTest := _core{
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
												objectUnderTest.Run(context.TODO(), "", &RunOpts{})
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

												fakePkgHandle := new(pkg.FakeHandle)
												fakePkg := new(pkg.Fake)
												fakePkg.ResolveReturns(fakePkgHandle, nil)
												fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

												fakeOpspecNodeAPIClient := new(client.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
												fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

												objectUnderTest := _core{
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
												objectUnderTest.Run(context.TODO(), "", &RunOpts{})
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

												fakePkgHandle := new(pkg.FakeHandle)
												fakePkg := new(pkg.Fake)
												fakePkg.ResolveReturns(fakePkgHandle, nil)
												fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

												fakeOpspecNodeAPIClient := new(client.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
												fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

												objectUnderTest := _core{
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
												objectUnderTest.Run(context.TODO(), "", &RunOpts{})
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

												fakePkgHandle := new(pkg.FakeHandle)
												fakePkg := new(pkg.Fake)
												fakePkg.ResolveReturns(fakePkgHandle, nil)
												fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

												fakeOpspecNodeAPIClient := new(client.Fake)
												eventChannel := make(chan model.Event, 10)
												eventChannel <- opEndedEvent
												defer close(eventChannel)
												fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
												fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

												objectUnderTest := _core{
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
												objectUnderTest.Run(context.TODO(), "", &RunOpts{})
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
