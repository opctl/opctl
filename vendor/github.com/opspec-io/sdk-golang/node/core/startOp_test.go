package core

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/core/containerruntime"
	"github.com/opspec-io/sdk-golang/op/dotyml"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"time"
)

var _ = Context("core", func() {
	Context("StartOp", func() {
		Context("req.Pkg nil", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("pkg required")

				objectUnderTest := _core{}

				/* act */
				_, actualErr := objectUnderTest.StartOp(
					context.Background(),
					model.StartOpReq{},
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("req.Pkg not nil", func() {
			It("should call data.Get w/ expected args", func() {
				/* arrange */
				fakeData := new(data.Fake)
				fakeDataHandle := new(data.FakeHandle)
				fakeData.ResolveReturns(fakeDataHandle, nil)

				fakeDotYmlGetter := new(dotyml.FakeGetter)
				// err to trigger immediate return
				fakeDotYmlGetter.GetReturns(nil, errors.New("dummyError"))

				objectUnderTest := _core{
					data:                fakeData,
					dotYmlGetter:        fakeDotYmlGetter,
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.StartOp(
					context.Background(),
					model.StartOpReq{Pkg: &model.DCGOpCallPkg{}},
				)

				/* assert */
				actualCtx,
					actualDataHandle := fakeDotYmlGetter.GetArgsForCall(0)

				Expect(actualCtx).To(Equal(actualCtx))
				Expect(actualDataHandle).To(Equal(fakeDataHandle))
			})
			Context("data.Get errs", func() {
				It("should return expected error", func() {
					/* arrange */
					fakeData := new(data.Fake)
					fakeDataHandle := new(data.FakeHandle)
					fakeData.ResolveReturns(fakeDataHandle, nil)

					fakeDotYmlGetter := new(dotyml.FakeGetter)
					expectedErr := errors.New("dummyError")
					fakeDotYmlGetter.GetReturns(&model.PkgManifest{}, expectedErr)

					objectUnderTest := _core{
						data:                fakeData,
						dotYmlGetter:        fakeDotYmlGetter,
						uniqueStringFactory: new(uniquestring.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.StartOp(
						context.Background(),
						model.StartOpReq{Pkg: &model.DCGOpCallPkg{}},
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("data.Get doesn't err", func() {
				It("should call opCaller.Call w/ expected args", func() {
					/* arrange */
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
						Pkg: &model.DCGOpCallPkg{
							Ref: "dummyPkgRef",
						},
					}

					fakeData := new(data.Fake)
					fakeDataHandle := new(data.FakeHandle)
					fakeData.ResolveReturns(fakeDataHandle, nil)

					pkgManifest := &model.PkgManifest{
						Outputs: map[string]*model.Param{
							"dummyOutput1": nil,
							"dummyOutput2": nil,
						},
					}

					fakeDotYmlGetter := new(dotyml.FakeGetter)
					fakeDotYmlGetter.GetReturns(pkgManifest, nil)

					expectedSCGOpCall := &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: fakeDataHandle.Ref(),
						},
						Inputs:  map[string]interface{}{},
						Outputs: map[string]string{},
					}
					for name := range providedReq.Args {
						expectedSCGOpCall.Inputs[name] = ""
					}
					for name := range pkgManifest.Outputs {
						expectedSCGOpCall.Outputs[name] = ""
					}

					expectedOpID := "dummyOpID"

					fakeOpCaller := new(fakeOpCaller)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					fakeUniqueStringFactory.ConstructReturns(expectedOpID, nil)

					objectUnderTest := _core{
						containerRuntime:    new(containerruntime.Fake),
						pubSub:              new(pubsub.Fake),
						data:                fakeData,
						dotYmlGetter:        fakeDotYmlGetter,
						opCaller:            fakeOpCaller,
						dcgNodeRepo:         new(fakeDCGNodeRepo),
						uniqueStringFactory: fakeUniqueStringFactory,
					}

					/* act */
					objectUnderTest.StartOp(
						context.Background(),
						providedReq,
					)

					/* assert */
					// Call happens in go routine; wait 500ms to allow it to occur
					time.Sleep(time.Millisecond * 500)
					actualInboundScope,
						actualOpID,
						actualOpHandle,
						actualRootOpID,
						actualSCGOpCall := fakeOpCaller.CallArgsForCall(0)

					Expect(actualInboundScope).To(Equal(providedReq.Args))
					Expect(actualOpID).To(Equal(expectedOpID))
					Expect(actualOpHandle).To(Equal(fakeDataHandle))
					Expect(actualRootOpID).To(Equal(actualOpID))
					Expect(actualSCGOpCall).To(Equal(expectedSCGOpCall))
				})
			})
		})
	})
})
