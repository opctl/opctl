package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
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
		Context("req.Pkg not nill", func() {
			It("should call opCaller.Call w/ expected args", func() {
				/* arrange */

				providedArg1String := "dummyArg1Value"
				providedArg2Dir := "dummyArg2Value"
				providedArg3Dir := "dummyArg3Value"
				providedArg4Dir := "dummyArg4Value"
				providedReq := model.StartOpReq{
					Args: map[string]*model.Data{
						"dummyArg1Name": {String: &providedArg1String},
						"dummyArg2Name": {Dir: &providedArg2Dir},
						"dummyArg3Name": {Dir: &providedArg3Dir},
						"dummyArg4Name": {Dir: &providedArg4Dir},
					},
					Pkg: &model.DCGOpCallPkg{
						Ref: "/something/dummyPkg",
					},
				}

				expectedPkgRef := filepath.Base(providedReq.Pkg.Ref)
				expectedPkgBasePath := filepath.Dir(providedReq.Pkg.Ref)

				expectedSCGOpCall := &model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref: expectedPkgRef,
					},
					Inputs: map[string]string{},
				}
				for name := range providedReq.Args {
					// map as passed
					expectedSCGOpCall.Inputs[name] = ""
				}

				expectedOpId := "dummyOpId"

				fakeOpCaller := new(fakeOpCaller)

				fakeUniqueStringFactory := new(uniquestring.Fake)
				fakeUniqueStringFactory.ConstructReturns(expectedOpId)

				objectUnderTest := _core{
					containerProvider:   new(containerprovider.Fake),
					pubSub:              new(pubsub.Fake),
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
					actualPkgBasePath,
					actualRootOpId,
					actualSCGOpCall := fakeOpCaller.CallArgsForCall(0)

				Expect(actualInboundScope).To(Equal(providedReq.Args))
				Expect(actualOpId).To(Equal(expectedOpId))
				Expect(actualPkgBasePath).To(Equal(expectedPkgBasePath))
				Expect(actualRootOpId).To(Equal(actualOpId))
				Expect(actualSCGOpCall).To(Equal(expectedSCGOpCall))
			})
		})
	})
})
