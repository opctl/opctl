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
	"path"
	"time"
)

var _ = Describe("runOp", func() {

	Context("Execute", func() {
		It("should call exiter with expected args when bundle.GetOp returns error", func() {
			/* arrange */
			fakeExiter := new(fakeExiter)
			returnedError := errors.New("dummyError")

			fakeBundle := new(bundle.FakeBundle)
			fakeBundle.GetOpReturns(model.OpView{}, returnedError)

			fakeEngineClient := new(engineclient.FakeEngineClient)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				engineClient:      fakeEngineClient,
				exiter:            fakeExiter,
				paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
				workDirPathGetter: new(fakeWorkDirPathGetter),
			}

			/* act */
			objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyName")

			/* assert */
			Expect(fakeExiter.ExitArgsForCall(0)).
				Should(Equal(ExitReq{Message: returnedError.Error(), Code: 1}))
		})
		It("should call bundle.GetOp with expected args", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)

			fakeWorkDirPathGetter := new(fakeWorkDirPathGetter)
			workDirPath := "dummyWorkDirPath"
			fakeWorkDirPathGetter.GetReturns(workDirPath)

			fakeEngineClient := new(engineclient.FakeEngineClient)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

			fakeExiter := new(fakeExiter)

			providedName := "dummyOpName"
			providedCollection := "dummyCollection"

			expectedPath := path.Join(workDirPath, providedCollection, providedName)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				engineClient:      fakeEngineClient,
				exiter:            fakeExiter,
				paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
				workDirPathGetter: fakeWorkDirPathGetter,
			}

			/* act */
			objectUnderTest.RunOp([]string{}, providedCollection, providedName)

			/* assert */
			Expect(fakeBundle.GetOpArgsForCall(0)).Should(Equal(expectedPath))
		})
		It("should call exiter with expected args when bundle.GetEventStream returns error", func() {
			/* arrange */
			fakeExiter := new(fakeExiter)
			returnedError := errors.New("dummyError")

			fakeBundle := new(bundle.FakeBundle)
			fakeBundle.GetOpReturns(model.OpView{}, nil)

			fakeEngineClient := new(engineclient.FakeEngineClient)
			fakeEngineClient.GetEventStreamReturns(nil, returnedError)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				engineClient:      fakeEngineClient,
				exiter:            fakeExiter,
				paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
				workDirPathGetter: new(fakeWorkDirPathGetter),
			}

			/* act */
			objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

			/* assert */
			Expect(fakeExiter.ExitArgsForCall(0)).
				Should(Equal(ExitReq{Message: returnedError.Error(), Code: 1}))
		})
		Describe("when op has params defined", func() {
			Describe("and corresponding args are provided explicitly with values", func() {
				It("should call engineClient.StartOp with provided arg values", func() {
					/* arrange */
					fakeExiter := new(fakeExiter)
					param1Name := "DUMMY_PARAM1_NAME"
					param1Value := &model.Data{String: "dummyParam1Value"}

					fakeBundle := new(bundle.FakeBundle)
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

					fakeEngineClient := new(engineclient.FakeEngineClient)
					fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

					objectUnderTest := _core{
						bundle:            fakeBundle,
						engineClient:      fakeEngineClient,
						exiter:            fakeExiter,
						paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
						workDirPathGetter: new(fakeWorkDirPathGetter),
					}

					expectedArgs := map[string]*model.Data{param1Name: param1Value}
					providedArgs := []string{fmt.Sprintf("%v=%v", param1Name, param1Value.String)}

					/* act */
					objectUnderTest.RunOp(providedArgs, "dummyCollection", "dummyOpName")

					/* assert */
					Expect(fakeEngineClient.StartOpArgsForCall(0).Args).To(BeEquivalentTo(expectedArgs))
				})
			})
			Describe("and corresponding args are provided explicitly without values", func() {
				It("should call engineClient.StartOp with arg values obtained from the environment", func() {
					/* arrange */
					fakeExiter := new(fakeExiter)
					param1Name := "DUMMY_PARAM1_NAME"
					param1Value := &model.Data{String: "dummyParam1Value"}

					fakeVos := new(vos.FakeVos)
					fakeVos.GetenvReturns(param1Value.String)

					fakeBundle := new(bundle.FakeBundle)
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

					fakeEngineClient := new(engineclient.FakeEngineClient)
					fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

					objectUnderTest := _core{
						bundle:            fakeBundle,
						engineClient:      fakeEngineClient,
						exiter:            fakeExiter,
						paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), fakeVos),
						workDirPathGetter: new(fakeWorkDirPathGetter),
					}

					expectedArgs := map[string]*model.Data{param1Name: param1Value}
					providedArgs := []string{param1Name}

					/* act */
					objectUnderTest.RunOp(providedArgs, "dummyCollection", "dummyOpName")

					/* assert */
					Expect(fakeEngineClient.StartOpArgsForCall(0).Args).To(BeEquivalentTo(expectedArgs))
				})
			})
			Describe("and corresponding args are not provided", func() {
				Describe("and defaults don't exist", func() {
					It("should call bundle.RunOp with arg values obtained from the environment", func() {
						/* arrange */
						fakeExiter := new(fakeExiter)
						param1Name := "DUMMY_PARAM1_NAME"
						param1Value := &model.Data{String: "dummyParam1Value"}

						fakeVos := new(vos.FakeVos)
						fakeVos.GetenvReturns(param1Value.String)

						fakeBundle := new(bundle.FakeBundle)
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

						fakeEngineClient := new(engineclient.FakeEngineClient)
						fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

						objectUnderTest := _core{
							bundle:            fakeBundle,
							engineClient:      fakeEngineClient,
							exiter:            fakeExiter,
							paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), fakeVos),
							workDirPathGetter: new(fakeWorkDirPathGetter),
						}

						expectedArgs := map[string]*model.Data{param1Name: param1Value}
						providedArgs := []string{}

						/* act */
						objectUnderTest.RunOp(providedArgs, "dummyCollection", "dummyOpName")

						/* assert */
						Expect(fakeEngineClient.StartOpArgsForCall(0).Args).To(BeEquivalentTo(expectedArgs))
					})
				})
				Describe("but defaults exist", func() {
					It("should call bundle.RunOp without args for defaulted params", func() {
						/* arrange */
						fakeExiter := new(fakeExiter)
						// unique name to ensure conflicting env var not present
						param1Name := string(time.Now().Unix())

						fakeBundle := new(bundle.FakeBundle)
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

						fakeEngineClient := new(engineclient.FakeEngineClient)
						fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

						objectUnderTest := _core{
							bundle:            fakeBundle,
							engineClient:      fakeEngineClient,
							exiter:            fakeExiter,
							paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
							workDirPathGetter: new(fakeWorkDirPathGetter),
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
		It("should call exiter with expected args when engineClient.StartOp returns error", func() {
			/* arrange */
			fakeExiter := new(fakeExiter)
			returnedError := errors.New("dummyError")

			fakeBundle := new(bundle.FakeBundle)
			fakeBundle.GetOpReturns(model.OpView{}, nil)

			fakeEngineClient := new(engineclient.FakeEngineClient)
			fakeEngineClient.StartOpReturns("dummyOpId", returnedError)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				engineClient:      fakeEngineClient,
				exiter:            fakeExiter,
				paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
				workDirPathGetter: new(fakeWorkDirPathGetter),
			}

			/* act */
			objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

			/* assert */
			Expect(fakeExiter.ExitArgsForCall(0)).
				Should(Equal(ExitReq{Message: returnedError.Error(), Code: 1}))
		})
		It("should call exiter with expected args when event channel closes unexpectedly", func() {
			/* arrange */
			fakeExiter := new(fakeExiter)

			fakeBundle := new(bundle.FakeBundle)
			fakeBundle.GetOpReturns(model.OpView{}, nil)

			fakeEngineClient := new(engineclient.FakeEngineClient)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				engineClient:      fakeEngineClient,
				exiter:            fakeExiter,
				paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
				workDirPathGetter: new(fakeWorkDirPathGetter),
			}

			/* act */
			objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")

			/* assert */
			Expect(fakeExiter.ExitArgsForCall(0)).
				Should(Equal(ExitReq{Message: "Event channel closed unexpectedly", Code: 1}))
		})
		Describe("when an OpEndedEvent is received for our root op run", func() {
			opGraphId := "dummyOpGraphId"
			It("should call exiter with expected args when it's Outcome is SUCCEEDED", func() {
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

				fakeBundle := new(bundle.FakeBundle)
				fakeBundle.GetOpReturns(model.OpView{}, nil)

				fakeEngineClient := new(engineclient.FakeEngineClient)
				eventChannel := make(chan model.Event, 10)
				eventChannel <- opEndedEvent
				defer close(eventChannel)
				fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
				fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

				objectUnderTest := _core{
					bundle:            fakeBundle,
					colorer:           colorer.New(),
					engineClient:      fakeEngineClient,
					exiter:            fakeExiter,
					output:            newOutput(colorer.New()),
					paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
					workDirPathGetter: new(fakeWorkDirPathGetter),
				}

				/* act/assert */
				objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Code: 0}))
			})
			It("should call exiter with expected args when it's Outcome is KILLED", func() {
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

				fakeBundle := new(bundle.FakeBundle)
				fakeBundle.GetOpReturns(model.OpView{}, nil)

				fakeEngineClient := new(engineclient.FakeEngineClient)
				eventChannel := make(chan model.Event, 10)
				eventChannel <- opEndedEvent
				defer close(eventChannel)
				fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
				fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

				objectUnderTest := _core{
					bundle:            fakeBundle,
					colorer:           colorer.New(),
					engineClient:      fakeEngineClient,
					exiter:            fakeExiter,
					output:            newOutput(colorer.New()),
					paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
					workDirPathGetter: new(fakeWorkDirPathGetter),
				}

				/* act/assert */
				objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Code: 137}))
			})
			It("should call exiter with expected args when it's Outcome is unexpected", func() {
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

				fakeBundle := new(bundle.FakeBundle)
				fakeBundle.GetOpReturns(model.OpView{}, nil)

				fakeEngineClient := new(engineclient.FakeEngineClient)
				eventChannel := make(chan model.Event, 10)
				eventChannel <- opEndedEvent
				defer close(eventChannel)
				fakeEngineClient.GetEventStreamReturns(eventChannel, nil)
				fakeEngineClient.StartOpReturns(opEndedEvent.OpEnded.OpGraphId, nil)

				objectUnderTest := _core{
					bundle:            fakeBundle,
					colorer:           colorer.New(),
					engineClient:      fakeEngineClient,
					exiter:            fakeExiter,
					output:            newOutput(colorer.New()),
					paramSatisfier:    newParamSatisfier(colorer.New(), fakeExiter, validate.New(), new(vos.FakeVos)),
					workDirPathGetter: new(fakeWorkDirPathGetter),
				}

				/* act/assert */
				objectUnderTest.RunOp([]string{}, "dummyCollection", "dummyOpName")
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Code: 1}))
			})
		})
	})
})
