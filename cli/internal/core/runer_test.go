package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clioutputFakes "github.com/opctl/opctl/cli/internal/clioutput/fakes"
	cliparamsatisfierFakes "github.com/opctl/opctl/cli/internal/cliparamsatisfier/fakes"
	dataresolver "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

var fsDataProvider = fs.New("testdata")

func getDummyOpDataHandle() model.DataHandle {
	dataHandle, err := fsDataProvider.TryResolve(context.TODO(), "dummy-op")
	if nil != err {
		panic(err)
	}
	return dataHandle
}

// errReadSeekCloser is a mock ReadSeekCloser that returns an error on Read
type errReadSeekCloser struct {
	err error // if not specified, will panic with bytes.ErrToLarge
}

func (e errReadSeekCloser) Read(p []byte) (n int, err error) {
	return 0, e.err
}
func (errReadSeekCloser) Close() error {
	return nil
}
func (errReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

type mockReadSeekCloser struct {
	io.ReadSeeker
}

func (mockReadSeekCloser) Close() error {
	return errors.New("not implemented")
}

var _ = Context("Runer", func() {
	It("can be constructed", func() {
		newRuner(
			new(clioutputFakes.FakeCliOutput),
			new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
			new(dataresolver.FakeDataResolver),
			new(nodeFakes.FakeOpNode),
		)
	})

	Context("Run", func() {
		It("dataResolver.Resolve call", func() {
			/* arrange */
			providedOpRef := "dummyOpRef"

			expected := errors.New("data resolution error")
			fakeDataResolver := new(dataresolver.FakeDataResolver)
			fakeDataResolver.ResolveReturns(nil, expected)

			objectUnderTest := _runer{
				dataResolver:      fakeDataResolver,
				cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
				core:              new(nodeFakes.FakeOpNode),
			}

			/* act */
			err := objectUnderTest.Run(context.TODO(), providedOpRef, &RunOpts{})

			/* assert */
			Expect(err).To(MatchError(expected))
			actualOpRef, actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)
			Expect(actualOpRef).To(Equal(providedOpRef))
			Expect(actualPullCreds).To(BeNil())
		})
		It("opfile.GetContent call", func() {
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
			err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})

			/* assert */
			Expect(err).To(MatchError(""))
		})
		Context("opfile.Get doesn't error", func() {
			It("opfile.GetContent reader failure", func() {
				/* arrange */
				expectedError := errors.New("expected")

				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.GetContentReturns(errReadSeekCloser{err: expectedError}, nil)

				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _runer{
					dataResolver:      fakeDataResolver,
					cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
				}

				/* act */
				err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})

				/* assert */
				Expect(err).To(MatchError(expectedError))
			})
			It("opfile.Unmarshal failure", func() {
				/* arrange */
				fakeOpHandle := new(FakeDataHandle)
				rs := bytes.NewReader([]byte("garbage"))
				fakeOpHandle.GetContentReturns(mockReadSeekCloser{rs}, nil)

				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _runer{
					dataResolver:      fakeDataResolver,
					cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
				}

				/* act */
				err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})

				/* assert */
				Expect(err).NotTo(BeNil())
			})
			It("cliParamSatisfier yml file failure", func() {
				/* arrange */
				expectedError := errors.New("expected")
				dummyOpDataHandle := getDummyOpDataHandle()
				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)
				fakeCliParamSatisfier := new(cliparamsatisfierFakes.FakeCLIParamSatisfier)
				fakeCliParamSatisfier.NewYMLFileInputSrcReturns(nil, expectedError)

				objectUnderTest := _runer{
					dataResolver:      fakeDataResolver,
					cliParamSatisfier: fakeCliParamSatisfier,
				}

				/* act */
				err := objectUnderTest.Run(context.TODO(), "", &RunOpts{ArgFile: "argfile"})

				/* assert */
				Expect(err).To(MatchError(fmt.Errorf("unable to load arg file at '%v'; error was: %v", "argfile", expectedError)))
			})
			It("cliParamSatisfier satisfaction failure", func() {
				/* arrange */
				expectedError := errors.New("expected")
				dummyOpDataHandle := getDummyOpDataHandle()
				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)
				fakeCliParamSatisfier := new(cliparamsatisfierFakes.FakeCLIParamSatisfier)
				fakeCliParamSatisfier.SatisfyReturns(nil, expectedError)

				objectUnderTest := _runer{
					dataResolver:      fakeDataResolver,
					cliParamSatisfier: fakeCliParamSatisfier,
				}

				/* act */
				err := objectUnderTest.Run(context.TODO(), "", &RunOpts{ArgFile: "argfile"})

				/* assert */
				Expect(err).To(MatchError(expectedError))
			})
			It("should call core.StartOp w/ expected args", func() {
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

				fakeCore := new(nodeFakes.FakeOpNode)

				// stub GetEventStream w/ closed channel so test doesn't wait for events indefinitely
				eventChannel := make(chan model.Event)
				close(eventChannel)
				fakeCore.GetEventStreamReturns(eventChannel, nil)

				fakeCliParamSatisfier := new(cliparamsatisfierFakes.FakeCLIParamSatisfier)
				fakeCliParamSatisfier.SatisfyReturns(expectedArgs.Args, nil)

				objectUnderTest := _runer{
					dataResolver:      fakeDataResolver,
					cliParamSatisfier: fakeCliParamSatisfier,
					core:              fakeCore,
				}

				/* act */
				objectUnderTest.Run(providedContext, "", &RunOpts{})

				/* assert */
				actualCtx, actualArgs := fakeCore.StartOpArgsForCall(0)
				Expect(actualCtx).To(Equal(expectedCtx))
				Expect(actualArgs).To(Equal(expectedArgs))
			})
			Context("core.StartOp errors", func() {
				It("should return expected error", func() {
					/* arrange */
					returnedError := errors.New("dummyError")

					dummyOpDataHandle := getDummyOpDataHandle()

					fakeDataResolver := new(dataresolver.FakeDataResolver)
					fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

					fakeCore := new(nodeFakes.FakeOpNode)
					fakeCore.StartOpReturns("dummyCallID", returnedError)

					objectUnderTest := _runer{
						dataResolver:      fakeDataResolver,
						cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
						core:              fakeCore,
					}

					/* act */
					err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})

					/* assert */
					Expect(err).To(MatchError(returnedError))
				})
			})
			Context("core.StartOp doesn't error", func() {
				It("should call core.GetEventStream w/ expected args", func() {
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

					fakeCore := new(nodeFakes.FakeOpNode)
					fakeCore.StartOpReturns(rootCallIDReturnedFromStartOp, nil)

					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeCore.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _runer{
						dataResolver:      fakeDataResolver,
						cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
						core:              fakeCore,
					}

					/* act */
					objectUnderTest.Run(providedCtx, "", &RunOpts{})

					/* assert */
					actualCtx,
						actualReq := fakeCore.GetEventStreamArgsForCall(0)

					// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
					Expect(*actualReq.Filter.Since).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
					// set temporal fields to expected vals since they're already asserted
					actualReq.Filter.Since = &startTime

					Expect(actualCtx).To(Equal(providedCtx))
					Expect(actualReq).To(Equal(expectedReq))
				})
				Context("core.GetEventStream errors", func() {
					It("should return expected error", func() {
						/* arrange */
						returnedError := errors.New("dummyError")

						dummyOpDataHandle := getDummyOpDataHandle()

						fakeDataResolver := new(dataresolver.FakeDataResolver)
						fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

						fakeCore := new(nodeFakes.FakeOpNode)
						fakeCore.GetEventStreamReturns(nil, returnedError)

						objectUnderTest := _runer{
							dataResolver:      fakeDataResolver,
							cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
							core:              fakeCore,
						}

						/* act */
						err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})

						/* assert */
						Expect(err).To(MatchError(returnedError))
					})
				})
				Context("core.GetEventStream doesn't error", func() {
					Context("event channel closes", func() {
						It("should return expected error", func() {
							/* arrange */
							dummyOpDataHandle := getDummyOpDataHandle()

							fakeDataResolver := new(dataresolver.FakeDataResolver)
							fakeDataResolver.ResolveReturns(dummyOpDataHandle, nil)

							fakeCore := new(nodeFakes.FakeOpNode)
							eventChannel := make(chan model.Event)
							close(eventChannel)
							fakeCore.GetEventStreamReturns(eventChannel, nil)

							objectUnderTest := _runer{
								dataResolver:      fakeDataResolver,
								cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
								core:              fakeCore,
							}

							/* act */
							err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})

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

										fakeCore := new(nodeFakes.FakeOpNode)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeCore.GetEventStreamReturns(eventChannel, nil)
										fakeCore.StartOpReturns(opEnded.CallEnded.Call.ID, nil)

										objectUnderTest := _runer{
											dataResolver:      fakeDataResolver,
											cliOutput:         new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											core:              fakeCore,
										}

										/* act/assert */
										err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})
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

										fakeCore := new(nodeFakes.FakeOpNode)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeCore.GetEventStreamReturns(eventChannel, nil)
										fakeCore.StartOpReturns(opEnded.CallEnded.Call.ID, nil)

										objectUnderTest := _runer{
											dataResolver:      fakeDataResolver,
											cliOutput:         new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											core:              fakeCore,
										}

										/* act/assert */
										err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})
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

										fakeCore := new(nodeFakes.FakeOpNode)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeCore.GetEventStreamReturns(eventChannel, nil)
										fakeCore.StartOpReturns(opEnded.CallEnded.Call.ID, nil)

										objectUnderTest := _runer{
											dataResolver:      fakeDataResolver,
											cliOutput:         new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											core:              fakeCore,
										}

										/* act/assert */
										err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})
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

										fakeCore := new(nodeFakes.FakeOpNode)
										eventChannel := make(chan model.Event, 10)
										eventChannel <- opEnded
										defer close(eventChannel)
										fakeCore.GetEventStreamReturns(eventChannel, nil)
										fakeCore.StartOpReturns(opEnded.CallEnded.Call.ID, nil)

										objectUnderTest := _runer{
											dataResolver:      fakeDataResolver,
											cliOutput:         new(clioutputFakes.FakeCliOutput),
											cliParamSatisfier: new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
											core:              fakeCore,
										}

										/* act/assert */
										err := objectUnderTest.Run(context.TODO(), "", &RunOpts{})
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
