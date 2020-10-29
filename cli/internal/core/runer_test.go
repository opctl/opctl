package core

import (
	"context"
	"errors"
	"io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	clioutputFakes "github.com/opctl/opctl/cli/internal/clioutput/fakes"
	cliparamsatisfierFakes "github.com/opctl/opctl/cli/internal/cliparamsatisfier/fakes"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	cliModel "github.com/opctl/opctl/cli/internal/model"
	modelFakes "github.com/opctl/opctl/cli/internal/model/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/data/provider/fs"
	"github.com/opctl/opctl/sdks/go/model"
	clientFakes "github.com/opctl/opctl/sdks/go/node/api/client/fakes"
	opfileFakes "github.com/opctl/opctl/sdks/go/opspec/opfile/fakes"
)

var fsDataProvider = fs.New("testdata")

func getDummyOpDataHandle() model.DataHandle {
	dataHandle, err := fsDataProvider.TryResolve(context.TODO(), "dummy-op")
	if nil != err {
		panic(err)
	}
	return dataHandle
}

var _ = Context("Runer", func() {

	Context("Run", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedOpRef := "dummyOpRef"

			fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
			// err to trigger immediate return
			fakeOpFileUnmarshaller.UnmarshalReturns(nil, errors.New("dummyError"))

			dummyOpDataHandle := getDummyOpDataHandle()

			fakeDataResolver := new(dataresolver.Fake)
			fakeDataResolver.ResolveReturns(dummyOpDataHandle)

			objectUnderTest := _runer{
				opFileUnmarshaller: fakeOpFileUnmarshaller,
				dataResolver:       fakeDataResolver,
				cliExiter:          new(cliexiterFakes.FakeCliExiter),
				cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
			}

			/* act */
			objectUnderTest.Run(context.TODO(), providedOpRef, &cliModel.RunOpts{})

			/* assert */
			actualOpRef, actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)
			Expect(actualOpRef).To(Equal(providedOpRef))
			Expect(actualPullCreds).To(BeNil())
		})
		It("should call opfile.Get w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()

			fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
			// error to trigger immediate return
			fakeOpFileUnmarshaller.UnmarshalReturns(nil, errors.New("dummyError"))

			dummyOpDataHandle := getDummyOpDataHandle()

			fileReader, err := dummyOpDataHandle.GetContent(context.TODO(), "op.yml")
			defer fileReader.Close()
			if nil != err {
				panic(err)
			}

			fileBytes, err := ioutil.ReadAll(fileReader)
			if nil != err {
				panic(err)
			}

			fakeDataResolver := new(dataresolver.Fake)
			fakeDataResolver.ResolveReturns(dummyOpDataHandle)

			objectUnderTest := _runer{
				opFileUnmarshaller: fakeOpFileUnmarshaller,
				dataResolver:       fakeDataResolver,
				cliExiter:          new(cliexiterFakes.FakeCliExiter),
				cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
			}

			/* act */
			objectUnderTest.Run(
				providedCtx,
				"",
				new(cliModel.RunOpts),
			)

			/* assert */
			actualOpFileBytes := fakeOpFileUnmarshaller.UnmarshalArgsForCall(0)

			Expect(actualOpFileBytes).To(Equal(fileBytes))
		})
		Context("opfile.Get errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				getManifestErr := errors.New("dummyError")

				fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
				fakeOpFileUnmarshaller.UnmarshalReturns(nil, getManifestErr)

				dummyOpDataHandle := getDummyOpDataHandle()

				fakeDataResolver := new(dataresolver.Fake)
				fakeDataResolver.ResolveReturns(dummyOpDataHandle)

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _runer{
					opFileUnmarshaller: fakeOpFileUnmarshaller,
					dataResolver:       fakeDataResolver,
					cliExiter:          fakeCliExiter,
					cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
				}

				/* act */
				objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: getManifestErr.Error(), Code: 1}))
			})
		})
		Context("opfile.Get doesn't error", func() {
			It("should call paramSatisfier.Satisfy w/ expected args", func() {
				/* arrange */
				param1Name := "DUMMY_PARAM1_NAME"
				opFile := &model.OpSpec{
					Inputs: map[string]*model.Param{
						param1Name: {
							String: &model.StringParam{},
						},
					},
				}

				expectedParams := opFile.Inputs

				fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
				fakeOpFileUnmarshaller.UnmarshalReturns(opFile, nil)

				dummyOpDataHandle := getDummyOpDataHandle()

				fakeDataResolver := new(dataresolver.Fake)
				fakeDataResolver.ResolveReturns(dummyOpDataHandle)

				// stub node provider
				fakeAPIClient := new(clientFakes.FakeClient)

				// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

				fakeNodeHandle := new(modelFakes.FakeNodeHandle)
				fakeNodeHandle.APIClientReturns(fakeAPIClient)

				fakeNodeProvider := new(nodeprovider.Fake)
				fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

				fakeCliParamSatisfier := new(cliparamsatisfierFakes.FakeCLIParamSatisfier)

				objectUnderTest := _runer{
					opFileUnmarshaller: fakeOpFileUnmarshaller,
					dataResolver:       fakeDataResolver,
					cliExiter:          new(cliexiterFakes.FakeCliExiter),
					cliParamSatisfier:  fakeCliParamSatisfier,
					nodeProvider:       fakeNodeProvider,
				}

				/* act */
				objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

				/* assert */
				_, actualParams := fakeCliParamSatisfier.SatisfyArgsForCall(0)

				Expect(actualParams).To(Equal(expectedParams))
			})
			It("should call nodeHandle.APIClient().StartOp w/ expected args", func() {
				/* arrange */
				dummyOpDataHandle := getDummyOpDataHandle()

				providedContext := context.TODO()
				expectedCtx := providedContext

				expectedArg1ValueString := "dummyArg1Value"
				expectedArgs := model.StartOpReq{
					Args: map[string]*model.Value{
						"dummyArg1Name": {String: &expectedArg1ValueString},
					},
					Op: model.StartOpReqOp{
						Ref: dummyOpDataHandle.Ref(),
					},
				}

				fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
				fakeOpFileUnmarshaller.UnmarshalReturns(&model.OpSpec{}, nil)

				fakeDataResolver := new(dataresolver.Fake)
				fakeDataResolver.ResolveReturns(dummyOpDataHandle)

				// stub node provider
				fakeAPIClient := new(clientFakes.FakeClient)

				// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

				fakeNodeHandle := new(modelFakes.FakeNodeHandle)
				fakeNodeHandle.APIClientReturns(fakeAPIClient)

				fakeNodeProvider := new(nodeprovider.Fake)
				fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

				fakeCliParamSatisfier := new(cliparamsatisfierFakes.FakeCLIParamSatisfier)
				fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args)

				objectUnderTest := _runer{
					opFileUnmarshaller: fakeOpFileUnmarshaller,
					dataResolver:       fakeDataResolver,
					cliExiter:          new(cliexiterFakes.FakeCliExiter),
					cliParamSatisfier:  fakeCliParamSatisfier,
					nodeProvider:       fakeNodeProvider,
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
					fakeCliExiter := new(cliexiterFakes.FakeCliExiter)
					returnedError := errors.New("dummyError")

					fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
					fakeOpFileUnmarshaller.UnmarshalReturns(&model.OpSpec{}, nil)

					dummyOpDataHandle := getDummyOpDataHandle()

					fakeDataResolver := new(dataresolver.Fake)
					fakeDataResolver.ResolveReturns(dummyOpDataHandle)

					// stub node provider
					fakeAPIClient := new(clientFakes.FakeClient)
					fakeAPIClient.StartOpReturns("dummyCallID", returnedError)

					fakeNodeHandle := new(modelFakes.FakeNodeHandle)
					fakeNodeHandle.APIClientReturns(fakeAPIClient)

					fakeNodeProvider := new(nodeprovider.Fake)
					fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

					objectUnderTest := _runer{
						opFileUnmarshaller: fakeOpFileUnmarshaller,
						dataResolver:       fakeDataResolver,
						cliExiter:          fakeCliExiter,
						cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
						nodeProvider:       fakeNodeProvider,
					}

					/* act */
					objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
				})
			})
			Context("apiClient.StartOp doesn't error", func() {
				It("should call nodeHandle.APIClient().GetEventStream w/ expected args", func() {
					/* arrange */
					providedCtx := context.Background()
					rootCallIDReturnedFromStartOp := "dummyRootCallID"
					startTime := time.Now().UTC()
					expectedReq := &model.GetEventStreamReq{
						Filter: model.EventFilter{
							Roots: []string{rootCallIDReturnedFromStartOp},
							Since: &startTime,
						},
					}

					fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
					fakeOpFileUnmarshaller.UnmarshalReturns(&model.OpSpec{}, nil)

					dummyOpDataHandle := getDummyOpDataHandle()

					fakeDataResolver := new(dataresolver.Fake)
					fakeDataResolver.ResolveReturns(dummyOpDataHandle)

					// stub node provider
					fakeAPIClient := new(clientFakes.FakeClient)
					fakeAPIClient.StartOpReturns(rootCallIDReturnedFromStartOp, nil)

					fakeNodeHandle := new(modelFakes.FakeNodeHandle)
					fakeNodeHandle.APIClientReturns(fakeAPIClient)

					fakeNodeProvider := new(nodeprovider.Fake)
					fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _runer{
						opFileUnmarshaller: fakeOpFileUnmarshaller,
						dataResolver:       fakeDataResolver,
						cliExiter:          new(cliexiterFakes.FakeCliExiter),
						cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
						nodeProvider:       fakeNodeProvider,
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
						fakeCliExiter := new(cliexiterFakes.FakeCliExiter)
						returnedError := errors.New("dummyError")

						fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
						fakeOpFileUnmarshaller.UnmarshalReturns(&model.OpSpec{}, nil)

						dummyOpDataHandle := getDummyOpDataHandle()

						fakeDataResolver := new(dataresolver.Fake)
						fakeDataResolver.ResolveReturns(dummyOpDataHandle)

						fakeAPIClient := new(clientFakes.FakeClient)
						fakeAPIClient.GetEventStreamReturns(nil, returnedError)

						fakeNodeHandle := new(modelFakes.FakeNodeHandle)
						fakeNodeHandle.APIClientReturns(fakeAPIClient)

						fakeNodeProvider := new(nodeprovider.Fake)
						fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

						objectUnderTest := _runer{
							opFileUnmarshaller: fakeOpFileUnmarshaller,
							dataResolver:       fakeDataResolver,
							cliExiter:          fakeCliExiter,
							cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
							nodeProvider:       fakeNodeProvider,
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
							fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

							fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
							fakeOpFileUnmarshaller.UnmarshalReturns(&model.OpSpec{}, nil)

							dummyOpDataHandle := getDummyOpDataHandle()

							fakeDataResolver := new(dataresolver.Fake)
							fakeDataResolver.ResolveReturns(dummyOpDataHandle)

							fakeAPIClient := new(clientFakes.FakeClient)
							eventChannel := make(chan model.Event)
							close(eventChannel)
							fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

							fakeNodeHandle := new(modelFakes.FakeNodeHandle)
							fakeNodeHandle.APIClientReturns(fakeAPIClient)

							fakeNodeProvider := new(nodeprovider.Fake)
							fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

							objectUnderTest := _runer{
								opFileUnmarshaller: fakeOpFileUnmarshaller,
								dataResolver:       fakeDataResolver,
								cliExiter:          fakeCliExiter,
								cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
								nodeProvider:       fakeNodeProvider,
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
							rootCallID := "dummyRootCallID"
							Context("CallEnded", func() {
								Context("Outcome==SUCCEEDED", func() {
									It("should call exiter w/ expected args", func() {
										/* arrange */
										opEnded := model.Event{
											Timestamp: time.Now(),
											CallEnded: &model.CallEnded{
												CallID:     rootCallID,
												Outcome:    model.OpOutcomeSucceeded,
												RootCallID: rootCallID,
											},
										}

										fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

										fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
										fakeOpFileUnmarshaller.UnmarshalReturns(&model.OpSpec{}, nil)

										dummyOpDataHandle := getDummyOpDataHandle()

										fakeDataResolver := new(dataresolver.Fake)
										fakeDataResolver.ResolveReturns(dummyOpDataHandle)

										fakeAPIClient := new(clientFakes.FakeClient)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEnded.CallEnded.RootCallID, nil)

										fakeNodeHandle := new(modelFakes.FakeNodeHandle)
										fakeNodeHandle.APIClientReturns(fakeAPIClient)

										fakeNodeProvider := new(nodeprovider.Fake)
										fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

										objectUnderTest := _runer{
											opFileUnmarshaller: fakeOpFileUnmarshaller,
											dataResolver:       fakeDataResolver,
											cliColorer:         clicolorer.New(),
											cliExiter:          fakeCliExiter,
											cliOutput:          new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											nodeProvider:       fakeNodeProvider,
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
										opEnded := model.Event{
											Timestamp: time.Now(),
											CallEnded: &model.CallEnded{
												CallID:     rootCallID,
												Outcome:    model.OpOutcomeKilled,
												RootCallID: rootCallID,
											},
										}

										fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

										fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
										fakeOpFileUnmarshaller.UnmarshalReturns(&model.OpSpec{}, nil)

										dummyOpDataHandle := getDummyOpDataHandle()

										fakeDataResolver := new(dataresolver.Fake)
										fakeDataResolver.ResolveReturns(dummyOpDataHandle)

										fakeAPIClient := new(clientFakes.FakeClient)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEnded.CallEnded.RootCallID, nil)

										fakeNodeHandle := new(modelFakes.FakeNodeHandle)
										fakeNodeHandle.APIClientReturns(fakeAPIClient)

										fakeNodeProvider := new(nodeprovider.Fake)
										fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

										objectUnderTest := _runer{
											opFileUnmarshaller: fakeOpFileUnmarshaller,
											dataResolver:       fakeDataResolver,
											cliColorer:         clicolorer.New(),
											cliExiter:          fakeCliExiter,
											cliOutput:          new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											nodeProvider:       fakeNodeProvider,
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
										opEnded := model.Event{
											Timestamp: time.Now(),
											CallEnded: &model.CallEnded{
												CallID:     rootCallID,
												Outcome:    model.OpOutcomeFailed,
												RootCallID: rootCallID,
											},
										}

										fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

										fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
										fakeOpFileUnmarshaller.UnmarshalReturns(&model.OpSpec{}, nil)

										dummyOpDataHandle := getDummyOpDataHandle()

										fakeDataResolver := new(dataresolver.Fake)
										fakeDataResolver.ResolveReturns(dummyOpDataHandle)

										fakeAPIClient := new(clientFakes.FakeClient)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEnded.CallEnded.RootCallID, nil)

										fakeNodeHandle := new(modelFakes.FakeNodeHandle)
										fakeNodeHandle.APIClientReturns(fakeAPIClient)

										fakeNodeProvider := new(nodeprovider.Fake)
										fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

										objectUnderTest := _runer{
											opFileUnmarshaller: fakeOpFileUnmarshaller,
											dataResolver:       fakeDataResolver,
											cliColorer:         clicolorer.New(),
											cliExiter:          fakeCliExiter,
											cliOutput:          new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											nodeProvider:       fakeNodeProvider,
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
										opEnded := model.Event{
											Timestamp: time.Now(),
											CallEnded: &model.CallEnded{
												CallID:     rootCallID,
												Outcome:    "some unexpected outcome",
												RootCallID: rootCallID,
											},
										}

										fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

										fakeOpFileUnmarshaller := new(opfileFakes.FakeUnmarshaller)
										fakeOpFileUnmarshaller.UnmarshalReturns(&model.OpSpec{}, nil)

										dummyOpDataHandle := getDummyOpDataHandle()

										fakeDataResolver := new(dataresolver.Fake)
										fakeDataResolver.ResolveReturns(dummyOpDataHandle)

										fakeAPIClient := new(clientFakes.FakeClient)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEnded.CallEnded.RootCallID, nil)

										fakeNodeHandle := new(modelFakes.FakeNodeHandle)
										fakeNodeHandle.APIClientReturns(fakeAPIClient)

										fakeNodeProvider := new(nodeprovider.Fake)
										fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

										objectUnderTest := _runer{
											opFileUnmarshaller: fakeOpFileUnmarshaller,
											dataResolver:       fakeDataResolver,
											cliColorer:         clicolorer.New(),
											cliExiter:          fakeCliExiter,
											cliOutput:          new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier:  new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											nodeProvider:       fakeNodeProvider,
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
