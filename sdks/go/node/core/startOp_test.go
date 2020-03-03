package core

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/data/fakes"
	"github.com/opctl/opctl/sdks/go/data/provider"
	. "github.com/opctl/opctl/sdks/go/data/provider/fakes"
	uniquestringFakes "github.com/opctl/opctl/sdks/go/internal/uniquestring/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	. "github.com/opctl/opctl/sdks/go/opspec/opfile/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("core", func() {
	Context("StartOp", func() {
		It("should call data.NewGitProvider w/ expected args", func() {
			/* arrange */
			providedStartOpReq := model.StartOpReq{
				Op: model.StartOpReqOp{
					PullCreds: &model.PullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				},
			}
			providedDataCachePath := "providedDataCachePath"

			fakeData := new(FakeData)
			// err to trigger immediate return
			fakeData.ResolveReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _core{
				data:          fakeData,
				dataCachePath: providedDataCachePath,
			}

			/* act */
			objectUnderTest.StartOp(
				context.Background(),
				providedStartOpReq,
			)

			/* assert */
			actualCachePath,
				actualPullCreds := fakeData.NewGitProviderArgsForCall(0)

			Expect(actualCachePath).To(Equal(providedDataCachePath))
			Expect(actualPullCreds).To(Equal(providedStartOpReq.Op.PullCreds))
		})
		It("should call data.Resolve w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedStartOpReq := model.StartOpReq{
				Op: model.StartOpReqOp{
					Ref: "dummyOpRef",
				},
			}

			fakeData := new(FakeData)
			// err to trigger immediate return
			fakeData.ResolveReturns(nil, errors.New("dummyErr"))

			fsProvider := new(FakeProvider)
			fakeData.NewFSProviderReturns(fsProvider)

			gitProvider := new(FakeProvider)
			fakeData.NewGitProviderReturns(gitProvider)

			expectedProviders := []provider.Provider{
				fsProvider,
				gitProvider,
			}

			objectUnderTest := _core{
				data: fakeData,
			}

			/* act */
			objectUnderTest.StartOp(
				providedCtx,
				providedStartOpReq,
			)

			/* assert */
			actualCtx,
				actualOpRef,
				actualProviders := fakeData.ResolveArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualOpRef).To(Equal(providedStartOpReq.Op.Ref))
			Expect(actualProviders).To(Equal(expectedProviders))
		})
		Context("data.Resolve errs", func() {
			It("should return expected result", func() {

				/* arrange */
				providedCtx := context.Background()
				providedStartOpReq := model.StartOpReq{
					Op: model.StartOpReqOp{
						Ref: "dummyOpRef",
					},
				}

				fakeData := new(FakeData)
				expectedErr := errors.New("dummyErr")
				fakeData.ResolveReturns(nil, expectedErr)

				objectUnderTest := _core{
					data: fakeData,
				}

				/* act */
				_, actualErr := objectUnderTest.StartOp(
					providedCtx,
					providedStartOpReq,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("data.Resolve doesn't err", func() {
			It("should call data.Get w/ expected args", func() {
				/* arrange */
				providedCtx := context.Background()
				fakeData := new(FakeData)
				fakeDataHandle := new(modelFakes.FakeDataHandle)
				fakeData.ResolveReturns(fakeDataHandle, nil)

				fakeOpFileGetter := new(FakeGetter)
				// err to trigger immediate return
				fakeOpFileGetter.GetReturns(nil, errors.New("dummyError"))

				objectUnderTest := _core{
					data:                fakeData,
					opFileGetter:        fakeOpFileGetter,
					uniqueStringFactory: new(uniquestringFakes.FakeUniqueStringFactory),
				}

				/* act */
				objectUnderTest.StartOp(
					providedCtx,
					model.StartOpReq{},
				)

				/* assert */
				actualCtx,
					actualDataHandle := fakeOpFileGetter.GetArgsForCall(0)

				Expect(actualCtx).To(Equal(providedCtx))
				Expect(actualDataHandle).To(Equal(fakeDataHandle))
			})
			Context("data.Get errs", func() {
				It("should return expected error", func() {
					/* arrange */
					fakeData := new(FakeData)
					fakeDataHandle := new(modelFakes.FakeDataHandle)
					fakeData.ResolveReturns(fakeDataHandle, nil)

					fakeOpFileGetter := new(FakeGetter)
					expectedErr := errors.New("dummyError")
					fakeOpFileGetter.GetReturns(&model.OpFile{}, expectedErr)

					objectUnderTest := _core{
						data:                fakeData,
						opFileGetter:        fakeOpFileGetter,
						uniqueStringFactory: new(uniquestringFakes.FakeUniqueStringFactory),
					}

					/* act */
					_, actualErr := objectUnderTest.StartOp(
						context.Background(),
						model.StartOpReq{},
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("data.Get doesn't err", func() {
				It("should call caller.Call w/ expected args", func() {
					/* arrange */
					providedCtx := context.Background()
					providedArg1String := "dummyArg1Value"
					providedArg2Dir := "dummyArg2Value"
					providedArg3Dir := "dummyArg3Value"
					providedArg4Dir := "dummyArg4Value"
					providedReq := model.StartOpReq{
						Args: map[string]*model.Value{
							"dummyArg1Name": {String: &providedArg1String},
							"dummyArg2Name": {Dir: &providedArg2Dir},
							"dummyArg3Name": {Dir: &providedArg3Dir},
							"dummyArg4Name": {Dir: &providedArg4Dir},
						},
						Op: model.StartOpReqOp{
							Ref: "dummyOpRef",
						},
					}

					fakeData := new(FakeData)
					expectedOpHandle := new(modelFakes.FakeDataHandle)
					fakeData.ResolveReturns(expectedOpHandle, nil)

					opFile := &model.OpFile{
						Outputs: map[string]*model.Param{
							"dummyOutput1": nil,
							"dummyOutput2": nil,
						},
					}

					fakeOpFileGetter := new(FakeGetter)
					fakeOpFileGetter.GetReturns(opFile, nil)

					expectedSCGOpCall := &model.SCGOpCall{
						Ref:     expectedOpHandle.Ref(),
						Inputs:  map[string]interface{}{},
						Outputs: map[string]string{},
					}
					for name := range providedReq.Args {
						expectedSCGOpCall.Inputs[name] = ""
					}
					for name := range opFile.Outputs {
						expectedSCGOpCall.Outputs[name] = ""
					}

					expectedID := "expectedID"
					expectedRootOpID := expectedID
					fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
					fakeUniqueStringFactory.ConstructReturns(expectedID, nil)

					fakeCaller := new(FakeCaller)

					objectUnderTest := _core{
						caller:              fakeCaller,
						pubSub:              new(FakePubSub),
						data:                fakeData,
						opFileGetter:        fakeOpFileGetter,
						uniqueStringFactory: fakeUniqueStringFactory,
					}

					/* act */
					objectUnderTest.StartOp(
						providedCtx,
						providedReq,
					)

					/* assert */
					// Call happens in go routine; wait 500ms to allow it to occur
					time.Sleep(time.Millisecond * 500)
					actualCtx,
						actualID,
						actualScope,
						actualSCG,
						actualOpHandle,
						actualParentCallID,
						actualRootOpID := fakeCaller.CallArgsForCall(0)

					Expect(actualCtx).To(Equal(providedCtx))
					Expect(actualID).To(Equal(expectedID))
					Expect(actualScope).To(Equal(providedReq.Args))
					Expect(actualSCG).To(Equal(&model.SCG{Op: expectedSCGOpCall}))
					Expect(actualOpHandle).To(Equal(expectedOpHandle))
					Expect(actualParentCallID).To(BeNil())
					Expect(actualRootOpID).To(Equal(expectedRootOpID))
				})
			})
		})
	})
})
