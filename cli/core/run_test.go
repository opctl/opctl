package core

import (
	"context"
	"errors"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		It("should call pkgResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"

			fakePkg := new(pkg.Fake)
			// err to trigger immediate return
			fakePkg.GetManifestReturns(nil, errors.New("dummyError"))

			fakePkgHandle := new(pkg.FakeHandle)
			fakePkgResolver := new(fakePkgResolver)
			fakePkgResolver.ResolveReturns(fakePkgHandle)

			objectUnderTest := _core{
				pkg:               fakePkg,
				pkgResolver:       fakePkgResolver,
				cliExiter:         new(cliexiter.Fake),
				cliParamSatisfier: new(cliparamsatisfier.Fake),
			}

			/* act */
			objectUnderTest.Run(context.TODO(), providedPkgRef, &RunOpts{})

			/* assert */
			actualPkgRef, actualPullCreds := fakePkgResolver.ResolveArgsForCall(0)
			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPullCreds).To(BeNil())
		})
		It("should call pkg.GetManifest w/ expected args", func() {
			/* arrange */
			fakePkg := new(pkg.Fake)
			// error to trigger immediate return
			fakePkg.GetManifestReturns(nil, errors.New("dummyError"))

			fakePkgHandle := new(pkg.FakeHandle)
			fakePkgResolver := new(fakePkgResolver)
			fakePkgResolver.ResolveReturns(fakePkgHandle)

			objectUnderTest := _core{
				pkg:                 fakePkg,
				pkgResolver:         fakePkgResolver,
				opspecNodeAPIClient: new(client.Fake),
				cliExiter:           new(cliexiter.Fake),
				cliParamSatisfier:   new(cliparamsatisfier.Fake),
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
				getManifestErr := errors.New("dummyError")

				fakePkg := new(pkg.Fake)
				fakePkg.GetManifestReturns(nil, getManifestErr)

				fakePkgHandle := new(pkg.FakeHandle)
				fakePkgResolver := new(fakePkgResolver)
				fakePkgResolver.ResolveReturns(fakePkgHandle)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:               fakePkg,
					pkgResolver:       fakePkgResolver,
					cliExiter:         fakeCliExiter,
					cliParamSatisfier: new(cliparamsatisfier.Fake),
					os:                new(ios.Fake),
				}

				/* act */
				objectUnderTest.Run(context.TODO(), "", &RunOpts{})

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: getManifestErr.Error(), Code: 1}))
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

				fakePkg := new(pkg.Fake)
				fakePkg.GetManifestReturns(pkgManifest, nil)

				fakePkgHandle := new(pkg.FakeHandle)
				fakePkgResolver := new(fakePkgResolver)
				fakePkgResolver.ResolveReturns(fakePkgHandle)

				// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
				fakeOpspecNodeAPIClient := new(client.Fake)
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

				fakeCliParamSatisfier := new(cliparamsatisfier.Fake)

				objectUnderTest := _core{
					pkg:                 fakePkg,
					pkgResolver:         fakePkgResolver,
					opspecNodeAPIClient: fakeOpspecNodeAPIClient,
					cliExiter:           new(cliexiter.Fake),
					cliParamSatisfier:   fakeCliParamSatisfier,
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

				fakePkg := new(pkg.Fake)
				fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

				fakePkgHandle := new(pkg.FakeHandle)
				fakePkgHandle.RefReturns(resolvedPkgRef)

				fakePkgResolver := new(fakePkgResolver)
				fakePkgResolver.ResolveReturns(fakePkgHandle)

				// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
				fakeOpspecNodeAPIClient := new(client.Fake)
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

				fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
				fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args)

				objectUnderTest := _core{
					pkg:                 fakePkg,
					pkgResolver:         fakePkgResolver,
					opspecNodeAPIClient: fakeOpspecNodeAPIClient,
					cliExiter:           new(cliexiter.Fake),
					cliParamSatisfier:   fakeCliParamSatisfier,
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

					fakePkg := new(pkg.Fake)
					fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

					fakePkgHandle := new(pkg.FakeHandle)
					fakePkgResolver := new(fakePkgResolver)
					fakePkgResolver.ResolveReturns(fakePkgHandle)

					fakeOpspecNodeAPIClient := new(client.Fake)
					fakeOpspecNodeAPIClient.StartOpReturns("dummyOpId", returnedError)

					objectUnderTest := _core{
						pkg:                 fakePkg,
						pkgResolver:         fakePkgResolver,
						opspecNodeAPIClient: fakeOpspecNodeAPIClient,
						cliExiter:           fakeCliExiter,
						cliParamSatisfier:   new(cliparamsatisfier.Fake),
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
					rootOpIdReturnedFromStartOp := "dummyRootOpId"
					startTime := time.Now().UTC()
					expectedReq := &model.GetEventStreamReq{
						Filter: model.EventFilter{
							Roots: []string{rootOpIdReturnedFromStartOp},
							Since: &startTime,
						},
					}

					fakePkg := new(pkg.Fake)
					fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

					fakePkgHandle := new(pkg.FakeHandle)
					fakePkgResolver := new(fakePkgResolver)
					fakePkgResolver.ResolveReturns(fakePkgHandle)

					fakeOpspecNodeAPIClient := new(client.Fake)
					fakeOpspecNodeAPIClient.StartOpReturns(rootOpIdReturnedFromStartOp, nil)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _core{
						pkg:                 fakePkg,
						pkgResolver:         fakePkgResolver,
						opspecNodeAPIClient: fakeOpspecNodeAPIClient,
						cliExiter:           new(cliexiter.Fake),
						cliParamSatisfier:   new(cliparamsatisfier.Fake),
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

						fakePkg := new(pkg.Fake)
						fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

						fakePkgHandle := new(pkg.FakeHandle)
						fakePkgResolver := new(fakePkgResolver)
						fakePkgResolver.ResolveReturns(fakePkgHandle)

						fakeOpspecNodeAPIClient := new(client.Fake)
						fakeOpspecNodeAPIClient.GetEventStreamReturns(nil, returnedError)

						objectUnderTest := _core{
							pkg:                 fakePkg,
							pkgResolver:         fakePkgResolver,
							opspecNodeAPIClient: fakeOpspecNodeAPIClient,
							cliExiter:           fakeCliExiter,
							cliParamSatisfier:   new(cliparamsatisfier.Fake),
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

							fakePkg := new(pkg.Fake)
							fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

							fakePkgHandle := new(pkg.FakeHandle)
							fakePkgResolver := new(fakePkgResolver)
							fakePkgResolver.ResolveReturns(fakePkgHandle)

							fakeOpspecNodeAPIClient := new(client.Fake)
							eventChannel := make(chan model.Event)
							close(eventChannel)
							fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

							objectUnderTest := _core{
								pkg:                 fakePkg,
								pkgResolver:         fakePkgResolver,
								opspecNodeAPIClient: fakeOpspecNodeAPIClient,
								cliExiter:           fakeCliExiter,
								cliParamSatisfier:   new(cliparamsatisfier.Fake),
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

										fakePkg := new(pkg.Fake)
										fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

										fakePkgHandle := new(pkg.FakeHandle)
										fakePkgResolver := new(fakePkgResolver)
										fakePkgResolver.ResolveReturns(fakePkgHandle)

										fakeOpspecNodeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

										objectUnderTest := _core{
											pkg:                 fakePkg,
											pkgResolver:         fakePkgResolver,
											cliColorer:          clicolorer.New(),
											opspecNodeAPIClient: fakeOpspecNodeAPIClient,
											cliExiter:           fakeCliExiter,
											cliOutput:           new(clioutput.Fake),
											cliParamSatisfier:   new(cliparamsatisfier.Fake),
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

										fakePkg := new(pkg.Fake)
										fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

										fakePkgHandle := new(pkg.FakeHandle)
										fakePkgResolver := new(fakePkgResolver)
										fakePkgResolver.ResolveReturns(fakePkgHandle)

										fakeOpspecNodeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

										objectUnderTest := _core{
											pkg:                 fakePkg,
											pkgResolver:         fakePkgResolver,
											cliColorer:          clicolorer.New(),
											opspecNodeAPIClient: fakeOpspecNodeAPIClient,
											cliExiter:           fakeCliExiter,
											cliOutput:           new(clioutput.Fake),
											cliParamSatisfier:   new(cliparamsatisfier.Fake),
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

										fakePkg := new(pkg.Fake)
										fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

										fakePkgHandle := new(pkg.FakeHandle)
										fakePkgResolver := new(fakePkgResolver)
										fakePkgResolver.ResolveReturns(fakePkgHandle)

										fakeOpspecNodeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

										objectUnderTest := _core{
											pkg:                 fakePkg,
											pkgResolver:         fakePkgResolver,
											cliColorer:          clicolorer.New(),
											opspecNodeAPIClient: fakeOpspecNodeAPIClient,
											cliExiter:           fakeCliExiter,
											cliOutput:           new(clioutput.Fake),
											cliParamSatisfier:   new(cliparamsatisfier.Fake),
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

										fakePkg := new(pkg.Fake)
										fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

										fakePkgHandle := new(pkg.FakeHandle)
										fakePkgResolver := new(fakePkgResolver)
										fakePkgResolver.ResolveReturns(fakePkgHandle)

										fakeOpspecNodeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpId, nil)

										objectUnderTest := _core{
											pkg:                 fakePkg,
											pkgResolver:         fakePkgResolver,
											cliColorer:          clicolorer.New(),
											opspecNodeAPIClient: fakeOpspecNodeAPIClient,
											cliExiter:           fakeCliExiter,
											cliOutput:           new(clioutput.Fake),
											cliParamSatisfier:   new(cliparamsatisfier.Fake),
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
