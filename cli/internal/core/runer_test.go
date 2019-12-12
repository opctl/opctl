package core

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	cliModel "github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

var _ = Context("Runer", func() {
	Context("Run", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedOpRef := "dummyOpRef"

			fakeOpFileGetter := new(opfile.FakeGetter)
			// err to trigger immediate return
			fakeOpFileGetter.GetReturns(nil, errors.New("dummyError"))

			fakeOpHandle := new(data.FakeHandle)
			fakeDataResolver := new(dataresolver.Fake)
			fakeDataResolver.ResolveReturns(fakeOpHandle)

			objectUnderTest := _runer{
				opFileGetter:      fakeOpFileGetter,
				dataResolver:      fakeDataResolver,
				cliExiter:         new(cliexiter.Fake),
				cliParamSatisfier: new(cliparamsatisfier.Fake),
			}

			/* act */
			objectUnderTest.Run(context.TODO(), providedOpRef, &cliModel.RunOpts{})

			/* assert */
			actualOpRef, actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)
			Expect(actualOpRef).To(Equal(providedOpRef))
			Expect(actualPullCreds).To(BeNil())
		})
		It("should call data.Get w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()

			fakeOpFileGetter := new(opfile.FakeGetter)
			// error to trigger immediate return
			fakeOpFileGetter.GetReturns(nil, errors.New("dummyError"))

			fakeOpHandle := new(data.FakeHandle)
			fakeDataResolver := new(dataresolver.Fake)
			fakeDataResolver.ResolveReturns(fakeOpHandle)

			objectUnderTest := _runer{
				opFileGetter:      fakeOpFileGetter,
				dataResolver:      fakeDataResolver,
				apiClient:         new(client.Fake),
				cliExiter:         new(cliexiter.Fake),
				cliParamSatisfier: new(cliparamsatisfier.Fake),
			}

			/* act */
			objectUnderTest.Run(
				providedCtx,
				"",
				new(cliModel.RunOpts),
			)

			/* assert */
			actualCtx,
				actualOpHandle := fakeOpFileGetter.GetArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualOpHandle).To(Equal(fakeOpHandle))
		})
		Context("data.Get errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				getManifestErr := errors.New("dummyError")

				fakeOpFileGetter := new(opfile.FakeGetter)
				fakeOpFileGetter.GetReturns(nil, getManifestErr)

				fakeOpHandle := new(data.FakeHandle)
				fakeDataResolver := new(dataresolver.Fake)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _runer{
					opFileGetter:      fakeOpFileGetter,
					dataResolver:      fakeDataResolver,
					cliExiter:         fakeCliExiter,
					cliParamSatisfier: new(cliparamsatisfier.Fake),
				}

				/* act */
				objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: getManifestErr.Error(), Code: 1}))
			})
		})
		Context("data.Get doesn't error", func() {
			It("should call paramSatisfier.Satisfy w/ expected args", func() {
				/* arrange */
				param1Name := "DUMMY_PARAM1_NAME"
				opFile := &model.OpFile{
					Inputs: map[string]*model.Param{
						param1Name: {
							String: &model.StringParam{},
						},
					},
				}

				expectedParams := opFile.Inputs

				fakeOpFileGetter := new(opfile.FakeGetter)
				fakeOpFileGetter.GetReturns(opFile, nil)

				fakeOpHandle := new(data.FakeHandle)
				fakeDataResolver := new(dataresolver.Fake)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
				fakeAPIClient := new(client.Fake)
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

				fakeCliParamSatisfier := new(cliparamsatisfier.Fake)

				objectUnderTest := _runer{
					opFileGetter:      fakeOpFileGetter,
					dataResolver:      fakeDataResolver,
					apiClient:         fakeAPIClient,
					cliExiter:         new(cliexiter.Fake),
					cliParamSatisfier: fakeCliParamSatisfier,
				}

				/* act */
				objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

				/* assert */
				_, actualParams := fakeCliParamSatisfier.SatisfyArgsForCall(0)

				Expect(actualParams).To(Equal(expectedParams))
			})
			It("should call apiClient.StartOp w/ expected args", func() {
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

				fakeOpFileGetter := new(opfile.FakeGetter)
				fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

				fakeOpHandle := new(data.FakeHandle)
				fakeOpHandle.RefReturns(resolvedOpRef)

				fakeDataResolver := new(dataresolver.Fake)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
				fakeAPIClient := new(client.Fake)
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

				fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
				fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args)

				objectUnderTest := _runer{
					opFileGetter:      fakeOpFileGetter,
					dataResolver:      fakeDataResolver,
					apiClient:         fakeAPIClient,
					cliExiter:         new(cliexiter.Fake),
					cliParamSatisfier: fakeCliParamSatisfier,
				}

				/* act */
				objectUnderTest.Run(providedContext, "", &cliModel.RunOpts{})

				/* assert */
				actualCtx, actualArgs := fakeAPIClient.StartOpArgsForCall(0)
				Expect(actualCtx).To(Equal(expectedCtx))
				Expect(actualArgs).To(Equal(expectedArgs))
			})
			Context("apiClient.StartOp errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiter.Fake)
					returnedError := errors.New("dummyError")

					fakeOpFileGetter := new(opfile.FakeGetter)
					fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

					fakeOpHandle := new(data.FakeHandle)
					fakeDataResolver := new(dataresolver.Fake)
					fakeDataResolver.ResolveReturns(fakeOpHandle)

					fakeAPIClient := new(client.Fake)
					fakeAPIClient.StartOpReturns("dummyOpID", returnedError)

					objectUnderTest := _runer{
						opFileGetter:      fakeOpFileGetter,
						dataResolver:      fakeDataResolver,
						apiClient:         fakeAPIClient,
						cliExiter:         fakeCliExiter,
						cliParamSatisfier: new(cliparamsatisfier.Fake),
					}

					/* act */
					objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
				})
			})
			Context("apiClient.StartOp doesn't error", func() {
				It("should call apiClient.GetEventStream w/ expected args", func() {
					/* arrange */
					providedCtx := context.Background()
					rootOpIDReturnedFromStartOp := "dummyRootOpID"
					startTime := time.Now().UTC()
					expectedReq := &model.GetEventStreamReq{
						Filter: model.EventFilter{
							Roots: []string{rootOpIDReturnedFromStartOp},
							Since: &startTime,
						},
					}

					fakeOpFileGetter := new(opfile.FakeGetter)
					fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

					fakeOpHandle := new(data.FakeHandle)
					fakeDataResolver := new(dataresolver.Fake)
					fakeDataResolver.ResolveReturns(fakeOpHandle)

					fakeAPIClient := new(client.Fake)
					fakeAPIClient.StartOpReturns(rootOpIDReturnedFromStartOp, nil)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _runer{
						opFileGetter:      fakeOpFileGetter,
						dataResolver:      fakeDataResolver,
						apiClient:         fakeAPIClient,
						cliExiter:         new(cliexiter.Fake),
						cliParamSatisfier: new(cliparamsatisfier.Fake),
					}

					/* act */
					objectUnderTest.Run(providedCtx, "", &cliModel.RunOpts{})

					/* assert */
					actualCtx,
						actualReq := fakeAPIClient.GetEventStreamArgsForCall(0)

					// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
					Expect(*actualReq.Filter.Since).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
					// set temporal fields to expected vals since they're already asserted
					actualReq.Filter.Since = &startTime

					Expect(actualCtx).To(Equal(providedCtx))
					Expect(actualReq).To(Equal(expectedReq))
				})
				Context("apiClient.GetEventStream errors", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						fakeCliExiter := new(cliexiter.Fake)
						returnedError := errors.New("dummyError")

						fakeOpFileGetter := new(opfile.FakeGetter)
						fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

						fakeOpHandle := new(data.FakeHandle)
						fakeDataResolver := new(dataresolver.Fake)
						fakeDataResolver.ResolveReturns(fakeOpHandle)

						fakeAPIClient := new(client.Fake)
						fakeAPIClient.GetEventStreamReturns(nil, returnedError)

						objectUnderTest := _runer{
							opFileGetter:      fakeOpFileGetter,
							dataResolver:      fakeDataResolver,
							apiClient:         fakeAPIClient,
							cliExiter:         fakeCliExiter,
							cliParamSatisfier: new(cliparamsatisfier.Fake),
						}

						/* act */
						objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

						/* assert */
						Expect(fakeCliExiter.ExitArgsForCall(0)).
							To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
					})
				})
				Context("apiClient.GetEventStream doesn't error", func() {
					Context("event channel closes", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeCliExiter := new(cliexiter.Fake)

							fakeOpFileGetter := new(opfile.FakeGetter)
							fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

							fakeOpHandle := new(data.FakeHandle)
							fakeDataResolver := new(dataresolver.Fake)
							fakeDataResolver.ResolveReturns(fakeOpHandle)

							fakeAPIClient := new(client.Fake)
							eventChannel := make(chan model.Event)
							close(eventChannel)
							fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

							objectUnderTest := _runer{
								opFileGetter:      fakeOpFileGetter,
								dataResolver:      fakeDataResolver,
								apiClient:         fakeAPIClient,
								cliExiter:         fakeCliExiter,
								cliParamSatisfier: new(cliparamsatisfier.Fake),
							}

							/* act */
							objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

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

										fakeOpFileGetter := new(opfile.FakeGetter)
										fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

										fakeOpHandle := new(data.FakeHandle)
										fakeDataResolver := new(dataresolver.Fake)
										fakeDataResolver.ResolveReturns(fakeOpHandle)

										fakeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpID, nil)

										objectUnderTest := _runer{
											opFileGetter:      fakeOpFileGetter,
											dataResolver:      fakeDataResolver,
											cliColorer:        clicolorer.New(),
											apiClient:         fakeAPIClient,
											cliExiter:         fakeCliExiter,
											cliOutput:         new(clioutput.Fake),
											cliParamSatisfier: new(cliparamsatisfier.Fake),
										}

										/* act/assert */
										objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})
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

										fakeOpFileGetter := new(opfile.FakeGetter)
										fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

										fakeOpHandle := new(data.FakeHandle)
										fakeDataResolver := new(dataresolver.Fake)
										fakeDataResolver.ResolveReturns(fakeOpHandle)

										fakeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpID, nil)

										objectUnderTest := _runer{
											opFileGetter:      fakeOpFileGetter,
											dataResolver:      fakeDataResolver,
											cliColorer:        clicolorer.New(),
											apiClient:         fakeAPIClient,
											cliExiter:         fakeCliExiter,
											cliOutput:         new(clioutput.Fake),
											cliParamSatisfier: new(cliparamsatisfier.Fake),
										}

										/* act/assert */
										objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})
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

										fakeOpFileGetter := new(opfile.FakeGetter)
										fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

										fakeOpHandle := new(data.FakeHandle)
										fakeDataResolver := new(dataresolver.Fake)
										fakeDataResolver.ResolveReturns(fakeOpHandle)

										fakeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpID, nil)

										objectUnderTest := _runer{
											opFileGetter:      fakeOpFileGetter,
											dataResolver:      fakeDataResolver,
											cliColorer:        clicolorer.New(),
											apiClient:         fakeAPIClient,
											cliExiter:         fakeCliExiter,
											cliOutput:         new(clioutput.Fake),
											cliParamSatisfier: new(cliparamsatisfier.Fake),
										}

										/* act/assert */
										objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})
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

										fakeOpFileGetter := new(opfile.FakeGetter)
										fakeOpFileGetter.GetReturns(&model.OpFile{}, nil)

										fakeOpHandle := new(data.FakeHandle)
										fakeDataResolver := new(dataresolver.Fake)
										fakeDataResolver.ResolveReturns(fakeOpHandle)

										fakeAPIClient := new(client.Fake)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEndedEvent
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEndedEvent.OpEnded.RootOpID, nil)

										objectUnderTest := _runer{
											opFileGetter:      fakeOpFileGetter,
											dataResolver:      fakeDataResolver,
											cliColorer:        clicolorer.New(),
											apiClient:         fakeAPIClient,
											cliExiter:         fakeCliExiter,
											cliOutput:         new(clioutput.Fake),
											cliParamSatisfier: new(cliparamsatisfier.Fake),
										}

										/* act/assert */
										objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})
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
