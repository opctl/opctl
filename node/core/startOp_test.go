package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/util/containerprovider"
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
				_, actualErr := objectUnderTest.StartOp(model.StartOpReq{})

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("req.Pkg not nil", func() {
			It("should call pkg.GetManifest w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)

				fakePkgHandle := new(pkg.FakeHandle)
				fakePkg.ResolveReturns(fakePkgHandle, nil)

				// err to trigger immediate return
				fakePkg.GetManifestReturns(nil, errors.New("dummyError"))

				objectUnderTest := _core{
					pkg:                 fakePkg,
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.StartOp(model.StartOpReq{Pkg: &model.DCGOpCallPkg{}})

				/* assert */

				Expect(fakePkg.GetManifestArgsForCall(0)).To(Equal(fakePkgHandle))
			})
			Context("pkg.GetManifest errs", func() {
				It("should return expected error", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)

					fakePkgHandle := new(pkg.FakeHandle)
					fakePkg.ResolveReturns(fakePkgHandle, nil)

					expectedErr := errors.New("dummyError")
					fakePkg.GetManifestReturns(&model.PkgManifest{}, expectedErr)

					objectUnderTest := _core{
						pkg:                 fakePkg,
						uniqueStringFactory: new(uniquestring.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.StartOp(model.StartOpReq{Pkg: &model.DCGOpCallPkg{}})

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("pkg.GetManifest doesn't err", func() {
				It("should call opCaller.Call w/ expected args", func() {
					/* arrange */
					arrayInputName := "dummyArrayInputName"
					arrayValue := []interface{}{}

					dirInputName := "dummyDirInputName"
					dirValue := "dummyDirValue"

					fileInputName := "dummyFileInputName"
					fileValue := "dummyFileValue"

					numberInputName := "dummyNumberInputName"
					numberValue := 2.0

					objectInputName := "dummyObjectInputName"
					objectValue := map[string]interface{}{}

					socketInputName := "dummySocketInputName"
					socketValue := "dummySocketValue"

					stringInputName := "dummyStringInputName"
					stringValue := "dummyStringValue"

					providedReq := model.StartOpReq{
						Args: map[string]interface{}{
							arrayInputName:  arrayValue,
							dirInputName:    dirValue,
							fileInputName:   fileValue,
							numberInputName: numberValue,
							objectInputName: objectValue,
							socketInputName: socketValue,
							stringInputName: stringValue,
						},
						Pkg: &model.DCGOpCallPkg{
							Ref: "dummyPkgRef",
						},
					}

					fakePkg := new(pkg.Fake)

					fakePkgHandle := new(pkg.FakeHandle)
					fakePkg.ResolveReturns(fakePkgHandle, nil)

					pkgManifest := &model.PkgManifest{
						Inputs: map[string]*model.Param{
							arrayInputName:  &model.Param{Array: &model.ArrayParam{}},
							dirInputName:    &model.Param{Dir: &model.DirParam{}},
							fileInputName:   &model.Param{File: &model.FileParam{}},
							numberInputName: &model.Param{Number: &model.NumberParam{}},
							objectInputName: &model.Param{Object: &model.ObjectParam{}},
							socketInputName: &model.Param{Socket: &model.SocketParam{}},
							stringInputName: &model.Param{String: &model.StringParam{}},
						},
						Outputs: map[string]*model.Param{
							"dummyOutput1": nil,
							"dummyOutput2": nil,
						},
					}
					fakePkg.GetManifestReturns(pkgManifest, nil)

					expectedSCGOpCall := &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: fakePkgHandle.Ref(),
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

					expectedInboundScope := map[string]*model.Value{
						arrayInputName:  &model.Value{Array: arrayValue},
						dirInputName:    &model.Value{Dir: &dirValue},
						fileInputName:   &model.Value{File: &fileValue},
						numberInputName: &model.Value{Number: &numberValue},
						objectInputName: &model.Value{Object: objectValue},
						socketInputName: &model.Value{Socket: &socketValue},
						stringInputName: &model.Value{String: &stringValue},
					}

					expectedOpId := "dummyOpId"

					fakeOpCaller := new(fakeOpCaller)

					fakeUniqueStringFactory := new(uniquestring.Fake)
					fakeUniqueStringFactory.ConstructReturns(expectedOpId)

					objectUnderTest := _core{
						containerProvider:   new(containerprovider.Fake),
						pubSub:              new(pubsub.Fake),
						pkg:                 fakePkg,
						opCaller:            fakeOpCaller,
						dcgNodeRepo:         new(fakeDCGNodeRepo),
						uniqueStringFactory: fakeUniqueStringFactory,
					}

					/* act */
					objectUnderTest.StartOp(providedReq)

					/* assert */
					// Call happens in go routine; wait 500ms to allow it to occur
					time.Sleep(time.Millisecond * 500)
					actualInboundScope,
						actualOpId,
						actualPkgHandle,
						actualRootOpId,
						actualSCGOpCall := fakeOpCaller.CallArgsForCall(0)

					Expect(actualInboundScope).To(Equal(expectedInboundScope))
					Expect(actualOpId).To(Equal(expectedOpId))
					Expect(actualPkgHandle).To(Equal(fakePkgHandle))
					Expect(actualRootOpId).To(Equal(actualOpId))
					Expect(actualSCGOpCall).To(Equal(expectedSCGOpCall))
				})
			})
		})
	})
})
