package core

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/core/containerruntime"
	"github.com/opspec-io/sdk-golang/pkg"
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
			It("should call data.GetManifest w/ expected args", func() {
				/* arrange */
				fakeData := new(data.Fake)
				fakeDataHandle := new(data.FakeHandle)
				fakeData.ResolveReturns(fakeDataHandle, nil)

				fakePkg := new(pkg.Fake)
				// err to trigger immediate return
				fakePkg.GetManifestReturns(nil, errors.New("dummyError"))

				objectUnderTest := _core{
					data:                fakeData,
					pkg:                 fakePkg,
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.StartOp(
					context.Background(),
					model.StartOpReq{Pkg: &model.DCGOpCallPkg{}},
				)

				/* assert */
				Expect(fakePkg.GetManifestArgsForCall(0)).To(Equal(fakeDataHandle))
			})
			Context("data.GetManifest errs", func() {
				It("should return expected error", func() {
					/* arrange */
					fakeData := new(data.Fake)
					fakeDataHandle := new(data.FakeHandle)
					fakeData.ResolveReturns(fakeDataHandle, nil)

					fakePkg := new(pkg.Fake)
					expectedErr := errors.New("dummyError")
					fakePkg.GetManifestReturns(&model.PkgManifest{}, expectedErr)

					objectUnderTest := _core{
						data:                fakeData,
						pkg:                 fakePkg,
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
			Context("data.GetManifest doesn't err", func() {
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

					fakePkg := new(pkg.Fake)
					fakePkg.GetManifestReturns(pkgManifest, nil)

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

					expectedOpId := "dummyOpId"

					fakeOpCaller := new(fakeOpCaller)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					fakeUniqueStringFactory.ConstructReturns(expectedOpId, nil)

					objectUnderTest := _core{
						containerRuntime:    new(containerruntime.Fake),
						pubSub:              new(pubsub.Fake),
						data:                fakeData,
						pkg:                 fakePkg,
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
						actualOpId,
						actualOpDirHandle,
						actualRootOpId,
						actualSCGOpCall := fakeOpCaller.CallArgsForCall(0)

					Expect(actualInboundScope).To(Equal(providedReq.Args))
					Expect(actualOpId).To(Equal(expectedOpId))
					Expect(actualOpDirHandle).To(Equal(fakeDataHandle))
					Expect(actualRootOpId).To(Equal(actualOpId))
					Expect(actualSCGOpCall).To(Equal(expectedSCGOpCall))
				})
			})
		})
	})
})
