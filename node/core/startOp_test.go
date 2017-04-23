package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
	"path"
	"time"
)

var _ = Context("core", func() {
	Context("StartOp", func() {
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
				PkgRef: "/something/dummyPkg",
			}

			expectedPkgRef := path.Base(providedReq.PkgRef)
			expectedPkgBasePath := path.Dir(providedReq.PkgRef)

			expectedSCGOpCall := &model.SCGOpCall{
				Pkg: &model.SCGOpCallPkg{
					Ref: expectedPkgRef,
				},
				Inputs: map[string]string{},
			}
			for name := range providedReq.Args {
				// map as passed
				expectedSCGOpCall.Inputs[name] = name
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
