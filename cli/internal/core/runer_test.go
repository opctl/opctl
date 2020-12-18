package core

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clioutputFakes "github.com/opctl/opctl/cli/internal/clioutput/fakes"
	cliparamsatisfierFakes "github.com/opctl/opctl/cli/internal/cliparamsatisfier/fakes"
	dataresolver "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	cliModel "github.com/opctl/opctl/cli/internal/model"
	modelFakes "github.com/opctl/opctl/cli/internal/model/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
	clientFakes "github.com/opctl/opctl/sdks/go/node/api/client/fakes"
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

			dummyOpDataHandle := getDummyOpDataHandle()

			fakeDataResolver := new(dataresolver.FakeDataResolver)
			fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

			fakeNodeProvider := new(nodeprovider.Fake)
			// error to trigger return
			fakeNodeProvider.CreateNodeIfNotExistsReturns(nil, errors.New(""))

			objectUnderTest := _runer{
				dataResolver:      fakeDataResolver,
				cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
				nodeProvider:      fakeNodeProvider,
			}

			/* act */
			objectUnderTest.Run(context.TODO(), providedOpRef, &cliModel.RunOpts{})

			/* assert */
			actualOpRef, actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)
			Expect(actualOpRef).To(Equal(providedOpRef))
			Expect(actualPullCreds).To(BeNil())
		})
		Context("opfile.GetContent errors", func() {
			It("should return expected error", func() {
				/* arrange */

				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.GetContentReturns(nil, errors.New(""))

				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _runer{
					dataResolver:      fakeDataResolver,
					cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
				}

				/* act */
				err := objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

				/* assert */
				Expect(err).To(MatchError(""))
			})
		})
		Context("opfile.Get doesn't error", func() {
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

				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

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
				fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args, nil)

				objectUnderTest := _runer{
					dataResolver:      fakeDataResolver,
					cliParamSatisfier: fakeCliParamSatisfier,
					nodeProvider:      fakeNodeProvider,
				}

				/* act */
				objectUnderTest.Run(providedContext, "", &cliModel.RunOpts{})

				/* assert */
				actualCtx, actualArgs := fakeAPIClient.StartOpArgsForCall(0)
				Expect(actualCtx).To(Equal(expectedCtx))
				Expect(actualArgs).To(Equal(expectedArgs))
			})
			Context("apiClient.StartOp errors", func() {
				It("should return expected error", func() {
					/* arrange */
					returnedError := errors.New("dummyError")

					dummyOpDataHandle := getDummyOpDataHandle()

					fakeDataResolver := new(dataresolver.FakeDataResolver)
					fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

					// stub node provider
					fakeAPIClient := new(clientFakes.FakeClient)
					fakeAPIClient.StartOpReturns("dummyCallID", returnedError)

					fakeNodeHandle := new(modelFakes.FakeNodeHandle)
					fakeNodeHandle.APIClientReturns(fakeAPIClient)

					fakeNodeProvider := new(nodeprovider.Fake)
					fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

					objectUnderTest := _runer{
						dataResolver:      fakeDataResolver,
						cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
						nodeProvider:      fakeNodeProvider,
					}

					/* act */
					err := objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

					/* assert */
					Expect(err).To(MatchError(returnedError))
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

					dummyOpDataHandle := getDummyOpDataHandle()

					fakeDataResolver := new(dataresolver.FakeDataResolver)
					fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

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
						dataResolver:      fakeDataResolver,
						cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
						nodeProvider:      fakeNodeProvider,
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
					It("should return expected error", func() {
						/* arrange */
						returnedError := errors.New("dummyError")

						dummyOpDataHandle := getDummyOpDataHandle()

						fakeDataResolver := new(dataresolver.FakeDataResolver)
						fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

						fakeAPIClient := new(clientFakes.FakeClient)
						fakeAPIClient.GetEventStreamReturns(nil, returnedError)

						fakeNodeHandle := new(modelFakes.FakeNodeHandle)
						fakeNodeHandle.APIClientReturns(fakeAPIClient)

						fakeNodeProvider := new(nodeprovider.Fake)
						fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

						objectUnderTest := _runer{
							dataResolver:      fakeDataResolver,
							cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
							nodeProvider:      fakeNodeProvider,
						}

						/* act */
						err := objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

						/* assert */
						Expect(err).To(MatchError(returnedError))
					})
				})
				Context("apiClient.GetEventStream doesn't error", func() {
					Context("event channel closes", func() {
						It("should return expected error", func() {
							/* arrange */
							dummyOpDataHandle := getDummyOpDataHandle()

							fakeDataResolver := new(dataresolver.FakeDataResolver)
							fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

							fakeAPIClient := new(clientFakes.FakeClient)
							eventChannel := make(chan model.Event)
							close(eventChannel)
							fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

							fakeNodeHandle := new(modelFakes.FakeNodeHandle)
							fakeNodeHandle.APIClientReturns(fakeAPIClient)

							fakeNodeProvider := new(nodeprovider.Fake)
							fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

							objectUnderTest := _runer{
								dataResolver:      fakeDataResolver,
								cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
								nodeProvider:      fakeNodeProvider,
							}

							/* act */
							err := objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})

							/* assert */
							Expect(err).To(MatchError("Event channel closed unexpectedly"))
						})
					})
					Context("event channel doesn't close", func() {
						Context("event received", func() {
							rootCallID := "dummyRootCallID"
							Context("CallEnded", func() {
								Context("Outcome==SUCCEEDED", func() {
									It("should return expected error", func() {
										/* arrange */
										opEnded := model.Event{
											Timestamp: time.Now(),
											CallEnded: &model.CallEnded{
												Call: model.Call{
													ID: rootCallID,
												},
												Outcome: model.OpOutcomeSucceeded,
											},
										}

										dummyOpDataHandle := getDummyOpDataHandle()

										fakeDataResolver := new(dataresolver.FakeDataResolver)
										fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

										fakeAPIClient := new(clientFakes.FakeClient)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEnded.CallEnded.Call.ID, nil)

										fakeNodeHandle := new(modelFakes.FakeNodeHandle)
										fakeNodeHandle.APIClientReturns(fakeAPIClient)

										fakeNodeProvider := new(nodeprovider.Fake)
										fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

										objectUnderTest := _runer{
											dataResolver:      fakeDataResolver,
											cliOutput:         new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											nodeProvider:      fakeNodeProvider,
										}

										/* act/assert */
										err := objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})
										Expect(err).To(BeNil())
									})
								})
								Context("Outcome==KILLED", func() {
									It("should return expected error", func() {
										/* arrange */
										opEnded := model.Event{
											Timestamp: time.Now(),
											CallEnded: &model.CallEnded{
												Call: model.Call{
													ID: rootCallID,
												},
												Outcome: model.OpOutcomeKilled,
											},
										}

										dummyOpDataHandle := getDummyOpDataHandle()

										fakeDataResolver := new(dataresolver.FakeDataResolver)
										fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

										fakeAPIClient := new(clientFakes.FakeClient)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEnded.CallEnded.Call.ID, nil)

										fakeNodeHandle := new(modelFakes.FakeNodeHandle)
										fakeNodeHandle.APIClientReturns(fakeAPIClient)

										fakeNodeProvider := new(nodeprovider.Fake)
										fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

										objectUnderTest := _runer{
											dataResolver:      fakeDataResolver,
											cliOutput:         new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											nodeProvider:      fakeNodeProvider,
										}

										/* act/assert */
										err := objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})
										Expect(err).To(MatchError(&RunError{ExitCode: 137}))
									})

								})
								Context("Outcome==FAILED", func() {
									It("should return expected error", func() {
										/* arrange */
										opEnded := model.Event{
											Timestamp: time.Now(),
											CallEnded: &model.CallEnded{
												Call: model.Call{
													ID: rootCallID,
												},
												Outcome: model.OpOutcomeFailed,
											},
										}

										dummyOpDataHandle := getDummyOpDataHandle()

										fakeDataResolver := new(dataresolver.FakeDataResolver)
										fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

										fakeAPIClient := new(clientFakes.FakeClient)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEnded.CallEnded.Call.ID, nil)

										fakeNodeHandle := new(modelFakes.FakeNodeHandle)
										fakeNodeHandle.APIClientReturns(fakeAPIClient)

										fakeNodeProvider := new(nodeprovider.Fake)
										fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

										objectUnderTest := _runer{
											dataResolver:      fakeDataResolver,
											cliOutput:         new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											nodeProvider:      fakeNodeProvider,
										}

										/* act/assert */
										err := objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})
										Expect(err).To(MatchError(&RunError{ExitCode: 1}))
									})
								})
								Context("Outcome==?", func() {
									It("should return expected error", func() {
										/* arrange */
										opEnded := model.Event{
											Timestamp: time.Now(),
											CallEnded: &model.CallEnded{
												Call: model.Call{
													ID: rootCallID,
												},
												Outcome: "some unexpected outcome",
											},
										}

										dummyOpDataHandle := getDummyOpDataHandle()

										fakeDataResolver := new(dataresolver.FakeDataResolver)
										fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

										fakeAPIClient := new(clientFakes.FakeClient)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeAPIClient.GetEventStreamReturns(eventChannel, nil)
										fakeAPIClient.StartOpReturns(opEnded.CallEnded.Call.ID, nil)

										fakeNodeHandle := new(modelFakes.FakeNodeHandle)
										fakeNodeHandle.APIClientReturns(fakeAPIClient)

										fakeNodeProvider := new(nodeprovider.Fake)
										fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

										objectUnderTest := _runer{
											dataResolver:      fakeDataResolver,
											cliOutput:         new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											nodeProvider:      fakeNodeProvider,
										}

										/* act/assert */
										err := objectUnderTest.Run(context.TODO(), "", &cliModel.RunOpts{})
										Expect(err).To(MatchError(&RunError{ExitCode: 1}))
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
