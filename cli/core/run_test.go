package core

import (
	"context"
	"errors"
	"time"

	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/clicolorer"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/clioutput"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/api/client"
	"github.com/opctl/sdk-golang/opspec/opfile"
)

var _ = Context("Run", func() {
	Context("Execute", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedOpRef := "dummyOpRef"

			fakeOpDotYmlGetter := new(dotyml.FakeGetter)
			// err to trigger immediate return
			fakeOpDotYmlGetter.GetReturns(nil, errors.New("dummyError"))

			fakeOpHandle := new(data.FakeHandle)
			fakeDataResolver := new(fakeDataResolver)
			fakeDataResolver.ResolveReturns(fakeOpHandle)

			objectUnderTest := _core{
				opDotYmlGetter:    fakeOpDotYmlGetter,
				dataResolver:      fakeDataResolver,
				cliExiter:         new(cliexiter.Fake),
				cliParamSatisfier: new(cliparamsatisfier.Fake),
			}

			/* act */
			objectUnderTest.Run(context.TODO(), providedOpRef, &RunOpts{})

			/* assert */
			actualOpRef, actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)
			Expect(actualOpRef).To(Equal(providedOpRef))
			Expect(actualPullCreds).To(BeNil())
		})
		It("should call data.Get w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()

			fakeOpDotYmlGetter := new(dotyml.FakeGetter)
			// error to trigger immediate return
			fakeOpDotYmlGetter.GetReturns(nil, errors.New("dummyError"))

			fakeOpHandle := new(data.FakeHandle)
			fakeDataResolver := new(fakeDataResolver)
			fakeDataResolver.ResolveReturns(fakeOpHandle)

			objectUnderTest := _core{
				opDotYmlGetter:      fakeOpDotYmlGetter,
				dataResolver:        fakeDataResolver,
				opspecNodeAPIClient: new(client.Fake),
				cliExiter:           new(cliexiter.Fake),
				cliParamSatisfier:   new(cliparamsatisfier.Fake),
				os:                  new(ios.Fake),
				ioutil:              new(iioutil.Fake),
			}

			/* act */
			objectUnderTest.Run(
				providedCtx,
				"",
				new(RunOpts),
			)

			/* assert */
			actualCtx,
				actualOpHandle := fakeOpDotYmlGetter.GetArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualOpHandle).To(Equal(fakeOpHandle))
		})
		Context("data.Get errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				getManifestErr := errors.New("dummyError")

				fakeOpDotYmlGetter := new(dotyml.FakeGetter)
				fakeOpDotYmlGetter.GetReturns(nil, getManifestErr)

				fakeOpHandle := new(data.FakeHandle)
				fakeDataResolver := new(fakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					opDotYmlGetter:    fakeOpDotYmlGetter,
					dataResolver:      fakeDataResolver,
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
		Context("data.Get doesn't error", func() {
			It("should call paramSatisfier.Satisfy w/ expected args", func() {
				/* arrange */
				param1Name := "DUMMY_PARAM1_NAME"
				opDotYml := &model.OpDotYml{
					Inputs: map[string]*model.Param{
						param1Name: {
							String: &model.StringParam{},
						},
					},
				}

				expectedParams := opDotYml.Inputs

				fakeOpDotYmlGetter := new(dotyml.FakeGetter)
				fakeOpDotYmlGetter.GetReturns(opDotYml, nil)

				fakeOpHandle := new(data.FakeHandle)
				fakeDataResolver := new(fakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
				fakeOpspecNodeAPIClient := new(client.Fake)
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

				fakeCliParamSatisfier := new(cliparamsatisfier.Fake)

				objectUnderTest := _core{
					opDotYmlGetter:      fakeOpDotYmlGetter,
					dataResolver:        fakeDataResolver,
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
				resolvedOpRef := "dummyOpRef"

				providedContext := context.TODO()
				expectedCtx := providedContext

				expectedArg1ValueString := "dummyArg1Value"
				expectedArgs := model.StartOpReq{
					Args: map[string]*model.Value{
						"dummyArg1Name": {String: &expectedArg1ValueString},
					},
					Op: model.StartOpReqOp{
						Ref: resolvedOpRef,
					},
				}

				fakeOpDotYmlGetter := new(dotyml.FakeGetter)
				fakeOpDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

				fakeOpHandle := new(data.FakeHandle)
				fakeOpHandle.RefReturns(resolvedOpRef)

				fakeDataResolver := new(fakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
				fakeOpspecNodeAPIClient := new(client.Fake)
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

				fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
				fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args)

				objectUnderTest := _core{
					opDotYmlGetter:      fakeOpDotYmlGetter,
					dataResolver:        fakeDataResolver,
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

					fakeOpDotYmlGetter := new(dotyml.FakeGetter)
					fakeOpDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

					fakeOpHandle := new(data.FakeHandle)
					fakeDataResolver := new(fakeDataResolver)
					fakeDataResolver.ResolveReturns(fakeOpHandle)

					fakeOpspecNodeAPIClient := new(client.Fake)
					fakeOpspecNodeAPIClient.StartOpReturns("dummyOpID", returnedError)

					objectUnderTest := _core{
						opDotYmlGetter:      fakeOpDotYmlGetter,
						dataResolver:        fakeDataResolver,
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
					rootOpIDReturnedFromStartOp := "dummyRootOpID"
					startTime := time.Now().UTC()
					expectedReq := &model.GetEventStreamReq{
						Filter: model.EventFilter{
							Roots: []string{rootOpIDReturnedFromStartOp},
							Since: &startTime,
						},
					}

					fakeOpDotYmlGetter := new(dotyml.FakeGetter)
					fakeOpDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

					fakeOpHandle := new(data.FakeHandle)
					fakeDataResolver := new(fakeDataResolver)
					fakeDataResolver.ResolveReturns(fakeOpHandle)

					fakeOpspecNodeAPIClient := new(client.Fake)
					fakeOpspecNodeAPIClient.StartOpReturns(rootOpIDReturnedFromStartOp, nil)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _core{
						opDotYmlGetter:      fakeOpDotYmlGetter,
						dataResolver:        fakeDataResolver,
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

						fakeOpDotYmlGetter := new(dotyml.FakeGetter)
						fakeOpDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

						fakeOpHandle := new(data.FakeHandle)
						fakeDataResolver := new(fakeDataResolver)
						fakeDataResolver.ResolveReturns(fakeOpHandle)

						fakeOpspecNodeAPIClient := new(client.Fake)
						fakeOpspecNodeAPIClient.GetEventStreamReturns(nil, returnedError)

						objectUnderTest := _core{
							opDotYmlGetter:      fakeOpDotYmlGetter,
							dataResolver:        fakeDataResolver,
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

							fakeOpDotYmlGetter := new(dotyml.FakeGetter)
							fakeOpDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

							fakeOpHandle := new(data.FakeHandle)
							fakeDataResolver := new(fakeDataResolver)
							fakeDataResolver.ResolveReturns(fakeOpHandle)

							fakeOpspecNodeAPIClient := new(client.Fake)
							eventChannel := make(chan model.Event)
							close(eventChannel)
							fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

							objectUnderTest := _core{
								opDotYmlGetter:      fakeOpDotYmlGetter,
								dataResolver:        fakeDataResolver,
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
							rootOpID := "dummyRootOpID"
							Context("OpEndedEvent", func() {
								Context("Outcome==SUCCEEDED", func() {
									It("should call exiter w/ expected args", func() {
										/* arrange */
										opEndedEvent := model.Event{
											Timestamp: time.Now(),
											OpEnded: &model.OpEndedEvent{
												OpID:     rootOpID,
												OpRef:    "dummyOpRef",
												Outcome:  model.OpOutcomeSucceeded,
												RootOpID: rootOpID,
											},
										}

										fakeCliExiter := new(cliexiter.Fake)

										fakeOpDotYmlGetter := new(dotyml.FakeGetter)
										fakeOpDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

										fakeOpHandle := new(data.FakeHandle)
										fakeDataResolver := new(fakeDataResolver)
										fakeDataResolver.ResolveReturns(fakeOpHandle)

										fakeOpspecNodeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpID, nil)

										objectUnderTest := _core{
											opDotYmlGetter:      fakeOpDotYmlGetter,
											dataResolver:        fakeDataResolver,
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
												OpID:     rootOpID,
												OpRef:    "dummyOpRef",
												Outcome:  model.OpOutcomeKilled,
												RootOpID: rootOpID,
											},
										}

										fakeCliExiter := new(cliexiter.Fake)

										fakeOpDotYmlGetter := new(dotyml.FakeGetter)
										fakeOpDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

										fakeOpHandle := new(data.FakeHandle)
										fakeDataResolver := new(fakeDataResolver)
										fakeDataResolver.ResolveReturns(fakeOpHandle)

										fakeOpspecNodeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpID, nil)

										objectUnderTest := _core{
											opDotYmlGetter:      fakeOpDotYmlGetter,
											dataResolver:        fakeDataResolver,
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
												OpID:     rootOpID,
												OpRef:    "dummyOpRef",
												Outcome:  model.OpOutcomeFailed,
												RootOpID: rootOpID,
											},
										}

										fakeCliExiter := new(cliexiter.Fake)

										fakeOpDotYmlGetter := new(dotyml.FakeGetter)
										fakeOpDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

										fakeOpHandle := new(data.FakeHandle)
										fakeDataResolver := new(fakeDataResolver)
										fakeDataResolver.ResolveReturns(fakeOpHandle)

										fakeOpspecNodeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpID, nil)

										objectUnderTest := _core{
											opDotYmlGetter:      fakeOpDotYmlGetter,
											dataResolver:        fakeDataResolver,
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
												OpID:     rootOpID,
												OpRef:    "dummyOpRef",
												Outcome:  "some unexpected outcome",
												RootOpID: rootOpID,
											},
										}

										fakeCliExiter := new(cliexiter.Fake)

										fakeOpDotYmlGetter := new(dotyml.FakeGetter)
										fakeOpDotYmlGetter.GetReturns(&model.OpDotYml{}, nil)

										fakeOpHandle := new(data.FakeHandle)
										fakeDataResolver := new(fakeDataResolver)
										fakeDataResolver.ResolveReturns(fakeOpHandle)

										fakeOpspecNodeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeOpspecNodeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpID, nil)

										objectUnderTest := _core{
											opDotYmlGetter:      fakeOpDotYmlGetter,
											dataResolver:        fakeDataResolver,
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
